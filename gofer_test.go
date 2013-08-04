package gofer

import (
  "fmt"
  "testing"
)

var (
  namespace = "test"
  register  = "register"
  chaining  = "chaining"
  preform   = "preform"
)

func TestRegister(t *testing.T) {

  Register(Task{
    Namespace: namespace,
    Label:     register,
    Action: func() error {
      fmt.Println("Register.")
      return nil
    },
  }).Register(Task{
    Namespace: namespace,
    Label:     chaining,
    Action: func() error {
      fmt.Println("Chain register.")
      return nil
    },
  })

  if namespaceCount := len(gofer); 1 != namespaceCount {
    t.Errorf("Expected 1 namspace to be registered, got %d.", namespaceCount)
  } else if taskCount := len(gofer[namespace]); 2 != taskCount {
    t.Errorf("Expected 2 tasks to be registered, got %d.", taskCount)
  }
}

func TestPreform(t *testing.T) {
  unpreformed := true

  Register(Task{
    Namespace: namespace,
    Label:     preform,
    Action: func() error {
      unpreformed = !unpreformed
      return nil
    },
  })

  Preform(fmt.Sprintf(`%s%s%s`, namespace, DEFINITION_SPLITTER, preform))

  if unpreformed {
    t.Error("Preform failed to flip boolean test flag.")
  }
}
