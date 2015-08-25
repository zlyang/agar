package main

import (
  "flag"
  "log"
  "math/rand"
  "time"

  "github.com/busyStone/agar/conn"
  "github.com/golang/protobuf/proto"
  "github.com/henrylee2cn/teleport"
)

const (
  AutoRunAction = "UDRL"
  RandInterval  = 10
  IntervalBase  = 10
)

var maxIntervalTime time.Duration

type AutoRunControl struct {
  tp teleport.Teleport
}

func main() {
  var ip = flag.String("ip", "127.0.0.1", "local or other") // 主机ip
  var concurrent = flag.Int("c", 10, "concurrent num")      // 并发数

  flag.Parse()

  log.Print(*ip, *concurrent)

  for i := 0; i < *concurrent; i++ {
    go func() {
      tp := conn.NewClient("127.0.0.1", teleport.API{
        conn.CDUpdateClientsType: new(ReadHandle),
        conn.CDSelfClientType:    new(Empty),
        conn.CDAllClientsType:    new(Empty),
      })

      var au AutoRunControl
      au.tp = tp

      tp.Request("", conn.CDConnectType)

      // 启动 读取 goroutine
      go au.AutoRun()
    }()

    time.Sleep(time.Millisecond * 3)
  }

  select {}
}

func (a *AutoRunControl) AutoRun() {
  r := rand.New(rand.NewSource(time.Now().UnixNano()))

  for {
    var act string
    // for i := 0; i < 100; i++ {
    act += string(AutoRunAction[r.Int31n(int32(len(AutoRunAction)))]) // 随机走一定的步数
    // }
    interval := r.Int31n(RandInterval) + IntervalBase
    // if a.ID != "" {
    action := conn.C2SAction{
      Action: proto.String(act),
    }

    data, err := proto.Marshal(&action)
    if err == nil {
      a.tp.Request(data, conn.CUMoveType)
    }
    // }

    time.Sleep(time.Duration(interval) * time.Millisecond)
  }
}

type ReadHandle struct{ last time.Time }

func (r *ReadHandle) Process(receive *teleport.NetData) *teleport.NetData {
  now := time.Now()

  // log.Println(now.Sub(r.last))

  r.last = now

  return nil
}

type Empty struct {
}

func (*Empty) Process(receive *teleport.NetData) *teleport.NetData { return nil }
