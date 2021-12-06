package main

import (
	"errors"
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

type memberInitialization Field

func (m memberInitialization) Read([]byte) (int, error) {
	return 0, errors.New("not implemented")
}

func (m memberInitialization) WriteTo(w io.Writer) (int64, error) {
	typ := camel(arrayName(m.Type))
	name := camel(m.Name)
	mw := mustWriter{w: w}
	mw.Fprintf(`
		%s: new%s(assertion),`, name, typ)
	return mw.total, nil
}

type initialField Field

func (initialField) Read([]byte) (int, error) {
	return 0, errors.New("not implemented")
}

func (f initialField) WriteTo(w io.Writer) (int64, error) {
	mw := mustWriter{w: w}
	mw.Fprintf(`

func(parent *%s) %s(checkers ...%s) *Assertion {
	parent.%sCheckers = checkers
	return parent.assertion
}`, f.Parent, camel(f.Name), interfaceTypeName(f.Type), privateCamel(f.Name))
	return mw.total, nil
}

type performCheck Field

func (m performCheck) Read([]byte) (int, error) {
	return 0, errors.New("not implemented")
}

func (field performCheck) WriteTo(w io.Writer) (int64, error) {
	mw := mustWriter{w: w}
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
	return mw.total, nil
}

type AssertionField Field

func (a AssertionField) WriteTo(w io.Writer) (int64, error) {
	mw := mustWriter{w: w}
	mw.Fprintf(`
	%s %s`, camel(a.Parent), a.Parent)
	return mw.total, nil
}

func (a AssertionField) Read([]byte) (int, error) {
	return 0, errors.New("not implemented")
}

type AssertionInitialization Field

func (a AssertionInitialization) WriteTo(w io.Writer) (int64, error) {
	mw := mustWriter{w: w}
	mw.Fprintf(`
	assertion.%s = new%s(&assertion)`, a.Parent, a.Parent)
	return mw.total, nil
}

func (a AssertionInitialization) Read([]byte) (int, error) {
	return 0, errors.New("not implemented")
}

type MyType Field

func (a MyType) Read([]byte) (int, error) {
	return 0, errors.New("not implemented")
}

func (a MyType) WriteTo(w io.Writer) (int64, error) {
	mw := mustWriter{w: w}
	mw.Fprintf(`
	return &assertion
}`)
	return mw.total, nil
}

type AssertionImplementation Field

func (m AssertionImplementation) Read([]byte) (int, error) {
	return 0, errors.New("not implemented")
}

func (m AssertionImplementation) WriteTo(w io.Writer) (int64, error) {
	mw := mustWriter{w: w}
	mw.Fprintf(`

func (a *Assertion) %s(t *testing.T, desc string, val %s){
	a.%s.%s(t, desc+" %s", val)
}`, interfaceMethodName(m.Parent), maybeQualify(m.Parent), m.Parent, interfaceMethodName(m.Parent), m.Parent)
	return mw.total, nil
}

type forwardCheck Field

func (m forwardCheck) Read([]byte) (int, error) {
	return 0, errors.New("not implemented")
}

func (field forwardCheck) WriteTo(w io.Writer) (int64, error) {
	mw := mustWriter{w: w}
	mw.Fprintf(`
	parent.%s.%s(t, desc + ".%s", val.%s)`,
		camel(field.Name),
		interfaceMethodName(field.Type),
		field.Name,
		field.Name,
	)
	return mw.total, nil
}

type arrayTypeDef struct {
	typ  string
	next string
}

func (a arrayTypeDef) Read([]byte) (int, error) {
	return 0, errors.New("not implemented")
}

func (a arrayTypeDef) WriteTo(w io.Writer) (int64, error) {
	mw := mustWriter{w: w}
	mw.Fprintf(`

type %s struct {
	assertion      *Assertion
	lengthCheckers []IntChecker
	nth            map[int]%s

	ForEach %s
}`, camel(arrayName(a.typ)), camel(arrayName(a.next)), camel(arrayName(a.next)))
	mw.Fprintf(`

func new%s(assertion *Assertion) %s {
	return %s {
		assertion: assertion,
		nth:       make(map[int]%s),
		ForEach:   new%s(assertion),
	}
}`, camel(arrayName(a.typ)), camel(arrayName(a.typ)), camel(arrayName(a.typ)), camel(arrayName(a.next)), camel(arrayName(a.next)))
	mw.Fprintf(`

func (a *%s) Nth(i int) %s {
	prev, ok := a.nth[i]
	if ok {
		return prev
	}
	result := new%s(a.assertion)
	a.nth[i] = result
	return result
}`, camel(arrayName(a.typ)), camel(arrayName(a.next)), camel(arrayName(a.next)))
	mw.Fprintf(`

func (a *%s) Length(checkers ...IntChecker) *Assertion {
	a.lengthCheckers = checkers
	return a.assertion
}`, camel(arrayName(a.typ)))
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
		camel(arrayName(a.typ)),
		interfaceMethodName(a.typ),
		maybeQualify(a.typ),
		interfaceMethodName(a.next),
		interfaceMethodName(a.next),
	)
	return mw.total, nil
}

