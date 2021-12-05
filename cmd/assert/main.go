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
	"strconv"
	"strings"
)

func isBasicType(val string) bool {
	return map[string]bool{
		"bool":       true,
		"byte":       true,
		"complex64":  true,
		"complex128": true,
		"float32":    true,
		"float64":    true,
		"int":        true,
		"int16":      true,
		"int32":      true,
		"int64":      true,
		"int8":       true,
		"rune":       true,
		"string":     true,
		"uint":       true,
		"uint16":     true,
		"uint32":     true,
		"uint64":     true,
		"uint8":      true,
		"uintptr":    true,
	}[val]
}

type mustWriter struct {
	total int64
	w     io.Writer
}

func (m *mustWriter) Fprintf(format string, a ...interface{}) {
	n, err := fmt.Fprintf(m.w, format, a...)
	m.total += int64(n)
	if err != nil {
		log.Fatal(err)
	}
}

func maybeQualify(val string) string {
	if isBasicType(val) {
		return val
	}
	if strings.HasPrefix(val, "[") {
		return val[:strings.LastIndex(val, "]")+1] + maybeQualify(val[strings.LastIndex(val, "]")+1:])
	}
	return "poc." + val
}

func camel(val string) string {
	parts := strings.SplitN(val, ".", 2)
	if len(parts) > 1 {
		return camel(parts[0]) + camel(parts[1])
	}
	return strings.ToUpper(val[:1]) + val[1:]
}

func privateCamel(val string) string {
	result := camel(val)
	return strings.ToLower(result[:1]) + result[1:]
}

func dimensions(val string) int {
	return len(strings.Split(val, "[")) - 1
}

func arrayName(val string) string {
	dims := dimensions(val)
	if dims == 0 {
		return val
	}
	return fmt.Sprintf("%sArray%dD", val[strings.LastIndex(val, "]")+1:], dims)
}

func interfaceTypeName(typ string) string {
	return fmt.Sprintf("%sChecker", camel(arrayName(typ)))
}

func interfaceMethodName(typ string) string {
	return fmt.Sprintf("Check%s", camel(arrayName(typ)))
}

// Implementations writes assertion structs.
type Implementations []Field

