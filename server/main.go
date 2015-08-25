package main

import (
  "flag"
  "log"
  "runtime"

  "github.com/busyStone/agar/conn"
  "github.com/busyStone/agar/server/core"
  "github.com/henrylee2cn/teleport"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
  log.SetFlags(log.Llongfile | log.LstdFlags | log.Lmicroseconds)

  runtime.GOMAXPROCS(runtime.NumCPU())

  core.H.Run()
  core.UpdateClientsRun()
  core.HandleLogicRun()

  conn.SocketServerInstance = conn.NewServer(teleport.API{conn.CDConnectType: new(ServerConnect),
    conn.CUMoveType: new(ClientMove)})

  // http.HandleFunc("/connect", core.ServeConnect)

  // err := http.ListenAndServe(":8080", nil)
  // if err != nil {
  // log.Fatal("ListenAndServe: ", err)
  // }
  select {}
}

type ServerConnect struct{}

func (*ServerConnect) Process(receive *teleport.NetData) *teleport.NetData {
  return core.ServeConnect(receive)
}

type ClientMove struct{}

func (*ClientMove) Process(receive *teleport.NetData) *teleport.NetData {
  return core.UserMove(receive)
}
