package core

import (
  "bytes"
  "compress/gzip"
  "encoding/json"
  "log"
  "time"

  "github.com/busyStone/agar/conn"
  "github.com/golang/protobuf/proto"
)

const (
  UpdateClientsPeriod = (1000 / 22) * time.Millisecond // 22Hz
)

// 周期性更新所有客户端的状态
func UpdateClientsRun() {
  go func() {
    for {
      time.Sleep(UpdateClientsPeriod)

      clients := conn.S2CClientInfo{
        Clients: make([]*conn.Logic, 0)}

      for n, u := range H.Users {
        if u.Update == true {
          clients.Clients = append(clients.Clients, u.LogicOb)
          H.Users[n].Update = false
        }
      }

      if len(clients.Clients) == 0 {
        continue
      }

      Send2Broadcast(&clients)
    }
  }()
}

func DeleteClient(u *User) {
  client := conn.S2CDeleteClient{
    Name: u.LogicOb.Name}

  go Send2Broadcast(&client)
}

func SendSelfInfo(u *User) {
  client := conn.S2CSelfInfo{
    Clients:      u.LogicOb,
    CanvasWidth:  proto.Int32(int32(CanvasWidth)),
    CanvasHeight: proto.Int32(int32(CanvasHeight))}

  Send2User(u, &client, conn.CDSelfClientType)
}

func SendAllClientsInfo(u *User) {
  clients := conn.S2CClientInfo{
    Clients: make([]*conn.Logic, 0)}

  for _, u := range H.Users {
    clients.Clients = append(clients.Clients, u.LogicOb)
  }

  Send2User(u, &clients, conn.CDAllClientsType)
}

func Send2Broadcast(s proto.Message) {
  s2cjson, err := json.Marshal(s)
  if err != nil {
    return
  }

  s2c, err := proto.Marshal(proto.Message(s))
  if err != nil {
    return
  }

  log.Println(len(s2cjson), len(s2c))

  // s2cc := Gzip(s2c)

  H.Broadcast <- s2c
}

func Send2User(u *User, s proto.Message, action string) {
  s2cjson, err := json.Marshal(s)
  if err != nil {
    return
  }

  // s2c, err := json.Marshal(s)
  s2c, err := proto.Marshal(s)
  if err != nil {
    return
  }

  log.Println(len(s2cjson), len(s2c), u.ID, action, u)

  // s2cc := Gzip(s2c)
  conn.SocketServerInstance.Request(s2c, action, u.ID)
}

func Gzip(s []byte) []byte {
  var b bytes.Buffer
  w := gzip.NewWriter(&b)
  defer w.Close()

  w.Write(s)
  w.Flush()

  return b.Bytes()
}
