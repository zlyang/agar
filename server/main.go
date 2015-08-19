package main

import (
  "log"
  "net/http"

  "github.com/zlyang/agar/server/hub"
  "github.com/zlyang/agar/server/user"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
  go hub.H.run()

  http.HandleFunc("/connect", user.ServeConnect)

  err := http.ListenAndServe(*addr, nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
