package core

import (
  "log"
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
      case u := <-h.Register:
        h.Users[u.ID] = u
        u.Finish <- ""
      case u := <-h.Unregister:
        if _, ok := h.Users[u.ID]; ok {
          go DeleteClient(u)
          delete(h.Users, u.ID)
          close(u.Send)
          log.Print(len(h.Users), h.Users)
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
