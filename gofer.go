package gofer

import (
  "errors"
  "fmt"
  "os"
  "strings"
)

// DEFINITION_SPLITTER splits an string into it's namespace and label.
const DEFINITION_SPLITTER = ":"

// Gofer error definitions.
var (
  ErrBadDef      = errors.New("Bad task definition, expected <namespace>:<task>") // The task definition was malformed.
  ErrNoNamesapce = errors.New("Namespace is undefnied")                           // The namespace was undefined.
  ErrNoTask      = errors.New("Task is undefined for namespace")                  // The task was undefined.
)

// Action is a function that returns an error.
type Action func() error

// Delegator is a used to map namespaces to tasks.
// delegator[<namespace>][<task>] => Action
type Delegator map[string]map[string]Action

// Task defines an action to preform based on a namespace and label.
type Task struct {
  Namespace string // The namespace the task lives under.
  Label     string // The label the task is lives under.
  Action    Action // The action preformed by executing the task.
}

var gofer = make(Delegator)

// Register accepts a Task, storing it for later usage
// in an unexported Delegator.
func Register(task Task) Delegator {
  var (
    namespace = task.Namespace
    label     = task.Label
    action    = task.Action
  )

  if _, ok := gofer[namespace]; !ok {
    gofer[namespace] = make(map[string]Action)
  }

  gofer[namespace][label] = action

  return gofer
}

// A Delegators Register is just calls Register, allows
// for chaining.
func (self Delegator) Register(task Task) Delegator {
  return Register(task)
}

// Preform looks at the first string in the `arguments` given,
// splitting it based on the constant `DEFINITION_SPLITTER`,
// calling the action defined.
func Preform(arguments ...string) {
  definition := arguments[0]
  split := strings.Split(definition, DEFINITION_SPLITTER)

  if 2 != len(split) {
    fmt.Fprintf(os.Stderr, "%s\n", ErrBadDef)
    os.Exit(0)
  }

  var namespace map[string]Action

  if n, ok := gofer[split[0]]; ok {
    namespace = n
  } else {
    fmt.Fprintf(os.Stderr, "%s\n", ErrNoNamesapce)
    os.Exit(0)
  }

  var action Action

  if a, ok := namespace[split[1]]; ok {
    action = a
  } else {
    fmt.Fprintf(os.Stderr, "%s\n", ErrNoTask)
    os.Exit(0)
  }

  err := action()

  if nil != err {
    panic(err)
  }
}
