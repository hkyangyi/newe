// Package ws 提供“主管理器 + 分组管理器”的 WebSocket 框架：
// - MainMag：仅管理分组（GroupMag）的生命周期，不直接持有连接
// - GroupMag：管理该分组下的所有连接（Client），广播严格限定在分组内
// - Client：连接实体，内置回调（Callback）处理该连接的业务消息
// 设计要点：
// 1) 路由决定使用哪个分组（MainMag.Get(name)）
// 2) 每条连接在 Serve 时绑定一个专属 Callback，互不干扰
// 3) 心跳策略：服务器周期性 Ping；文本 "ping" -> "pong" 也被兼容
package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// 基本参数（超时、缓冲等）
const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 2 << 20
	sendBufSize    = 256
)

// Upgrader：按需配置来源校验
// Upgrader：用于将 HTTP 连接升级为 WebSocket；默认放开跨域校验（生产可自定义）
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// 统一的消息格式
type Frame struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
	T    int64       `json:"t"`
}

// -------- Client --------

// Client 表示一个 WebSocket 连接（属于某个 GroupMag）。
// - Key：该连接在分组内的唯一标识（例如用户/会话 ID 组合）
// - Send：向该连接异步发送消息的通道（TextMessage）
// - Callback：每条消息进入时的业务处理回调（按连接维度独享）
type Client struct {
	Key      string
	Conn     *websocket.Conn
	Send     chan []byte
	Callback func(*Client, []byte)
	group    *GroupMag
	closed   bool
	mu       sync.Mutex
}

// ReadPump 持续读取客户端消息：
// - 设置 Read 限制与 Pong 处理（维持心跳）
// - 文本 "ping" 将被立即回复 "pong"
// - 其他消息进入该连接的 Callback 进行业务处理
func (c *Client) ReadPump() {
	defer c.closeWithReason(nil)
	c.Conn.SetReadLimit(maxMessageSize)
	_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	c.Conn.SetCloseHandler(func(int, string) error { return nil })

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		fmt.Println("ws message recv:", string(message))
		var f Frame
		//解析消息
		if err := json.Unmarshal(message, &f); err != nil {
			//解析失败当纯文本处理
			f = Frame{
				Type: "text",
				Data: string(message),
				T:    time.Now().Unix(),
			}

		}

		// 兼容纯文本心跳
		if f.Type == "ping" {
			_ = c.safeWrite(websocket.TextMessage, []byte("pong"))
			continue
		}
		if c.Callback != nil {
			c.Callback(c, message)
		}
	}
}

// WritePump 持续写入客户端数据：
// - 定期发送 Ping 维持心跳（服务端驱动）
// - 从 Send 通道消费消息并写入客户端
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() { ticker.Stop(); c.closeWithReason(nil) }()
	for {
		select {
		case msg, ok := <-c.Send:
			fmt.Println("ws message send:", string(msg))
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// safeWrite 在持有互斥锁的情况下写入消息，避免并发写冲突
func (c *Client) safeWrite(t int, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.Conn.WriteMessage(t, data)
}

// closeWithReason 关闭连接并通知所在分组注销该客户端；
// err 参数仅用于问题定位（当前未保存，可在业务中扩展日志）
func (c *Client) closeWithReason(err error) {
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return
	}
	c.closed = true
	c.mu.Unlock()
	if c.group != nil {
		c.group.unregister <- c
	}
	_ = c.Conn.Close()
	// 回调由业务自行在 Callback 中处理需要的清理
}

// -------- GroupMag（单分组管理器） --------

// GroupMag 管理单个分组内的所有连接：
// - register/unregister：连接的注册与注销
// - broadcast：向分组内所有连接发送消息
// - stop：停止该分组并关闭所有连接
type GroupMag struct {
	name       string
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	stop       chan struct{}
	mu         sync.RWMutex
}

// NewGroup 创建一个新的分组管理器（未启动 Run 循环）
func NewGroup(name string) *GroupMag {
	return &GroupMag{
		name:       name,
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte, 64),
		stop:       make(chan struct{}),
	}
}

// Run 启动分组主循环：串行处理注册/注销/广播/停止，确保状态一致
func (g *GroupMag) Run() {
	for {
		select {
		case c := <-g.register:
			g.mu.Lock()
			if old, ok := g.clients[c.Key]; ok {
				old.closeWithReason(errors.New("duplicate key"))
			}
			g.clients[c.Key] = c
			g.mu.Unlock()
		case c := <-g.unregister:
			g.mu.Lock()
			if cur, ok := g.clients[c.Key]; ok && cur == c {
				delete(g.clients, c.Key)
			}
			g.mu.Unlock()
			closeSafe(c.Send)
		case msg := <-g.broadcast:
			g.mu.RLock()
			for _, cli := range g.clients {
				select {
				case cli.Send <- msg:
				default:
					go cli.closeWithReason(errors.New("send buffer full"))
				}
			}
			g.mu.RUnlock()
		case <-g.stop:
			g.mu.Lock()
			for _, cli := range g.clients {
				cli.closeWithReason(nil)
			}
			g.clients = map[string]*Client{}
			g.mu.Unlock()
			return
		}
	}
}

