package core

import (
  "errors"
  "log"
  "net/http"
  "time"

  "github.com/gorilla/websocket"
)

type Conn struct {
  Ws              *websocket.Conn
  WriteWait       time.Duration
  PongWait        time.Duration
  PingPeriod      time.Duration
  MaxMessageSize  int64
  ReadBufferSize  int
  WriteBufferSize int
}

func (c *Conn) ReadPump(handle func(m []byte), exception func()) {
  defer func() {
    exception()
    c.Ws.Close()
  }()

  c.Ws.SetReadLimit(c.MaxMessageSize)
  c.Ws.SetReadDeadline(time.Now().Add(c.PongWait))
  c.Ws.SetPongHandler(func(string) error { c.Ws.SetReadDeadline(time.Now().Add(c.PongWait)); return nil })
  for {
    _, message, err := c.Ws.ReadMessage()
    if err != nil {
      break
    }

    handle(message)
  }
}

func (c *Conn) write(mt int, payload []byte) error {
  c.Ws.SetWriteDeadline(time.Now().Add(c.WriteWait))
  return c.Ws.WriteMessage(mt, payload)
}

func (c *Conn) WritePump(send chan []byte) {
  ticker := time.NewTicker(c.PingPeriod)
  defer func() {
    ticker.Stop()
    c.Ws.Close()
  }()

  for {
    select {
    case message, ok := <-send:
      if !ok { // 如果send channel出错，则关闭连接
        c.write(websocket.CloseMessage, []byte{})
        return
      }
      if err := c.write(websocket.TextMessage, message); err != nil {
        return
      }
    case <-ticker.C: // 超时，发送Ping包
      if err := c.write(websocket.PingMessage, []byte{}); err != nil {
        return
      }
    }
  }
}

func (c *Conn) Init(w http.ResponseWriter, r *http.Request) error {
  if r.Method != "GET" {
    http.Error(w, "Method not allowed", 405)
    return errors.New("Method not allowed")
  }

  var err error
  c.Ws, err = upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Println(err)
    return err
  }

  return nil
}
