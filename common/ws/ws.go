package ws

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

// websocket管理器
type WsManager struct {
	Clients    map[string]*Client
	Broadcast  chan []byte
	Register   chan *Client //链接
	Unregister chan *Client //断开链接
}

// websocket 链接对象
type Client struct {
	key      string
	Socket   *websocket.Conn
	Send     chan []byte
	OutClose chan bool
	Lock     sync.Mutex
	Status   bool
}

// 消息结构提
type Message struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Content interface{} `json:"data,omitempty"`
}

var WsMag WsManager

func NewMag() *WsManager {
	WsMag = WsManager{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
	return &WsMag
}

// 注册ws链接
func (wsmag *WsManager) Reg(uuid string, conn *websocket.Conn) *Client {

	if _, ok := wsmag.Clients[uuid]; ok {
		wsmag.Clients[uuid].Close()
	}
	var c = Client{
		key:      uuid,
		Socket:   conn,
		Send:     make(chan []byte, 1024),
		OutClose: make(chan bool),
		Status:   true,
	}

	wsmag.Clients[uuid] = &c
	go c.Write()
	return &c
}

// var Manager = ClientManager{
// 	Broadcast:  make(chan []byte),
// 	// Register:   make(chan *Client),
// 	// Unregister: make(chan *Client),
// 	// Clients:    make(map[string]*Client),
// }

func (wsmag *WsManager) SetUp() {
	for {
		select {
		case conn := <-wsmag.Register:
			//加入链接管理器
			wsmag.Clients[conn.key] = conn
		case conn := <-wsmag.Unregister:
			if _, ok := wsmag.Clients[conn.key]; ok {
				conn.Close()
			}
			//广播
		case message := <-wsmag.Broadcast:
			for _, conn := range wsmag.Clients {
				conn.Send <- message
			}
		}
	}
}

func (c *Client) Read(ch chan []byte) {
	defer func() {
		WsMag.Unregister <- c
		c.Socket.Close()
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			c.Close()
			break
		}
		if string(message) == "ping" {
			if c.Status {
				c.Send <- message
			}
		}

		ch <- message
	}
}

func (c *Client) Write() {
	defer func() {
		c.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				return
			}
			if c.Status {
				c.Socket.WriteMessage(websocket.TextMessage, message)
			}

		}
	}
}

func (c *Client) Close() {
	c.Lock.Lock()
	if _, ok := WsMag.Clients[c.key]; ok {
		c.OutClose <- true
		c.Status = false
		c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
		//回收通道
		close(c.Send)
		close(c.OutClose)
		//从管理器中删除
		delete(WsMag.Clients, c.key)
	}
	c.Lock.Unlock()
}

type SendMessage struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (wsmag *WsManager) SendMain(code int, data interface{}) {
	var msg = SendMessage{
		Code: code,
		Msg:  "ok",
		Data: data,
	}
	var msgbyte, _ = json.Marshal(msg)

	for k, v := range wsmag.Clients {
		if len(k) > 4 && k[0:4] == "MAIN" {
			v.Send <- msgbyte
		}
	}
}
