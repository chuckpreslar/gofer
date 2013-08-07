package gofer

import (
  "errors"
  "strings"
)

const (
  DELIMITER = ":"
)

var (
  errBadLabel            = errors.New("Bad label for task, unexpected section delimiter.")
  errRegistrationFailure = errors.New("Registration for task failed unexpectedly.")
  errUnknownTask         = errors.New("Unable to look up task.")
  errNoAction            = errors.New("No action defined for task.")
)

type action func(arguments ...interface{}) error

type Task struct {
  Section      string
  Label        string
  Description  string
  Dependencies []string
  Action       action
  manual       manual
  location     string
}

type manual []*Task

// gofer variable used for storing tasks.
var gofer = make(manual, 0)

// index searches through the manual, returning a task
// found with the label and in the section (namepsace) defined
// by the definition.
func (self *manual) index(definition string) (task *Task) {
  sections := strings.Split(definition, DELIMITER)
  entries := self

  for _, section := range sections {
    for i := 0; i < len(*entries); i++ {
      if (*entries)[i].Label == section {
        task = (*entries)[i]
        entries = &task.manual // adjust `entries` pointer for next iteration.
        break
      }
    }

    if nil == task {
      return
    } else if section != task.Label {
      return nil
    }
  }

  return
}

// sectionalize creates potential missing spaces in a manual
// based on the `definition`.
func (self *manual) sectionalize(definition string) (task *Task) {
  task = self.index(definition)

  if nil != task {
    return
  }

  sections := strings.Split(definition, DELIMITER)

  task = new(Task)
  task.Label = sections[0]

  *self = append(*self, task)

  for i := 1; i < len(sections); i++ {
    temp := new(Task)
    temp.Section = strings.Join(sections[:i], DELIMITER)
    temp.Label = sections[i]

    task.manual = append(task.manual, temp)
    task = temp // update task to temp for next iteration.
  }

  return
}

// Register accepts a `Task`, storing it for later.
func Register(task Task) (err error) {
  if index := strings.Index(task.Label, DELIMITER); -1 != index {
    return errBadLabel
  }

  parent := gofer.sectionalize(task.Section)

  if nil == parent {
    if 0 != len(task.Section) {
      return errRegistrationFailure
    }

    gofer = append(gofer, &task)
  } else {
    parent.manual = append(parent.manual, &task)
  }

  return
}

func Preform(definition string) (err error) {
  task := gofer.index(definition)

  if nil == task {
    return errUnknownTask
  } else if nil == task.Action {
    return errNoAction
  }

  for _, dependency := range task.Dependencies {
    err = Preform(dependency)

    if nil != err {
      return
    }
  }

  return task.Action()
}
