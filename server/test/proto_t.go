package main

import (
  "encoding/base64"
  "log"

  "github.com/busyStone/agar/conn"
  "github.com/golang/protobuf/proto"
  "github.com/henrylee2cn/teleport"
)

const (
  base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

var coder = base64.NewEncoding(base64Table)

func main() {
  log.SetFlags(log.Llongfile | log.LstdFlags | log.Lmicroseconds)

  tp := conn.NewClient("127.0.0.1", teleport.API{conn.CDSelfClientType: Self,
    "all":    func(receive *teleport.NetData) *teleport.NetData { return nil },
    "update": func(receive *teleport.NetData) *teleport.NetData { return nil }})
  tp.Request("", conn.CDConnectType)

  select {}
}

func Self(receive *teleport.NetData) *teleport.NetData {
  body, ok := receive.Body.(string)
  if !ok {
    log.Fatalln("string error")
    return nil
  }

  message, err := coder.DecodeString(body)
  if err != nil {
    log.Fatalln(body, err)
    return nil
  }

  var x conn.S2CSelfInfo
  err = proto.Unmarshal(message, &x)
  if err != nil {
    log.Fatalln("fail", err)
  }

  log.Fatalln("success", x)

  return nil
}
