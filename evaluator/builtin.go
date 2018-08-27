package evaluator

import "go_interpreter/object"

var builtins = map[string]*object.BuiltIn{
	"len":   object.GetBuiltin("len"),
	"first": object.GetBuiltin("first"),
	"last":  object.GetBuiltin("last"),
	"tail":  object.GetBuiltin("tail"),
	"push":  object.GetBuiltin("push"),
	"print": object.GetBuiltin("print"),
}
