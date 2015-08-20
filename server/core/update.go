package core

import (
  "encoding/json"
  "time"
)

const (
  UpdateClientsPeriod = (1000 / 22) * time.Millisecond // 22Hz
)

// 周期性更新所有客户端的状态
func UpdateClientsRun() {
  go func() {
    for {
      time.Sleep(UpdateClientsPeriod)

      clients := S2CClientInfo{Type: CDUpdateClientsType, Clients: make([]Logic, 0)}
      for n, u := range H.Users {
        if u.Update == true {
          clients.Clients = append(clients.Clients, *u.LogicOb)
          H.Users[n].Update = false
        }
      }

      if len(clients.Clients) == 0 {
        continue
      }

      sendClients, err := json.Marshal(clients)
      if err != nil {
        continue
      }

      H.Broadcast <- sendClients
    }
  }()
}

func DeleteClient(u *User) {
  client := S2CDeleteClient{Type: CDDeleteClientType, Name: u.LogicOb.Name}

  deleteClient, err := json.Marshal(client)
  if err != nil {
    continue
  }

  H.Broadcast <- deleteClient
}

func SendSelfInfo(u *User) {
  client := S2CSelfInfo{Type: CDSelfClientType, ID: u.ID, Clients: u.LogicOb}

  selfClient, err := json.Marshal(client)
  if err != nil {
    continue
  }

  H.Broadcast <- selfClient
}

func SendAllClientsInfo() {
  clients := S2CClientInfo{Type: CDUpdateClientsType, Clients: make([]Logic, 0)}
  for n, u := range H.Users {
    clients.Clients = append(clients.Clients, *u.LogicOb)
  }

  sendClients, err := json.Marshal(clients)
  if err != nil {
    continue
  }

  H.Broadcast <- sendClients
}
