package core

import (
  "encoding/json"
  "time"

  "github.com/busyStone/agar/conn"
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
        Type:    conn.CDUpdateClientsType,
        Clients: make([]conn.Logic, 0)}

      for n, u := range H.Users {
        if u.Update == true {
          clients.Clients = append(clients.Clients, *u.LogicOb)
          H.Users[n].Update = false
        }
      }

      if len(clients.Clients) == 0 {
        continue
      }

      Send2Broadcast(clients)
    }
  }()
}

func DeleteClient(u *User) {
  client := conn.S2CDeleteClient{
    Type: conn.CDDeleteClientType,
    Name: u.LogicOb.Name}

  go Send2Broadcast(client)
}

func SendSelfInfo(u *User) {
  client := conn.S2CSelfInfo{
    Type:         conn.CDSelfClientType,
    ID:           u.ID,
    Clients:      *u.LogicOb,
    CanvasWidth:  CanvasWidth,
    CanvasHeight: CanvasHeight}

  Send2User(u, client)
}

func SendAllClientsInfo(u *User) {
  clients := conn.S2CClientInfo{
    Type:    conn.CDAllClientsType,
    Clients: make([]conn.Logic, 0)}

  for _, u := range H.Users {
    clients.Clients = append(clients.Clients, *u.LogicOb)
  }

  Send2User(u, clients)
}

func Send2Broadcast(s interface{}) {
  s2c, err := json.Marshal(s)
  if err != nil {
    return
  }

  H.Broadcast <- s2c
}

func Send2User(u *User, s interface{}) {
  s2c, err := json.Marshal(s)
  if err != nil {
    return
  }

  u.Send <- s2c
}
