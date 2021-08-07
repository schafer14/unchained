package main

import (
	_ "embed"
	"go/doc"
	"html/template"
	"io"
)

//go:embed static/main.go.tmpl
var mainTemplate string

type RuntimeIR struct {
	Module   string
	Handlers []HandlerGroup
	Deps     []*doc.Func
}

func MakeMain(fnsPkg *doc.Package, deps *doc.Package, w io.Writer) {
	cfg := RuntimeIR{
		Module:   "assets",
		Handlers: makeHandlers(fnsPkg),
		Deps:     getDepTypes(deps),
	}

	tmpl, err := template.New("main.go").Parse(mainTemplate)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, cfg)
	if err != nil {
		panic(err)
	}
}
