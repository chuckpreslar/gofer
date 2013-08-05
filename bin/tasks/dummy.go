package tasks

import (
  "fmt"
)

type Dummy struct{}

const (
  CONST_TEST = 1
)

var (
  ExportedVarTest   = 1
  unexportedVarTest = 2
)

var UngroupedVar = 3

var UnassignedVar int

func ExportedFunction() {}
