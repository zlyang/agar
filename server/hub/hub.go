package hub

import (
  "github.com/zlyang/agar/server/user"
)

type hub struct {
  Users      map[string]user.user
  Broadcast  chan []byte
  Register   chan *connection
  Unregister chan *connection
}

var H = hub{
  Broadcast:  make(chan []byte),
  Register:   make(chan *connection),
  Unregister: make(chan *connection),
  Users:      make(map[string]user.user),
}

func (h *hub) run() {
  for {
    select {
    case c := <-h.Register:
      h.Users[c] = true
    case c := <-h.Unregister:
      if _, ok := h.Users[c]; ok {
        delete(h.Users, c)
        close(c.send)
      }
    case m := <-h.Broadcast:
      for c := range h.Users {
        select {
        case c.send <- m:
        default:
          close(c.send)
          delete(h.Users, c)
        }
      }
    }
  }
}
