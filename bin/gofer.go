package main

import (
  "github.com/chuckpreslar/gofer"

  "fmt"
  "os"
)

func main() {
  var definition string

  if 1 >= len(os.Args) {
    definition = ""
  } else {
    definition = os.Args[1]
  }

  if err := gofer.LoadAndPreform(definition); nil != err {
    fmt.Fprintf(os.Stderr, "%s\n", err)
  }
}
