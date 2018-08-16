package object

import (
	"bytes"
	"fmt"
	"go_interpreter/ast"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJECT  = "INTEGER"
	BOOLEAN_OBJECT  = "BOOLEAN"
	NULL_OBJECT     = "NULL"
	RETURN_OBJECT   = "RETURN"
	ERROR_OBJECT    = "ERROR"
	FUNCTION_OBJECT = "FUNCTION"
	STRING_OBJECT   = "STRING"
	BUILTIN_OBJECT  = "BUILTIN"
	ARRAY_OBJECT    = "ARRAY"
)

// Generic object
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer type
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJECT
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// Boolean type
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJECT
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

// Null type
type Null struct{}

func (n *Null) Type() ObjectType {
	return NULL_OBJECT
}

func (n *Null) Inspect() string {
	return "null"
}

// Return type
type Return struct {
	Value Object
}

func (r *Return) Type() ObjectType {
	return RETURN_OBJECT
}

func (r *Return) Inspect() string {
	return r.Value.Inspect()
}

// Error type
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJECT
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

// Function type
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJECT
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// String type
type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return STRING_OBJECT
}

func (s *String) Inspect() string {
	return s.Value
}

// Built in function type
type BuiltInFunction func(args ...Object) Object

type BuiltIn struct {
	Function BuiltInFunction
}

func (b *BuiltIn) Type() ObjectType {
	return BUILTIN_OBJECT
}

func (b *BuiltIn) Inspect() string {
	return "built in function"
}

// Array type
type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType {
	return ARRAY_OBJECT
}

func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
