package main

import (
  // "encoding/json"
  "log"
  "net/http"
  "testing"
  "time"

  "github.com/gorilla/websocket"
  // "github.com/zlyang/agar/server/core"
)

func TestConnect(t *testing.T) {
  conn, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:8080/connect", http.Header{})
  if err != nil {
    t.Fatal(err)
  }

  messageType, content, err := conn.ReadMessage()
  if err != nil {
    t.Fatal(err)
  }

  log.Print(messageType, string(content))

  time.Sleep(time.Minute * 2)

  conn.Close()
}
