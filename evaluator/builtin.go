package evaluator

import "go_interpreter/object"

var builtins = map[string]*object.BuiltIn{
	"len": &object.BuiltIn{
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments (expected = 1)")
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return NewError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
}
