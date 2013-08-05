package main

import (
  "fmt"
  "go/ast"
  "go/parser"
  "go/token"
  "io"
  "os"
  "os/exec"
  "path"
  "path/filepath"
  "strings"
  "text/template"
  "time"
)

type Import struct {
  Path string
}

type TemplateData struct {
  Imports []Import
}

// Gofer binary constants.
const (
  SOURCE_PREFIX        = "/src/"
  PACKAGE_NAME         = "tasks"
  EXPECTED_IMPORT      = "gofer"
  TEMPLATE_DESTINATION = "gofer_task_definitions_%v.go"
)

// Gofer binary variables.
var (
  GoPath          = os.Getenv("GOPATH") // The users GOPATH environment variable.
  TaskDirectories = make([]string, 0)   // Task directories.
  Template        = TemplateData{}
)

func WalkGoPath() error {
  err := filepath.Walk(GoPath, func(path string, info os.FileInfo, err error) error {
    if info.IsDir() && strings.HasSuffix(path, PACKAGE_NAME) {
      TaskDirectories = append(TaskDirectories, path)
    }

    return err
  })

  return err
}

func ParseTaskDirectories() error {
  for _, dir := range TaskDirectories {
    fset := token.NewFileSet()
    packages, err := parser.ParseDir(fset, dir, nil, parser.AllErrors)

    if nil != err {
      return err
    }

    if err = ParsePackages(packages, dir); nil != err {
      return err
    }
  }

  return nil
}

func ParsePackages(packages map[string]*ast.Package, dir string) (err error) {
  for _, pkg := range packages {
    file := ast.MergePackageFiles(pkg, ast.FilterImportDuplicates)

    if IsGoferTaskFile(file) {
      AddImport(dir)
    }
  }
  return
}

func IsGoferTaskFile(file *ast.File) bool {
  for _, imprt := range file.Imports {
    if PACKAGE_NAME == file.Name.String() && strings.ContainsAny(imprt.Path.Value, EXPECTED_IMPORT) {
      return true
    }
  }

  return false
}

func AddImport(dir string) {
  imprt := strings.TrimLeft(strings.Replace(dir, GoPath, "", 1), SOURCE_PREFIX)
  Template.Imports = append(Template.Imports, Import{imprt})
}

func CompileTemplate() (tmpl *template.Template) {
  tmpl = template.Must(template.New("gofer").Parse(`
    package main

    import (
      "os"
      "github.com/chuckpreslar/gofer"

      // Imported task files.
    {{range .Imports}}
      _ "{{.Path}}"
    {{end}}
    )

    func main() {
      gofer.Preform(os.Args[1:]...)
    }
  `))

  return
}

func WriteTemplate(destination string, tmpl *template.Template) (err error) {
  file, err := os.Create(destination)

  if nil != err {
    return
  }

  defer file.Close()

  err = tmpl.Execute(file, Template)

  return
}

func RemoveTemplate(location string) {
  os.Remove(location)
}

func main() {
  err := WalkGoPath()

  if nil != err {
    panic(err)
  }

  err = ParseTaskDirectories()

  if nil != err {
    panic(err)
  }

  tmpl := CompileTemplate()
  dir := path.Join(os.TempDir(), fmt.Sprintf(TEMPLATE_DESTINATION, time.Now().Unix()))

  err = WriteTemplate(dir, tmpl)

  if nil != err {
    panic(err)
  }

  defer RemoveTemplate(dir)

  arguments := append([]string{"run", dir}, os.Args[1:]...)

  command := exec.Command("go", arguments...)

  stdout, err := command.StdoutPipe()

  if nil != err {
    panic(err)
  }

  stderr, err := command.StderrPipe()

  if nil != err {
    panic(err)
  }

  err = command.Start()

  if nil != err {
    panic(err)
  }

  go io.Copy(os.Stdout, stdout)
  go io.Copy(os.Stderr, stderr)

  err = command.Wait()

  if nil != err {
    panic(err)
  }
}
