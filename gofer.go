package gofer

import (
  "errors"
  "fmt"
  "io"
  "strings"
  "text/template"
)

const (
  DELIMITER = ":"     // Delimeter used to split string into namespace and task.
  VERSION   = "0.0.1" // Gofer version constant.
)

// Gofer error definitions.
var (
  ErrBadDef      = errors.New("Bad task definition, expected <namespace>:<task>") // The task definition was malformed.
  ErrNoNamesapce = errors.New("Namespace is undefnied")                           // The namespace was undefined.
  ErrNoTask      = errors.New("Task is undefined for namespace")                  // The task was undefined.
)

var TaskTemplate = template.Must(template.New("task").Parse(`
  Namespace:  {{.Namespace}}
  Label:      {{.Label}}
  {{if .Description}}
  Description:
  {{.Description}}
  {{end}}
`))

// Action is a function that returns an error.
type Action func() error

// Delegator is a used to map namespaces to tasks.
// delegator[<namespace>][<task>] => Action
type Delegator map[string]map[string]Task

// Task defines an action to preform based on a namespace and label.
type Task struct {
  Namespace   string // The namespace the task lives under.
  Label       string // The label the task is lives under.
  Description string // The description of the task's action.
  Action      Action // The action preformed by executing the task.
}

var gofer = make(Delegator)

// Register accepts a Task, storing it for later usage
// in an unexported Delegator.
func Register(task Task) Delegator {
  var (
    namespace = task.Namespace
    label     = task.Label
  )

  if _, ok := gofer[namespace]; !ok {
    gofer[namespace] = make(map[string]Task)
  }

  gofer[namespace][label] = task

  return gofer
}

// A Delegators Register is just calls Register, allows
// for chaining.
func (self Delegator) Register(task Task) Delegator {
  return Register(task)
}

// Preform looks at the first string in the `arguments` given,
// splitting it based on the constant `DELIMITER`,
// calling the action defined.
func Preform(arguments ...string) error {
  definition := arguments[0]
  split := strings.Split(definition, DELIMITER)

  if 2 != len(split) {
    return ErrBadDef
  }

  var namespace map[string]Task

  if n, ok := gofer[split[0]]; ok {
    namespace = n
  } else {
    return ErrNoNamesapce
  }

  var action Action

  if task, ok := namespace[split[1]]; ok {
    action = task.Action
  } else {
    return ErrNoNamesapce
  }

  return action()
}

func ListTasks(writter io.Writer, arguments ...string) error {
  if 1 == len(arguments) {
    return ListAllTasks(writter)
  }

  for _, namespace := range arguments[1:] {
    if err := ListTasksFor(writter, namespace); nil != err {
      return err
    }
  }

  return nil
}

func ListAllTasks(writter io.Writer) error {
  for namespace, _ := range gofer {
    if err := ListTasksFor(writter, namespace); nil != err {
      return err
    }
  }

  return nil
}

func ListTasksFor(writter io.Writer, namespace string) error {
  if n, ok := gofer[namespace]; ok {
    for _, t := range n {
      err := TaskTemplate.Execute(writter, t)

      if nil != err {
        return err
      }
    }
  } else {
    return ErrNoNamesapce
  }

  return nil
}

func PrintVersion(writter io.Writer) {
  fmt.Fprintf(writter, "%s\n", VERSION)
}
