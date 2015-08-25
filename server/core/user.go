package core

import (
  "log"

  "github.com/busyStone/agar/conn"
  "github.com/golang/protobuf/proto"
  "github.com/henrylee2cn/teleport"
)

type User struct { // 以map[string]user的形式保存用户信息
  LogicOb *conn.Logic
  Update  bool // 标识是否需要推送
  ID      string
  Finish  chan interface{} // 无缓冲通道，用于通知操作完成
}

func UserMove(receive *teleport.NetData) *teleport.NetData {
  var m conn.C2SAction
  err := proto.Unmarshal(receive.Body.([]byte), &m)
  if err != nil {
    log.Println(err)
    return nil
  }

  HandleLogicChan <- ActionHandleLog{ID: receive.From, Action: m.GetAction()} // 处理动作

  return nil
}

func ServeConnect(receive *teleport.NetData) *teleport.NetData {
  u, err := NewUser()
  if err != nil {
    return nil
  }
  u.ID = receive.From

  H.Register <- u
  <-u.Finish // 等待添加进用户集合中再进行操作

  // 发送给用户更新信息
  SendSelfInfo(u)
  SendAllClientsInfo(u)

  return nil
}

func NewUser() (*User, error) {
  l, err := NewLogicObject()
  if err != nil {
    return nil, err
  }

  return &User{LogicOb: l, Update: true, Finish: make(chan interface{})}, nil
}
