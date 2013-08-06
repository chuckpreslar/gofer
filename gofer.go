package gofer

import (
  "errors"
  "fmt"
  "io"
  "os"
  "path"
  "runtime"
  "strings"
)

// Gofer package constants.
const (
  DELIMITER           = ":"     // Delimeter used to split string into namespace and task.
  VERSION             = "0.0.1" // Gofer version constant.
  PROTECTED_NAMESPACE = "tasks" // No namespace should be called `tasks`

  LEFT_ALIGN = iota // Alignment for printing columns
  RIGHT_ALIGN
)

// Gofer error definitions.
var (
  ErrBadDef             = errors.New("Bad task definition.")                  // The task definition was malformed.
  ErrNoNamesapce        = errors.New("Namespace is undefnied")                // The namespace was undefined.
  ErrProtectedNamesapce = errors.New("Namespace is protected")                // The namespace used is `tasks`.
  ErrNoTask             = errors.New("Task is undefined for namespace")       // The task was undefined.
  ErrDefinedTask        = errors.New("Task is already defined for namespace") // The task was previously defined.
)

// Useful strings
var (
  GoPath        = os.Getenv("GOPATH")                // The GOPATH environtment variable.
  SourcePrefix  = "/src/"                            // Go packages live in a `src` directory.
  LabelColumn   = "            "                     // Label columns have a max length of 12.
  PackageColumn = "                                " // Package columns have a max length of 32.
)

// Action is a function that returns an error.
type Action func(arguments ...string) error

// Delegator is a used to map namespaces to tasks.
// delegator[<namespace>][<task>] => Action
type Delegator map[string]map[string]Task

// Task defines an action to preform based on a namespace and label.
type Task struct {
  Namespace    string   // The namespace the task lives under.
  Label        string   // The label the task is lives under.
  Description  string   // The description of the task's action.
  Dependencies []string // Slice of tasks that will be preformed prior to this tasks execution.
  Action       Action   // The action preformed by executing the task.
  pkg          string
}

var gofer = make(Delegator)

// Register accepts a Task, storing it for later usage
// in an unexported Delegator.
func Register(task Task) error {
  var (
    namespace = task.Namespace
    label     = task.Label
  )

  // A package attempted to register the PROTECTED_NAMESPACE of `tasks`
  if PROTECTED_NAMESPACE == namespace {
    return ErrProtectedNamesapce
  }

  // Attach to directory of the caller to the package.
  if _, file, _, ok := runtime.Caller(1); ok {
    task.pkg = path.Dir(strings.TrimLeft(file, path.Join(GoPath, SourcePrefix)))
  }

  // If this is a new namespace, initialize it on the map.
  if _, ok := gofer[namespace]; !ok {
    gofer[namespace] = make(map[string]Task)
  }

  // The task has already been defined for the namespace.
  if _, ok := gofer[namespace][label]; ok {
    return ErrDefinedTask
  }

  gofer[namespace][label] = task

  return nil
}

// Preform looks at the first string in the `arguments` given,
// splitting it based on the constant `DELIMITER`,
// calling the action defined.
func Preform(arguments ...string) error {
  definition := arguments[0]
  split := strings.Split(definition, DELIMITER) // Split the difinition by the DELIMITER constant.

  // Expect the split to result in an array with length of 2.
  if 2 != len(split) {
    return ErrBadDef
  }

  // Remove the `definition` argument from the array,
  // the remainder will be passed to the action.
  arguments = arguments[1:]

  var namespace map[string]Task

  //  Ensure the namespace exists.
  if n, ok := gofer[split[0]]; ok {
    namespace = n
  } else {
    return ErrNoNamesapce
  }

  task, ok := namespace[split[1]]

  // Ensure the task lives in the namespace.
  if !ok {
    return ErrNoNamesapce
  }

  // If the task has any dependencies, preform them first.
  for _, dep := range task.Dependencies {
    if err := Preform(append([]string{dep}, arguments...)...); nil != err {
      return err
    }
  }

  return task.Action(arguments...)
}

// ListTasks prints a lists of the registered tasks to the writter.
func ListTasks(writter io.Writer, arguments ...string) error {
  // Generated template from binary provides `tasks` from the CLI
  // as the first argument.
  if PROTECTED_NAMESPACE == arguments[0] {
    arguments = arguments[1:]
  }

  printTaskListHeader(writter)

  // Assume all tasks are to be printed.
  if 0 == len(arguments) {
    return ListAllTasks(writter)
  }

  // Assume an only specific tasks are to be printed.
  for _, namespace := range arguments {
    if err := ListTasksFor(writter, namespace); nil != err {
      return err
    }
  }

  return nil
}

// ListAllTasks prints a list of all tasks to the writter.
func ListAllTasks(writter io.Writer) error {
  for namespace, _ := range gofer {
    if err := ListTasksFor(writter, namespace); nil != err {
      return err
    }
  }

  return nil
}

// ListTasksFor prints a tasks for a specific namespace to the writter.
func ListTasksFor(writter io.Writer, namespace string) error {
  if n, ok := gofer[namespace]; ok {
    index := 0

    for _, t := range n {
      var namespaceColumn string

      // If this is the first iteration through the delegator for the task,
      // print the namespace in the first column.
      if index == 0 {
        namespaceColumn = generateColumn(t.Namespace, LabelColumn, RIGHT_ALIGN)
      } else {
        namespaceColumn = LabelColumn
      }

      // Generate the other columns to be printed.
      labelColumn := generateColumn(t.Label, LabelColumn, LEFT_ALIGN)
      pkgColumn := generateColumn(t.pkg, PackageColumn, LEFT_ALIGN)
      descriptionColumn := t.Description

      fmt.Fprintf(writter, "%s %s %s %s\n", namespaceColumn, labelColumn, pkgColumn, descriptionColumn)
      index += 1
    }
  } else {
    return ErrNoNamesapce
  }

  return nil
}

func printTaskListHeader(writter io.Writer) {
  namespaceColumn := generateColumn("Namespace", LabelColumn, RIGHT_ALIGN)
  labelColumn := generateColumn("Task", LabelColumn, LEFT_ALIGN)
  pkgColumn := generateColumn("Directory", PackageColumn, LEFT_ALIGN)
  descriptionColumn := "Description"
  fmt.Fprintf(writter, "%s %s %s %s\n", namespaceColumn, labelColumn, pkgColumn, descriptionColumn)
}

// generateColumn is a helper for aligning output.
func generateColumn(data, column string, alignment int) (result string) {
  // If the data causes the column to overflow, truncate it.
  if len(column) < len(data) {
    data = fmt.Sprintf("...%s", data[len(data)-len(column)+3:])
  }

  if alignment == RIGHT_ALIGN {
    result = fmt.Sprintf("%s%s", column[len(data):], data)
  } else if alignment == LEFT_ALIGN {
    result = fmt.Sprintf("%s%s", data, column[len(data):])
  }

  return
}
