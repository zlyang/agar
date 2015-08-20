package main

import (
  "flag"
  "log"
  "net/http"
  "runtime"

  "github.com/busyStone/agar/server/core"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
  log.SetFlags(log.Lshortfile | log.LstdFlags)

  runtime.GOMAXPROCS(runtime.NumCPU())

  core.H.Run()
  core.UpdateClientsRun()
  core.HandleLogicRun()

  http.HandleFunc("/connect", core.ServeConnect)

  err := http.ListenAndServe(*addr, nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