// Implementations writes assertion structs.
type Implementations []Field

func (imp Implementations) Read([]byte) (int, error) {
	return 0, errors.New("not implemented")
}

// WriteTo the given writer.
func (imp Implementations) WriteTo(w io.Writer) (int64, error) {
	mw := mustWriter{w: w}

	seen := make(map[string]bool)

	mw.Fprintf(`

type Assertion struct {`)
	var lastParent string

	for _, field := range imp {
		if field.Parent == lastParent {
			continue
		}
		lastParent = field.Parent
		AssertionField(field).WriteTo(w)
	}

	mw.Fprintf(`
}

func New() *Assertion {
	var assertion Assertion`)
	lastParent = ""

	for _, field := range imp {
		if field.Parent == lastParent {
			continue
		}
		lastParent = field.Parent
		AssertionInitialization(field).WriteTo(w)
	}
	mw.Fprintf(`
	return &assertion
}`)

	lastParent = ""

	for _, field := range imp {
		if field.Parent == lastParent {
			continue
		}
		lastParent = field.Parent
		AssertionImplementation(field).WriteTo(w)
	}
	lastParent = ""

	for _, f := range imp {
		if f.Parent == lastParent {
			continue
		}
		lastParent = f.Parent
		mw.Fprintf(`

type %s struct {
	assertion *Assertion
`, f.Parent)
		for _, field := range imp {
			if field.Parent != f.Parent {
				continue
			}
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
		assertion: assertion,`, f.Parent, f.Parent, f.Parent)
		for _, field := range imp {
			if field.Parent != f.Parent {
				continue
			}
			if isBasicType(field.Type) {
				continue
			}
			memberInitialization(field).WriteTo(w)
		}
		mw.Fprintf(`
	}
}`)
		for _, field := range imp {
			if field.Parent != f.Parent {
				continue
			}
			if !isBasicType(field.Type) {
				continue
			}
			initialField(field).WriteTo(w)
		}

		mw.Fprintf(`

func(parent *%s) %s(t *testing.T, desc string, val %s) {`, f.Parent, interfaceMethodName(f.Parent), maybeQualify(f.Parent))
		for _, field := range imp {
			if field.Parent != f.Parent {
				continue
			}
			if isBasicType(field.Type) {
				performCheck(field).WriteTo(w)
				continue
			}
			forwardCheck(field).WriteTo(w)
		}
		mw.Fprintf(`
}`)
	}
	seen = make(map[string]bool)

	for _, field := range imp {
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
			arrayTypeDef{typ: typ, next: next}.WriteTo(w)
		}
	}
	return mw.total, nil
}

// Field is a field in a struct.
//
// Name uses selector syntax for nested structs:
//
//   A {
//   	B struct {
//   		C int // Field { Parent: "A", Name: "A.B.C", Type: "int" }
//   	}
//   }
//
// Type uses array/slice declaration syntax:
//
// A {
//  B []int // Field { Parent: "A", Name: "B", Type: "[]int" }
//  C [2][]Person // Field { Parent: "A", Name: "C", Type: "[2][]Person"
// }
//
type Field struct {
	Parent string
	Name   string
	Type   string
}

func fieldFromIdent(parent string, name string, e *ast.Ident) Field {
	curr := Field{Parent: parent, Name: name, Type: e.Name}
	if name == "" {
		curr.Name = curr.Type
	}
	return curr
}

func fieldFromSelectorExpr(parent string, name string, e *ast.SelectorExpr) Field {
	return Field{Parent: parent, Name: name, Type: e.Sel.Name}
}

type arrayType struct {
	Type   string
	Length int
}

func (a arrayType) String() string {
	if a.Length == 0 {
		return "[]" + a.Type
	}
	return fmt.Sprintf("[%d]", a.Length) + a.Type
}

func mustAtoi(val string) int {
	result, err := strconv.Atoi(val)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func fieldsFromArray(parent string, name string, e *ast.ArrayType) []Field {
	result := fieldsFromExpr(parent, name, e.Elt)
	for i := range result {
		a := arrayType{Type: result[i].Type}
		if e.Len != nil {
			a.Length = mustAtoi(e.Len.(*ast.BasicLit).Value)
		}
		result[i].Type = a.String()
	}
	return result
}

func fieldsFromStruct(parent string, name string, e *ast.StructType) []Field {
	var result []Field
	for _, field := range e.Fields.List {
		currName := name
		for _, n := range field.Names {
			if currName != "" {
				currName += "."
			}
			currName += n.Name
		}
		result = append(result, fieldsFromExpr(parent, currName, field.Type)...)
	}
	return result
}

func fieldsFromExpr(parent string, name string, e ast.Expr) []Field {
	switch e.(type) {
	case *ast.Ident:
		return []Field{fieldFromIdent(parent, name, e.(*ast.Ident))}
	case *ast.SelectorExpr:
		return []Field{fieldFromSelectorExpr(parent, name, e.(*ast.SelectorExpr))}
	case *ast.ParenExpr:
		return fieldsFromExpr(parent, name, e.(*ast.ParenExpr).X)
	case *ast.StarExpr:
		return fieldsFromExpr(parent, name, e.(*ast.StarExpr).X)
	case *ast.ArrayType:
		return fieldsFromArray(parent, name, e.(*ast.ArrayType))
	case *ast.StructType:
		return fieldsFromStruct(parent, name, e.(*ast.StructType))
	default:
		log.Fatalf("%T is not supported", e)
	}
	return nil
}

func fieldsFromFile(parent string, decls []ast.Decl) []Field {
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
			result = append(result, fieldsFromExpr(ts.Name.Name, "", ts.Type)...)
		}
	}
	return result
}

var header = `// Code generated by cmd/assert; DO NOT EDIT.

package assert

import (
	"testing"
	"fmt"

	"github.com/slcjordan/poc"
)`

func layout(fields []Field) []io.Reader {
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

	var result []io.Reader
	result = append(result, strings.NewReader(header))
	result = append(result, Implementations(fields))
	return result
}

func parse(filename string) *ast.File {
	result, err := parser.ParseFile(token.NewFileSet(), filename, nil, 0)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func main() {
	var filename string
	flag.StringVar(&filename, "filename", "", "")
	flag.Parse()
	objs := fieldsFromFile("", parse(filename).Decls)
	readers := layout(objs)

	base := filepath.Base(filename)
	file, err := os.Create(filepath.Join("test/assert", base))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for _, r := range readers {
		_, err := io.Copy(file, r)
		if err != nil {
			log.Fatal(err)
		}
	}
}