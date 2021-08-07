package main

import (
	_ "embed"
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

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

type MiddlewareDefinition struct{}