// Stop 停止该分组（Run 循环退出，并关闭所有连接）
func (g *GroupMag) Stop() { close(g.stop) }

// Serve：升级连接并注册到当前分组，绑定 per-connection Callback。
// 注意：该方法会阻塞直到 ReadPump 退出（即连接关闭）。
func (g *GroupMag) Serve(w http.ResponseWriter, r *http.Request, key string, cb func(*Client, []byte)) error {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	c := &Client{Key: key, Conn: conn, Send: make(chan []byte, sendBufSize), Callback: cb, group: g}
	g.register <- c
	go c.WritePump()
	c.ReadPump()
	return nil
}

// SendJSON 向指定 key 的连接发送 JSON 数据
func (g *GroupMag) SendJSON(key string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	g.mu.RLock()
	cli, ok := g.clients[key]
	g.mu.RUnlock()
	if !ok {
		return errors.New("client not found")
	}
	select {
	case cli.Send <- b:
		return nil
	default:
		go cli.closeWithReason(errors.New("send buffer full"))
		return errors.New("client busy")
	}
}

// BroadcastJSON 向分组内所有连接广播 JSON 数据
func (g *GroupMag) BroadcastJSON(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	select {
	case g.broadcast <- b:
		return nil
	default:
		return errors.New("broadcast buffer full")
	}
}

// SendText 向指定 key 的连接发送文本消息
func (g *GroupMag) SendText(key, text string) error {
	g.mu.RLock()
	cli, ok := g.clients[key]
	g.mu.RUnlock()
	if !ok {
		return errors.New("client not found")
	}
	select {
	case cli.Send <- []byte(text):
		return nil
	default:
		go cli.closeWithReason(errors.New("send buffer full"))
		return errors.New("client busy")
	}
}

// BroadcastText 向分组内所有连接广播文本消息
func (g *GroupMag) BroadcastText(text string) error {
	select {
	case g.broadcast <- []byte(text):
		return nil
	default:
		return errors.New("broadcast buffer full")
	}
}

// OnlineCount 返回分组内在线连接数量
func (g *GroupMag) OnlineCount() int { g.mu.RLock(); defer g.mu.RUnlock(); return len(g.clients) }

// OnlineKeys 返回分组内在线连接的 key 列表
func (g *GroupMag) OnlineKeys() []string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	ks := make([]string, 0, len(g.clients))
	for k := range g.clients {
		ks = append(ks, k)
	}
	return ks
}

// Has 判断某个 key 的连接是否在线
func (g *GroupMag) Has(key string) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	_, ok := g.clients[key]
	return ok
}

// Disconnect 主动断开指定 key 的连接
func (g *GroupMag) Disconnect(key string) error {
	g.mu.RLock()
	cli, ok := g.clients[key]
	g.mu.RUnlock()
	if !ok {
		return errors.New("client not found")
	}
	cli.closeWithReason(nil)
	return nil
}

// -------- MainMag（主管理器，仅管理分组） --------

// MainMag 只负责分组（GroupMag）的创建、查询、删除，不直接管理连接。
type MainMag struct {
	groups map[string]*GroupMag
	mu     sync.RWMutex
}

// 默认导出：全局主管理器实例
var MainHub *MainMag

// NewMainMag 创建一个新的主管理器
func NewMainMag() *MainMag {
	MainHub = &MainMag{groups: make(map[string]*GroupMag)}
	return MainHub
}

// Get 获取（或延迟创建并启动）指定名称的分组管理器。
// name 为空时使用 "DEFAULT" 分组。
func (m *MainMag) Get(name string) *GroupMag {
	if name == "" {
		name = "DEFAULT"
	}
	m.mu.RLock()
	g, ok := m.groups[name]
	m.mu.RUnlock()
	if ok {
		return g
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if g2, ok2 := m.groups[name]; ok2 {
		return g2
	}
	gm := NewGroup(name)
	m.groups[name] = gm
	go gm.Run()
	return gm
}

// Groups 返回当前已存在的分组名称列表
func (m *MainMag) Groups() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]string, 0, len(m.groups))
	for k := range m.groups {
		out = append(out, k)
	}
	return out
}

// Delete 删除并停止指定分组（所有连接将被关闭）
func (m *MainMag) Delete(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if g, ok := m.groups[name]; ok {
		g.Stop()
		delete(m.groups, name)
	}
}

// 工具方法
// SetCheckOrigin 自定义升级器的跨域校验逻辑
func SetCheckOrigin(fn func(*http.Request) bool) { Upgrader.CheckOrigin = fn }

// closeSafe 在可能多次关闭的场景下安全关闭通道
func closeSafe(ch chan []byte) { defer func() { _ = recover() }(); close(ch) }
