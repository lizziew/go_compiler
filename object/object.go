package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJECT = "INTEGER"
	BOOLEAN_OBJECT = "BOOLEAN"
	NULL_OBJECT    = "NULL"
	RETURN_OBJECT  = "RETURN"
	ERROR_OBJECT   = "ERROR"
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
