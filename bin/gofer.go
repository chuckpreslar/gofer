package main

import (
  "fmt"
  _ "github.com/chuckpreslar/gofer"

  "go/ast"
  "go/parser"
  "go/token"

  "os"

  "path/filepath"

  "strings"
)

type Import string

type Name struct {
  original string // Original name for declaration.
  unique   string // Unique name for declaration.
}

type Declaration struct {
  constant bool
  assigned bool
  value    string
  typ      string
  name     Name
}

type Function struct {
  reciever string
  block    string
  name     Name
}

// Gofer binary variables.
var (
  Debug  = true
  GoPath = os.Getenv("GOPATH") // The users GOPATH environment variable.

  TaskFiles     = make([]string, 0)      // Files within task directories.
  TaskImports   = make([]Import, 0)      // Unique imports used within the task files.
  TaskVariables = make([]Declaration, 0) // Variables from task files.
  TaskConstants = make([]Declaration, 0) // Constants from task files.
  TaskFunctions = make([]Function, 0)    // Functions from task files.

  TaskMap = make(map[string]map[string]string) // Map of original names to unqiue names.
)

func WalkGoPath() error {
  err := filepath.Walk(GoPath, func(path string, info os.FileInfo, err error) error {
    if info.IsDir() && strings.HasSuffix(path, "tasks") {
      TaskFiles = append(TaskFiles, path)
    }

    return err
  })

  return err
}

func ParseTaskFiles() error {
  for _, dir := range TaskFiles {
    fset := token.NewFileSet()
    packages, err := parser.ParseDir(fset, dir, nil, parser.AllErrors)

    if nil != err {
      return err
    }

    if err = ParsePackages(packages); nil != err {
      return err
    }
  }

  return nil
}

func ParsePackages(packages map[string]*ast.Package) error {
  for _, pkg := range packages {
    for location, file := range pkg.Files {
      if err := ParseFile(file, location); nil != err {
        return err
      }
    }
  }

  return nil
}

func ParseFile(file *ast.File, location string) error {
  for _, decl := range file.Decls {
    if err := ParseDeclaration(decl); nil != err {
      return err
    }
  }

  return nil
}

func ParseDeclaration(decl ast.Decl) (err error) {
  switch decl.(type) {
  case *ast.GenDecl:
    err = ParseGeneralDeclaration(decl.(*ast.GenDecl))
  case *ast.FuncDecl:
    err = ParseFunctionDeclaration(decl.(*ast.FuncDecl))
  default:
    panic(fmt.Sprintf("%T", decl))
  }

  return
}

func ParseGeneralDeclaration(decl *ast.GenDecl) (err error) {
  for _, spec := range decl.Specs {
    switch spec.(type) {
    case *ast.ImportSpec:
      err = ParseImport(spec.(*ast.ImportSpec))
    case *ast.ValueSpec:
      err = ParseValue(spec.(*ast.ValueSpec))
    }
  }
  return
}

func ParseFunctionDeclaration(decl *ast.FuncDecl) (err error) {
  return
}

func ParseImport(imprt *ast.ImportSpec) (err error) {
  if IsNewImport(imprt) {
    AppendImport(imprt)
  }

  return
}

func ParseValue(value *ast.ValueSpec) (err error) {
  fmt.Println(value)
  return
}

func IsNewImport(imprt *ast.ImportSpec) bool {
  for _, i := range TaskImports {
    if i == Import(imprt.Path.Value) {
      return false
    }
  }

  return true
}

func AppendImport(imprt *ast.ImportSpec) {
  if Debug {
    fmt.Printf("Appending import %s\n", imprt.Path.Value)
  }

  TaskImports = append(TaskImports, Import(imprt.Path.Value))
}

func main() {
  err := WalkGoPath()

  if nil != err {
    panic(err)
  }

  err = ParseTaskFiles()

  if nil != err {
    panic(err)
  }
}
