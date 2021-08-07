package main

import (
	"go/ast"
	"go/doc"
	"regexp"
	"strings"
)

type HandlerGroup struct {

	// Name is the name of the handler function.
	Name string

	// RoutingKey defines the routing key prefix for these handlers.
	RoutingKey string

	// Middlewares is a list of middlewares that apply to this handler
	Middlewares []MiddlewareDefinition

	// Fns is a list of fns applied to that handler
	Fns []HandlerFunction
}

type HandlerFunction struct {

	// Type is the variant of handler (ie. http)
	Type string

	// Name is the calling name for this function.
	Name string

	// RoutingKey defines the routing key prefix for this handler.
	RoutingKey string

	// Method defines the method associated with the hand handler.
	Method string

	// Middlewares is a list of middlewares that apply to this handler.
	Middlewares []MiddlewareDefinition

	// Args is a list of function arguments in order.
	Args []FunctionArg

	// Returns is a list of function return arguments.
	Return []FunctionReturns
}

// FunctionArg is a function argument
type FunctionArg struct {

	// Name is the argument name.
	Name string

	// Type is the argument type.
	Type string
}

type FunctionReturns struct{}

// makeHandler creates a set of handler groups from a go package.
func makeHandlers(fnsPkg *doc.Package) []HandlerGroup {
	handlerGroups := []HandlerGroup{}

	for _, h := range fnsPkg.Types {
		if strings.HasSuffix(h.Name, "Handler") {
			handlerGroups = append(handlerGroups, makeHandler(h))
		}
	}

	return handlerGroups
}

func makeHandler(t *doc.Type) HandlerGroup {
	var hg HandlerGroup

	docs := strings.Split(t.Doc, "\n")
	hg.Name = t.Name

LOOP:
	for _, d := range docs {
		if !strings.HasPrefix(d, "@unchained") {
			continue
		}

		words := strings.Split(d, " ")
		parts := strings.Split(words[0], ":")
		if len(parts) < 1 {
			continue
		}

		switch parts[1] {
		case "routingKey":
			if len(words) < 2 {
				hg.RoutingKey = ""
				continue LOOP
			}
			hg.RoutingKey = words[1]
		}
	}

	for _, fn := range t.Methods {
		hg.Fns = append(hg.Fns, makeFn(fn))
	}

	return hg
}

func makeFn(fn *doc.Func) HandlerFunction {
	var hf HandlerFunction
	hf.Type = "http"
	hf.Name = fn.Name

	docs := strings.Split(fn.Doc, "\n")
LOOP:
	for _, d := range docs {
		if !strings.HasPrefix(d, "@unchained") {
			continue
		}

		words := strings.Split(d, " ")
		parts := strings.Split(words[0], ":")
		if len(parts) < 1 {
			continue
		}

		switch parts[1] {
		case "routingKey":
			if len(words) < 2 {
				hf.RoutingKey = ""
				continue LOOP
			}
			hf.RoutingKey = words[1]
		case "method":
			if len(words) < 2 {
				continue LOOP
			}
			hf.Method = words[1]
		}
	}

	for _, a := range fn.Decl.Type.Params.List {
		for _, name := range a.Names {
			var arg FunctionArg
			arg.Name = name.Name
			switch a.Type.(type) {
			case *ast.Ident:
				arg.Type = a.Type.(*ast.Ident).Name
			case *ast.SelectorExpr:
				arg.Type = a.Type.(*ast.SelectorExpr).Sel.Name
			}

			hf.Args = append(hf.Args, arg)
		}
	}

	return hf
}

func (hf HandlerFunction) BodyArg() string {
	urlParams := hf.UrlParams()

LOOP:
	for _, a := range hf.Args {
		if a.Type == "Context" {
			continue
		}

		for _, p := range urlParams {
			if p == a.Name {
				continue LOOP
			}
		}

		return a.Type
	}
	return ""
}

func (hf HandlerFunction) UrlParams() []string {
	var partsRegex = regexp.MustCompile(`{.+?}`)
	b := partsRegex.FindAll([]byte(hf.RoutingKey), -1)

	p := []string{}
	for _, b := range b {
		s := string(b)
		p = append(p, s[1:len(s)-1])
	}

	return p
}

func (hf HandlerFunction) CallArgs() []string {
	urlParams := hf.UrlParams()
	args := []string{}

LOOP:
	for _, a := range hf.Args {
		if a.Type == "Context" {
			args = append(args, "r.Context()")
			continue
		}

		for _, p := range urlParams {
			if p == a.Name {
				args = append(args, p)
				continue LOOP
			}
		}

		args = append(args, "body")
	}

	return args
}

func (hf HandlerFunction) CallArgsStr() string {
	return strings.Join(hf.CallArgs(), ", ")
}
