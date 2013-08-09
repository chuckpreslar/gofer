package tasks

import (
  "fmt"
  "github.com/chuckpreslar/gofer"
)

var Dummy = gofer.Register(gofer.Task{
  Section: "dummy",
  Label:   "task",
  //Dependencies: []string{"other:task"},
  Description: "This is a dummy task that prints a message when executed.",
  Action: func() error {
    fmt.Println("dummy:task was executed.")
    return nil
  },
})

var Another = gofer.Register(gofer.Task{
  Section:      "dummy",
  Label:        "another",
  Dependencies: []string{"dummy:fake"},
  Description:  "This is a dummy task that prints a message when executed.",
  Action: func() error {
    fmt.Println("dummy:another was executed.")
    return nil
  },
})
