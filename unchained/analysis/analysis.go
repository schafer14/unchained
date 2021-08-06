package main

import (
	_ "embed"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

//go:embed static/main.go.tmpl
var mainTemplate string

func main() {

	if len(os.Args) < 2 {
		fmt.Println("provide a path")
		os.Exit(1)
	}

	root := os.Args[1]
	fnsDir := path.Join(root, "/fns")
	depsDir := path.Join(root, "/dependencies")

	fnsPkg := mustMakePackage(fnsDir)
	depsPkg := mustMakePackage(depsDir)

	f, err := os.Create(path.Join(root, "main.go"))
	if err != nil {
		panic(err)
	}

	MakeMain(fnsPkg, depsPkg, f)

	f.Close()
}

type MainConfig struct {
	Module   string
	Handlers []*doc.Type
	Deps     []*doc.Func
}

func MakeMain(fnsPkg *doc.Package, deps *doc.Package, w io.Writer) {
	funcMap := template.FuncMap{
		"MkHandler": func(s string) string {
			return strings.ReplaceAll(strings.ToLower(s), "handler", "")
		},
		"MkPath": func(s string) string {
			return strings.ReplaceAll(strings.ToLower(s), "handle", "")
		},
		"FirstArg": func(x *doc.Func) string {
			return x.Decl.Type.Params.List[1].Type.(*ast.Ident).Name
		},
	}

	cfg := MainConfig{
		Module:   "assets",
		Handlers: getHandlerTypes(fnsPkg),
		Deps:     getDepTypes(deps),
	}

	tmpl, err := template.New("main.go").Funcs(funcMap).Parse(mainTemplate)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, cfg)
	if err != nil {
		panic(err)
	}
}

func getHandlerTypes(fnsPkg *doc.Package) []*doc.Type {
	types := []*doc.Type{}

	for _, t := range fnsPkg.Types {
		if strings.HasSuffix(t.Name, "Handler") {
			types = append(types, t)
		}
	}

	return types
}

func getDepTypes(depPkg *doc.Package) []*doc.Func {
	funcs := []*doc.Func{}

	for _, f := range depPkg.Funcs {
		if strings.HasPrefix(f.Name, "Provide") {
			funcs = append(funcs, f)
		}
	}

	return funcs
}

func MakePackage(dir string) (*doc.Package, error) {
	fset := token.NewFileSet()
	var astFiles []*ast.File
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		p := path.Join(dir, f.Name())
		if f.IsDir() {
			continue
		}

		if !strings.HasSuffix(p, ".go") {
			continue
		}

		bytes, err := ioutil.ReadFile(p)
		if err != nil {
			panic(err)
		}

		astFiles = append(astFiles, mustParse(fset, p, string(bytes)))

	}

	// Compute package documentation with examples.
	p, err := doc.NewFromFiles(fset, astFiles, "package")
	if err != nil {
		panic(err)
	}

	return p, nil
}

func mustMakePackage(dir string) *doc.Package {
	p, err := MakePackage(dir)
	if err != nil {
		panic(err)
	}

	return p
}

func mustParse(fset *token.FileSet, filename, src string) *ast.File {
	f, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	return f
}
