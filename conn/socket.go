package conn

import (
  "github.com/henrylee2cn/teleport"
)

var SocketServerInstance teleport.Teleport

type Socket struct {
}

func NewServer(apis teleport.API) teleport.Teleport {
  tp := teleport.New()
  tp.SetAPI(apis).Server(":8080")

  return tp
}

func NewClient(ip string, apis teleport.API) teleport.Teleport {
  tp := teleport.New()
  tp.SetAPI(apis).Client(ip, ":8080")

  return tp
}