// WriteTo the given writer.
func (imp Implementations) WriteTo(w io.Writer) (int64, error) {
	fields := make([]Field, len(imp))
	copy(fields, imp)
	byName := func(i int, j int) bool { return fields[i].Name < fields[j].Name }
	byIsBasicType := func(i int, j int) bool {
		if isBasicType(fields[i].Type) == isBasicType(fields[j].Type) {
			return false
		}
		return isBasicType(fields[i].Type)
	}
	byParent := func(i int, j int) bool { return fields[i].Parent < fields[j].Parent }
	sort.SliceStable(fields, byName)
	sort.SliceStable(fields, byIsBasicType)
	sort.SliceStable(fields, byParent)
	mw := mustWriter{w: w}

	declarations := make(map[string][]Field)
	seen := make(map[string]bool)
	var parents []string

	for _, field := range fields {
		if !seen[field.Parent] {
			parents = append(parents, field.Parent)
			seen[field.Parent] = true
		}
		declarations[field.Parent] = append(declarations[field.Parent], field)
	}
	mw.Fprintf(`

type Assertion struct {`)

	for _, parent := range parents {
		mw.Fprintf(`
	%s %s`, camel(parent), parent)
	}

	mw.Fprintf(`
}

func New() *Assertion {
	var assertion Assertion`)
	for _, parent := range parents {
		mw.Fprintf(`
	assertion.%s = new%s(&assertion)`, parent, parent)
	}
	mw.Fprintf(`
	return &assertion
}`)

	for _, parent := range parents {
		mw.Fprintf(`

func (a *Assertion) %s(t *testing.T, desc string, val %s){
	a.%s.%s(t, desc+" %s", val)
}`, interfaceMethodName(parent), maybeQualify(parent), parent, interfaceMethodName(parent), parent)
	}

	for _, parent := range parents {
		mw.Fprintf(`

type %s struct {
	assertion *Assertion
`, parent)
		for _, field := range declarations[parent] {
			if isBasicType(field.Type) {
				mw.Fprintf(`	%sCheckers []%s
`, privateCamel(field.Name), interfaceTypeName(field.Type))
			} else {
				mw.Fprintf(`
	%s %s`, camel(field.Name), camel(arrayName(field.Type)))
			}
		}
		mw.Fprintf(`
}`)
		mw.Fprintf(`

func new%s(assertion *Assertion) %s {
	return %s {
		assertion: assertion,`, parent, parent, parent)
		for _, field := range declarations[parent] {
			typ := camel(arrayName(field.Type))
			name := camel(field.Name)
			if isBasicType(field.Type) {
				continue
			}
			mw.Fprintf(`
		%s: new%s(assertion),`, name, typ)
		}
		mw.Fprintf(`
	}
}`)
		for _, field := range declarations[parent] {
			if !isBasicType(field.Type) {
				continue
			}
			mw.Fprintf(`

func(parent *%s) %s(checkers ...%s) *Assertion {
	parent.%sCheckers = checkers
	return parent.assertion
}`, parent, camel(field.Name), interfaceTypeName(field.Type), privateCamel(field.Name))
		}

		mw.Fprintf(`

func(parent *%s) %s(t *testing.T, desc string, val %s) {`, parent, interfaceMethodName(parent), maybeQualify(parent))
		for _, field := range declarations[parent] {
			if isBasicType(field.Type) {
				currName := "val." + field.Name
				if isBasicType(field.Name) {
					currName = field.Name + "(val)"
				}
				mw.Fprintf(`
	for _, checker := range  parent.%sCheckers {
		checker.%s(t, desc + ".%s", %s)
	}`,
					privateCamel(field.Name),
					interfaceMethodName(field.Type),
					field.Name,
					currName,
				)
				continue
			}
			mw.Fprintf(`
	parent.%s.%s(t, desc + ".%s", val.%s)`,
				camel(field.Name),
				interfaceMethodName(field.Type),
				field.Name,
				field.Name,
			)
		}
		mw.Fprintf(`
}`)
	}
	seen = make(map[string]bool)

	for _, field := range fields {
		if !strings.HasPrefix(field.Type, "[") {
			continue
		}
		for dims := dimensions(field.Type); dims > 0; dims-- {
			typ := strings.Join(strings.Split(field.Type, "]")[dims-1:], "]")
			if seen[typ] {
				continue
			}
			seen[typ] = true
			next := strings.Join(strings.Split(field.Type, "]")[dims:], "]")
			mw.Fprintf(`

type %s struct {
	assertion      *Assertion
	lengthCheckers []IntChecker
	nth            map[int]%s

	ForEach %s
}`, camel(arrayName(typ)), camel(arrayName(next)), camel(arrayName(next)))
			mw.Fprintf(`

func new%s(assertion *Assertion) %s {
	return %s {
		assertion: assertion,
		nth:       make(map[int]%s),
		ForEach:   new%s(assertion),
	}
}`, camel(arrayName(typ)), camel(arrayName(typ)), camel(arrayName(typ)), camel(arrayName(next)), camel(arrayName(next)))
			mw.Fprintf(`

func (a *%s) Nth(i int) %s {
	prev, ok := a.nth[i]
	if ok {
		return prev
	}
	result := new%s(a.assertion)
	a.nth[i] = result
	return result
}`, camel(arrayName(typ)), camel(arrayName(next)), camel(arrayName(next)))
			mw.Fprintf(`

func (a *%s) Length(checkers ...IntChecker) *Assertion {
	a.lengthCheckers = checkers
	return a.assertion
}`, camel(arrayName(typ)))
			mw.Fprintf(`

func (a *%s) %s(t *testing.T, desc string, val %s) {
	for _, checker := range a.lengthCheckers {
		checker.CheckInt(t, desc+".length", len(val))
	}
	for i, checker := range a.nth {
		checker.%s(t, desc+fmt.Sprintf("[%%d]", i), val[i])
	}
	for _, curr := range val {
		a.ForEach.%s(t, desc+".ForEach", curr)
	}
}`,
				camel(arrayName(typ)),
				interfaceMethodName(typ),
				maybeQualify(typ),
				interfaceMethodName(next),
				interfaceMethodName(next),
			)
		}
	}
	return mw.total, nil
}

// Field is a field in an object.
type Field struct {
	Name   string
	Type   string
	Parent string
	Len    int
}

func expr(parent string, name string, e ast.Expr) []Field {
	curr := Field{Parent: parent, Name: name}
	switch e.(type) {
	case *ast.Ident:
		curr.Type = e.(*ast.Ident).Name
		if name == "" {
			curr.Name = curr.Type
		}
	case *ast.SelectorExpr:
		curr.Type = e.(*ast.SelectorExpr).Sel.Name
	case *ast.ParenExpr:
		return expr(parent, name, e.(*ast.ParenExpr).X)
	case *ast.StarExpr:
		return expr(parent, name, e.(*ast.StarExpr).X)
	case *ast.ArrayType:
		result := expr(parent, name, e.(*ast.ArrayType).Elt)
		for i := range result {
			if e.(*ast.ArrayType).Len != nil {
				length, err := strconv.Atoi(e.(*ast.ArrayType).Len.(*ast.BasicLit).Value)
				if err != nil {
					log.Fatal(err)
				}
				result[i].Type = fmt.Sprintf("[%d]", length) + result[i].Type
			} else {
				result[i].Type = "[]" + result[i].Type
			}
		}
		return result
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

var header = []byte(`// Code generated by cmd/assert; DO NOT EDIT.

package assert

import (
	"testing"
	"fmt"

	"github.com/slcjordan/poc"
)`)

func main() {
	var filename string
	flag.StringVar(&filename, "filename", "", "")
	flag.Parse()
	nodes := fields("", parse(filename).Decls)

	base := filepath.Base(filename)
	file, err := os.Create(filepath.Join("test/assert", base))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	_, err = file.Write(header)
	if err != nil {
		log.Fatal(err)
	}
	_, err = Implementations(nodes).WriteTo(file)
	if err != nil {
		log.Fatal(err)
	}
}
