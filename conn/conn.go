package conn

import (
  "errors"
  "fmt"
  "log"
  "net"
  "net/http"
  "net/url"
  "time"

  "github.com/gorilla/websocket"
)

type Conn struct {
  ws              *websocket.Conn
  WriteWait       time.Duration
  PongWait        time.Duration
  PingPeriod      time.Duration
  MaxMessageSize  int64
  ReadBufferSize  int
  WriteBufferSize int
}

var DefaultConn = Conn{
  WriteWait:       10 * time.Second,
  PongWait:        60 * time.Second,
  PingPeriod:      (60 * time.Second * 9) / 10,
  MaxMessageSize:  1024 * 100,
  ReadBufferSize:  1024 * 100,
  WriteBufferSize: 1024 * 100,
}

func (c *Conn) ReadPump(handle func(m []byte), exception func()) {
  defer func() {
    if exception != nil {
      exception()
    }
    c.ws.Close()
  }()

  c.ws.SetReadLimit(c.MaxMessageSize)
  c.ws.SetReadDeadline(time.Now().Add(c.PongWait))
  c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(c.PongWait)); return nil })
  for {
    _, message, err := c.ws.ReadMessage()
    if err != nil {
      break
    }

    if handle != nil {
      handle(message)
    }
  }
}

func (c *Conn) write(mt int, payload []byte) error {
  c.ws.SetWriteDeadline(time.Now().Add(c.WriteWait))
  return c.ws.WriteMessage(mt, payload)
}

func (c *Conn) WritePump(send chan []byte) {
  ticker := time.NewTicker(c.PingPeriod)
  defer func() {
    ticker.Stop()
    c.ws.Close()
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

func (c *Conn) Serve(w http.ResponseWriter, r *http.Request) error {
  if r.Method != "GET" {
    http.Error(w, "Method not allowed", 405)
    return errors.New("Method not allowed")
  }

  var upgrader = websocket.Upgrader{
    ReadBufferSize:  c.ReadBufferSize,
    WriteBufferSize: c.WriteBufferSize,
  }

  var err error
  c.ws, err = upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Println(err)
    return err
  }

  return nil
}

func (c *Conn) IsConnValid() bool {
  return c.ws != nil
}

func (c *Conn) Client(connectURL string, readBufSize, writeBufSize int) error {

  u, err := url.Parse(connectURL)
  if err != nil {
    return err
  }

  rawConn, err := net.Dial("tcp", u.Host)
  if err != nil {
    return err
  }

  wsHeaders := http.Header{
    "Origin": {u.Scheme + "://" + u.Host},
    // your milage may differ
    "Sec-WebSocket-Extensions": {"permessage-deflate; client_max_window_bits, x-webkit-deflate-frame"},
  }

  wsConn, resp, err := websocket.NewClient(rawConn, u, wsHeaders, readBufSize, writeBufSize)
  if err != nil {
    return fmt.Errorf("websocket.NewClient Error: %s\nResp:%+v", err, resp)
  }

  c.ws = wsConn

  return nil
}
