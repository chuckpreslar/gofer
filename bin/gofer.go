package main

import (
  "github.com/chuckpreslar/gofer"

  "fmt"
  "os"
)

const VERSION = "0.1.0"

func main() {
  var definition string

  if 1 >= len(os.Args) {
    definition = ""
  } else {
    definition = os.Args[1]
  }

  if "version" == definition {
    fmt.Fprintf(os.Stdout, "%s", VERSION)
    os.Exit(0)
  }

  if err := gofer.LoadAndPerform(definition); nil != err {
    fmt.Fprintf(os.Stderr, "%s\n", err)
  }
}
