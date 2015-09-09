package core

import (
  "github.com/busyStone/agar/conn"
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

ddd

func (h *hub) Run() {
  go func() {
    for {
      select {
      case u := <-h.Register:
        h.Users[u.ID] = u
        u.Finish <- ""
      case u := <-h.Unregister:
        if _, ok := h.Users[u.ID]; ok {
          DeleteClient(u)
          delete(h.Users, u.ID)
        }
      case m := <-h.Broadcast:
        for _, u := range h.Users {
          conn.SocketServerInstance.Request(m, "update", u.ID)
        }
      }
    }
  }()
}
