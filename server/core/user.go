package core

import (
  "crypto/sha1"
  "encoding/json"
  "log"
  "net/http"
  "strconv"
  "time"

  "github.com/gorilla/websocket"
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
  ID      string
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

    var m C2SAction
    if err := json.Unmarshal(message, &m); err != nil {
      log.Println(string(message), err)
      continue
    }

    HandleLogicChan <- ActionHandleLog{Name: m.Name, Action: m.Action} // 处理动作
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

  // 组织回传的数据，返回所有的
  clients := S2CClientInfo{Type: CDAllClientsType, Clients: make([]Logic, 0)}
  clients.Clients = append(clients.Clients, *u.LogicOb)
  for _, u := range H.Users {
    clients.Clients = append(clients.Clients, *u.LogicOb)
  }
  sendClients, err := json.Marshal(clients)
  if err != nil {
    return
  }

  H.Register <- u

  go u.writePump()

  u.Send <- sendClients

  u.readPump()

  return
}

func NewUser(ws *websocket.Conn) (*User, error) {
  l, err := NewLogicObject()
  if err != nil {
    return nil, err
  }

  h := sha1.New()
  id := h.Sum([]byte(l.Name + time.Now().String()))

  return &User{Ws: ws, LogicOb: l, Update: true, Send: make(chan []byte, 256), ID: id}, nil
}
