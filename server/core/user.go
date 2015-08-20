package core

import (
  "crypto/sha1"
  "encoding/hex"
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
  Finish  chan interface{} // 无缓冲通道，用于通知操作完成
}

func (u *User) readPump() {
  defer func() {
    H.Unregister <- u
    u.Ws.Close()
  }()
  u.Ws.SetReadLimit(maxMessageSize)
  u.Ws.SetReadDeadline(time.Now().Add(pongWait))
  u.Ws.SetPongHandler(func(string) error { u.Ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
  for {
    _, message, err := u.Ws.ReadMessage()
    if err != nil {
      break
    }

    var m C2SAction
    if err := json.Unmarshal(message, &m); err != nil {
      log.Println(string(message), err)
      continue
    }

    HandleLogicChan <- ActionHandleLog{ID: m.ID, Action: m.Action} // 处理动作
  }
}

func (u *User) write(mt int, payload []byte) error {
  u.Ws.SetWriteDeadline(time.Now().Add(writeWait))
  return u.Ws.WriteMessage(mt, payload)
}

func (u *User) writePump() {
  ticker := time.NewTicker(pingPeriod)
  defer func() {
    ticker.Stop()
    u.Ws.Close()
  }()

  for {
    select {
    case message, ok := <-u.Send:
      if !ok { // 如果send channel出错，则关闭连接
        u.write(websocket.CloseMessage, []byte{})
        return
      }
      if err := u.write(websocket.TextMessage, message); err != nil {
        return
      }
    case <-ticker.C: // 超时，发送Ping包
      if err := u.write(websocket.PingMessage, []byte{}); err != nil {
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

  <-u.Finish // 等待添加进用户集合中再进行操作

  // 发送给用户更新信息
  SendSelfInfo(u)
  SendAllClientsInfo(u)

  u.readPump()

  return
}

func NewUser(ws *websocket.Conn) (*User, error) {
  l, err := NewLogicObject()
  if err != nil {
    return nil, err
  }

  h := sha1.New()
  h.Write([]byte(l.Name + time.Now().String()))
  id := hex.EncodeToString(h.Sum(nil))

  return &User{Ws: ws, LogicOb: l, Update: true, Send: make(chan []byte, 256), ID: id, Finish: make(chan interface{})}, nil
}
