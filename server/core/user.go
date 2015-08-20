package core

import (
  "crypto/sha1"
  "encoding/hex"
  "encoding/json"
  "net/http"
  "time"

  "github.com/gorilla/websocket"
)

type User struct { // 以map[string]user的形式保存用户信息
  LogicOb *Logic
  Update  bool // 标识是否需要推送
  ID      string
  Send    chan []byte
  Finish  chan interface{} // 无缓冲通道，用于通知操作完成
}

func ServeConnect(w http.ResponseWriter, r *http.Request) {
  conn := Conn{
    WriteWait:       10 * time.Second,
    PongWait:        60 * time.Second,
    PingPeriod:      (60 * time.Second * 9) / 10,
    MaxMessageSize:  512,
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
  }

  err := conn.Init(w, r)
  if err != nil {
    return
  }

  u, err := NewUser()
  if err != nil {
    http.Error(w, "Server not enough source", 503)
    return
  }

  H.Register <- u

  go conn.WritePump(u.Send)

  <-u.Finish // 等待添加进用户集合中再进行操作

  // 发送给用户更新信息
  SendSelfInfo(u)
  SendAllClientsInfo(u)

  conn.ReadPump(func(message []byte) {
    var m C2SAction
    if err := json.Unmarshal(message, &m); err != nil {
      return
    }

    HandleLogicChan <- ActionHandleLog{ID: m.ID, Action: m.Action} // 处理动作
  }, func() {
    H.Unregister <- u
  })

  return
}

func NewUser() (*User, error) {
  l, err := NewLogicObject()
  if err != nil {
    return nil, err
  }

  h := sha1.New()
  h.Write([]byte(l.Name + time.Now().String()))
  id := hex.EncodeToString(h.Sum(nil))

  return &User{LogicOb: l, Update: true, Send: make(chan []byte, 256), ID: id, Finish: make(chan interface{})}, nil
}
