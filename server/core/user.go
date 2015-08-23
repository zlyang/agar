package core

import (
  "crypto/sha1"
  "encoding/hex"
  "encoding/json"
  "net/http"
  "time"

  "github.com/busyStone/agar/conn"
)

type User struct { // 以map[string]user的形式保存用户信息
  LogicOb *conn.Logic
  Update  bool // 标识是否需要推送
  ID      string
  Send    chan []byte
  Finish  chan interface{} // 无缓冲通道，用于通知操作完成
}

func ServeConnect(w http.ResponseWriter, r *http.Request) {
  wsConn := conn.DefaultConn

  err := wsConn.Serve(w, r)
  if err != nil {
    return
  }

  u, err := NewUser()
  if err != nil {
    http.Error(w, "Server not enough source", 503)
    return
  }

  H.Register <- u

  go wsConn.WritePump(u.Send)

  <-u.Finish // 等待添加进用户集合中再进行操作

  // 发送给用户更新信息
  SendSelfInfo(u)
  SendAllClientsInfo(u)

  wsConn.ReadPump(func(message []byte) {
    var m conn.C2SAction
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

  return &User{LogicOb: l, Update: true, Send: make(chan []byte, 100*1024), ID: id, Finish: make(chan interface{})}, nil
}
