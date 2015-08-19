package core

import (
  "time"
)

const (
  UpdateClientsPeriod = (1000 / 22) * time.Millisecond // 22Hz
)

type hub struct {
  Users      map[string]*User
  Broadcast  chan []byte
  Register   chan *User
  Unregister chan *User
}

var H = hub{
  Broadcast:  make(chan []byte),
  Register:   make(chan *User),
  Unregister: make(chan *User),
  Users:      make(map[string]*User),
}

func (h *hub) Run() {
  go func() {
    for {
      select {
      case c := <-h.Register:
        h.Users[c.LogicOb.Name] = c
      case c := <-h.Unregister:
        if _, ok := h.Users[c.LogicOb.Name]; ok {
          delete(h.Users, c.LogicOb.Name)
          close(c.Send)
        }
      case m := <-h.Broadcast:
        for k, u := range h.Users {
          select {
          case u.Send <- m:
          default:
            close(u.Send)
            delete(h.Users, k)
          }
        }
      }
    }
  }()
}

// 周期性更新所有客户端的状态
func UpdateClientsRun() {
  go func() {
    for {
      time.Sleep(UpdateClientsPeriod)

      for _, u := range H.Users {

      }
    }
  }()
}
