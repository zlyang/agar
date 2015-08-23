package main

import (
  "encoding/json"
  "flag"
  "log"
  "math/rand"
  "time"

  "github.com/busyStone/agar/conn"
)

const (
  AutoRunAction = "UDRL"
  RandInterval  = 10
  IntervalBase  = 10
)

var maxIntervalTime time.Duration

type AutoRunControl struct {
  last time.Time
  ID   string
}

func main() {
  var ip = flag.String("ip", "127.0.0.1", "local or other") // 主机ip
  var concurrent = flag.Int("c", 10, "concurrent num")      // 并发数

  flag.Parse()

  log.Print(*ip, *concurrent)

  for i := 0; i < *concurrent; i++ {
    go func() {
      wsConn := conn.DefaultConn
      err := wsConn.Client(`ws://`+(*ip)+`:8080/connect`, 1024, 1024)
      if err != nil {
        log.Println(err)
        return
      }

      var au AutoRunControl
      // 启动 读取 goroutine
      wsConnWriteChan := make(chan []byte, 256)
      go wsConn.ReadPump(au.HandleRead, nil)
      go wsConn.WritePump(wsConnWriteChan)
      go au.AutoRun(wsConnWriteChan)
    }()
  }

  <-make(chan interface{})
}

func (a *AutoRunControl) AutoRun(writeChane chan []byte) {
  r := rand.New(rand.NewSource(time.Now().UnixNano()))

  for {
    var act string
    // for i := 0; i < 100; i++ {
    act += string(AutoRunAction[r.Int31n(int32(len(AutoRunAction)))]) // 随机走一定的步数
    // }
    interval := r.Int31n(RandInterval) + IntervalBase
    // if a.ID != "" {
    action := conn.C2SAction{
      Type:   conn.CUActionType,
      ID:     a.ID,
      Action: act,
    }

    data, err := json.Marshal(action)
    if err == nil {
      writeChane <- data
    }
    // }

    time.Sleep(time.Duration(interval) * time.Millisecond)
  }
}

// MessageType 用于确定消息类型
type MessageType struct {
  Type int
}

func (a *AutoRunControl) HandleRead(m []byte) {
  if a.ID == "" {
    var mt MessageType
    if err := json.Unmarshal(m, &mt); err != nil {
      return
    }

    if mt.Type == conn.CDSelfClientType {
      var player conn.S2CSelfInfo
      if err := json.Unmarshal(m, &player); err != nil {
        return
      }
      a.ID = player.ID
    }

    a.last = time.Now()

    return
  }

  now := time.Now()

  if now.Sub(a.last) > maxIntervalTime {
    maxIntervalTime = now.Sub(a.last)

  }

  log.Println(maxIntervalTime)

  a.last = now
}
