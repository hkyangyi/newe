# WebSocket 使用说明

## 初始化

```go
import "github.com/hkyangyi/newe/common/ws"

// 推荐：在启动阶段初始化一次（设置回调）
func init() {
  if ws.Hub == nil {
    ws.InitOnce()
  }
  ws.Hub.SetOptions(ws.Options{
    OnMessage: func(c *ws.Client, msg []byte) {
      // TODO: 业务处理
    },
    OnRegister: func(c *ws.Client) {},
    OnClose: func(c *ws.Client, err error) {},
  })
  // 启用简单协议：支持 ping/msg/who
  ws.Hub.EnableSimpleProtocol()
}
```

## 路由注册

```go
// 在你的路由汇总处
wsGroup := r.Group("/api/ws")
v2.RegisterWsRoutes(wsGroup) // GET /api/ws/:key
```

## 客户端连接

- 连接地址：`ws://<host>/api/ws/<key>`，其中 `<key>` 由业务自定义（如 `MAIN:admin` 或 `USER:123`）。
- 心跳：服务端定期 Ping；客户端可发送文本 `ping` 获得 `pong`。

## 发送消息

```go
// 单发
_ = ws.Hub.SendJSON("USER:123", ws.Envelope{Code:0, Msg:"ok", Data: map[string]any{"x":1}})

// 广播
_ = ws.Hub.BroadcastJSON(ws.Envelope{Code:0, Msg:"broadcast"})

// 按前缀群发
_ = ws.Hub.BroadcastToPrefix("MAIN", ws.Envelope{Code:0, Msg:"hi"})
```

## 前端示例

```js
const key = 'MAIN:admin';
const sock = new WebSocket(`ws://${location.host}/api/ws/${key}`);

sock.onopen = () => {
  console.log('connected');
  sock.send('ping');
};

sock.onmessage = (ev) => {
  try { const msg = JSON.parse(ev.data); console.log(msg); } catch { console.log(ev.data); }
};

sock.onclose = () => console.log('closed');
```

## 注意事项
- 同一 `key` 二次连接会自动替换旧连接，避免僵尸会话。
- 当客户端发送队列满时，服务端会断开该连接，防止阻塞。
- 如需严格校验来源，请实现 `Upgrader.CheckOrigin`。

## 简单协议（可选）

启用后支持以下下行/上行交互：

- 心跳：
  - 上行：`{"type":"ping"}`
  - 下行：`{"type":"pong","ts":1690000000}`
- 查询自身（who）：
  - 上行：`{"type":"who"}`
  - 下行：`{"type":"who","key":"MAIN:uuid","ts":...}`
- 业务消息（msg）：
  - 上行：`{"type":"msg","to":"USER:1","data":{...}}`（带 to 单发；不带 to 广播）
  - 下行：`{"type":"msg","from":"MAIN:uuid","data":{...},"ts":...}`
  - 发送方会收到 `{"type":"ack","ts":...}`

## 常用便捷方法

```go
// 在线数量 / 列表
cnt := ws.Hub.OnlineCount()
keys := ws.Hub.OnlineKeys()

// 判断在线、断开
_ = ws.Hub.Disconnect("USER:1")

// 文本发送
_ = ws.Hub.SendText("USER:1", "hello")
_ = ws.Hub.BroadcastText("system maintenance")
```
