package core

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
      case u := <-h.Unregister:
        if _, ok := h.Users[u.LogicOb.Name]; ok {
          DeleteClient(u)
          delete(h.Users, u.LogicOb.Name)
          close(u.Send)
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
