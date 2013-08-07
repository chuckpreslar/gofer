package main

import (
  "github.com/chuckpreslar/gofer"

  "fmt"
  "os"
)

type command struct {
  name        string
  trigger     string
  description string
}

func main() {
  definition := os.Args[1]

  if err := gofer.LoadAndPreform(definition); nil != err {
    fmt.Fprintf(os.Stderr, "%s\n", err)
  }
}
