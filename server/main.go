package main

import (
  "flag"
  "log"
  "net/http"

  "github.com/zlyang/agar/server/core"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
  core.H.Run()
  core.UpdateClientsRun()
  core.HandleLogicRun()

  http.HandleFunc("/connect", core.ServeConnect)

  err := http.ListenAndServe(*addr, nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
