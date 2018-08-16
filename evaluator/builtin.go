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
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return NewError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": &object.BuiltIn{
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments (expected = 1)")
			}

			if args[0].Type() != object.ARRAY_OBJECT {
				return NewError("argument to `first` must be array")
			}

			array := args[0].(*object.Array)
			if len(array.Elements) > 0 {
				return array.Elements[0]
			} else {
				return NULL
			}
		},
	},
	"last": &object.BuiltIn{
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments (expected = 1)")
			}

			if args[0].Type() != object.ARRAY_OBJECT {
				return NewError("argument to `first` must be array")
			}

			array := args[0].(*object.Array)
			if len(array.Elements) > 0 {
				return array.Elements[len(array.Elements)-1]
			} else {
				return NULL
			}
		},
	},
	"tail": &object.BuiltIn{
		Function: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return NewError("wrong number of arguments (expected = 1)")
			}

			if args[0].Type() != object.ARRAY_OBJECT {
				return NewError("argument to `first` must be array")
			}

			array := args[0].(*object.Array)
			length := len(array.Elements)

			if length > 0 {
				tailElements := make([]object.Object, length-1, length-1)
				copy(tailElements, array.Elements[1:length])
				return &object.Array{Elements: tailElements}
			} else {
				return NULL
			}
		},
	},
	"push": &object.BuiltIn{
		Function: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return NewError("wrong number of arguments (expected = 2)")
			}

			if args[0].Type() != object.ARRAY_OBJECT {
				return NewError("argument to `first` must be array")
			}

			array := args[0].(*object.Array)
			length := len(array.Elements)

			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, array.Elements)
			newElements[length] = args[1]
			return &object.Array{Elements: newElements}
		},
	},
}
