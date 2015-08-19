package core

import (
  "github.com/gorilla/websocket"
  "log"
  "net/http"
  "strconv"
  "time"
)

const (
  writeWait      = 10 * time.Second
  pongWait       = 60 * time.Second
  pingPeriod     = (pongWait * 9) / 10
  maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
  ReadBufferSize:  1024,
  WriteBufferSize: 1024,
}

type User struct { // 以map[string]user的形式保存用户信息
  Ws      *websocket.Conn
  LogicOb *Logic
  Update  bool // 标识是否需要推送
  Send    chan []byte
}

func (c *User) readPump() {
  defer func() {
    H.Unregister <- c
    c.Ws.Close()
  }()
  c.Ws.SetReadLimit(maxMessageSize)
  c.Ws.SetReadDeadline(time.Now().Add(pongWait))
  c.Ws.SetPongHandler(func(string) error { c.Ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
  for {
    _, message, err := c.Ws.ReadMessage()
    if err != nil {
      break
    }
    // H.Broadcast <- message    // 由于传送的控制信息，需要定时处理，所以不需要广播

    // TODO: 先放在队列中，定时进行处理逻辑
  }
}

func (c *User) write(mt int, payload []byte) error {
  c.Ws.SetWriteDeadline(time.Now().Add(writeWait))
  return c.Ws.WriteMessage(mt, payload)
}

func (c *User) writePump() {
  ticker := time.NewTicker(pingPeriod)
  defer func() {
    ticker.Stop()
    c.Ws.Close()
  }()

  for {
    select {
    case message, ok := <-c.Send:
      if !ok { // 如果send channel出错，则关闭连接
        c.write(websocket.CloseMessage, []byte{})
        return
      }
      if err := c.write(websocket.TextMessage, message); err != nil {
        return
      }
    case <-ticker.C: // 超时，发送Ping包
      if err := c.write(websocket.PingMessage, []byte{}); err != nil {
        return
      }
    }
  }
}

func (u *User) GetInfoString() string {
  return u.LogicOb.Name + ",(" + strconv.Itoa(u.LogicOb.Position.X) + "," + strconv.Itoa(u.LogicOb.Position.Y) + ")"
}

func ServeConnect(w http.ResponseWriter, r *http.Request) {
  if r.Method != "GET" {
    http.Error(w, "Method not allowed", 405)
    return
  }

  ws, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Println(err)
    return
  }

  u, err := NewUser(ws)
  if err != nil {
    http.Error(w, "Server not enough source", 503)
    return
  }

  H.Register <- u

  go u.writePump()

  u.Send <- []byte(u.GetInfoString())

  u.readPump()

  return
}

func NewUser(ws *websocket.Conn) (*User, error) {
  l, err := NewLogicObject()
  if err != nil {
    return nil, err
  }

  return &User{Ws: ws, LogicOb: l, Update: false, Send: make(chan []byte, 256)}, nil
}
