package main

import (
  "log"

  "github.com/busyStone/agar/conn"
  "github.com/golang/protobuf/proto"
  "github.com/henrylee2cn/teleport"
)

func main() {
  log.SetFlags(log.Llongfile | log.LstdFlags | log.Lmicroseconds)

  tp := conn.NewClient("127.0.0.1", teleport.API{conn.CDSelfClientType: new(Self),
    "all":    new(All),
    "update": new(Update)})
  tp.Request("", conn.CDConnectType)

  select {}
}

type All struct {
}

func (*All) Process(receive *teleport.NetData) *teleport.NetData { return nil }

type Update struct {
}

func (*Update) Process(receive *teleport.NetData) *teleport.NetData { return nil }

type Self struct{}

func (*Self) Process(receive *teleport.NetData) *teleport.NetData {
  var x conn.S2CSelfInfo
  err := proto.Unmarshal(receive.Body.([]byte), &x)
  if err != nil {
    log.Fatalln("fail", err)
  }
  log.Println(x.GetCanvasWidth())
  log.Fatalln("success", x)

  return nil
}
