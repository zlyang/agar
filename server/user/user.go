package user

import (
  "github.com/gorilla/websocket"
  "log"
  "net/http"
  "time"

  "github.com/zlyang/agar/server/logic"
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
  LogicOb logic.Logic
  Update  bool // 标识是否需要推送
}

func (c *User) readPump() {
  defer func() {
    h.unregister <- c
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
    h.broadcast <- message
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
    case message, ok := <-c.send:
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

  c := &User{Ws: ws}
  h.register <- c

  go c.writePump()
  c.readPump()
}

func NewUser(ws *websocket.Conn) *User {
}
