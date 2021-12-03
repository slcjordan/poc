package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// InterfaceDeclarations writes assertion interface declarations.
type InterfaceDeclarations []Field

func camel(val string) string {
	return strings.ToUpper(val[:1]) + val[1:]
}

func interfaceTypeName(typ string) string {
	return fmt.Sprintf("%sChecker", camel(typ))
}

func interfaceMethodName(typ string) string {
	return fmt.Sprintf("Check%s", camel(typ))
}

func maybeQualify(val string) string {
	switch val {
	case "int", "array", "string":
		return val
	default:
		return "poc." + val
	}
}

// WriteTo the given writer.
func (d InterfaceDeclarations) WriteTo(w io.Writer) (int64, error) {
	var types []string
	seen := make(map[string]bool)
	for _, field := range d {
		if seen[field.Type] {
			continue
		}
		seen[field.Type] = true
		types = append(types, field.Type)
	}
	sort.Strings(types)
	var total int64

	for _, typ := range types {
		n, err := fmt.Fprintf(w, `

type %s interface {
	%s(*testing.T, %s)
}`, interfaceTypeName(typ), interfaceMethodName(typ), maybeQualify(typ))
		total += int64(n)
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

// Field is a field in an object.
type Field struct {
	Name   string
	Type   string
	Parent string
}

func expr(parent string, name string, e ast.Expr) []Field {
	curr := Field{Parent: parent, Name: name}
	switch e.(type) {
	case *ast.Ident:
		curr.Type = e.(*ast.Ident).Name
	case *ast.SelectorExpr:
		curr.Type = e.(*ast.SelectorExpr).Sel.Name
	case *ast.ParenExpr:
		return expr(parent, name, e.(*ast.ParenExpr).X)
	case *ast.StarExpr:
		return expr(parent, name, e.(*ast.StarExpr).X)
	case *ast.ArrayType:
		curr.Type = "array"
		return append(expr("array", name, e.(*ast.ArrayType).Elt), curr)
	case *ast.StructType:
		var result []Field
		for _, field := range e.(*ast.StructType).Fields.List {
			currName := name
			for _, fn := range field.Names {
				if currName != "" {
					currName += "."
				}
				currName += fn.Name
			}
			result = append(result, expr(parent, currName, field.Type)...)
		}
		return result
	case *ast.ChanType:
		curr.Type = "chan"
		return expr("chan", name, e.(*ast.ChanType).Value)
	default:
		log.Fatalf("%T is not supported", e)
	}
	return []Field{curr}
}

func fields(parent string, decls []ast.Decl) []Field {
	var result []Field
	for _, val := range decls {
		decl, ok := val.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range decl.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			result = append(result, expr(ts.Name.Name, "", ts.Type)...)
		}
	}
	return result
}

func parse(filename string) *ast.File {
	result, err := parser.ParseFile(token.NewFileSet(), filename, nil, 0)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

var header = []byte(`package assert

import "testing"`)

func main() {
	var filename string
	flag.StringVar(&filename, "filename", "", "")
	flag.Parse()
	nodes := fields("", parse(filename).Decls)

	ifaceFilename := strings.Replace(filepath.Base(filename), ".go", "_interfaces.go", 1)
	file, err := os.Create(filepath.Join("../../test/assert", ifaceFilename))
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.Write(header)
	if err != nil {
		log.Fatal(err)
	}
	_, err = InterfaceDeclarations(nodes).WriteTo(file)
	if err != nil {
		log.Fatal(err)
	}
}
