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
    Action: func(arguments ...string) error {
      fmt.Println("Register.")
      return nil
    },
  })

  if namespaceCount := len(gofer); 1 != namespaceCount {
    t.Errorf("Expected 1 namspace to be registered, got %d.", namespaceCount)
  }
}

func TestPreform(t *testing.T) {
  unpreformed := true

  Register(Task{
    Namespace: namespace,
    Label:     preform,
    Action: func(arguments ...string) error {
      unpreformed = !unpreformed
      return nil
    },
  })

  Preform(fmt.Sprintf(`%s%s%s`, namespace, DELIMITER, preform))

  if unpreformed {
    t.Error("Preform failed to flip boolean test flag.")
  }
}
