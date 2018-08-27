package object

import "fmt"

var Builtins = []struct {
	Name    string
	Builtin *BuiltIn
}{
	{
		"len",
		&BuiltIn{
			Function: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments (expected = 1)")
				}

				switch arg := args[0].(type) {
				case *String:
					return &Integer{Value: int64(len(arg.Value))}
				case *Array:
					return &Integer{Value: int64(len(arg.Elements))}
				default:
					return newError("argument to `len` not supported, got %s", args[0].Type())
				}
			},
		},
	},
	{
		"first",
		&BuiltIn{
			Function: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments (expected = 1)")
				}

				if args[0].Type() != ARRAY_OBJECT {
					return newError("argument to `first` must be array")
				}

				array := args[0].(*Array)
				if len(array.Elements) > 0 {
					return array.Elements[0]
				} else {
					return nil
				}
			},
		},
	},
	{
		"last",
		&BuiltIn{
			Function: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments (expected = 1)")
				}

				if args[0].Type() != ARRAY_OBJECT {
					return newError("argument to `first` must be array")
				}

				array := args[0].(*Array)
				if len(array.Elements) > 0 {
					return array.Elements[len(array.Elements)-1]
				} else {
					return nil
				}
			},
		},
	},
	{
		"tail",
		&BuiltIn{
			Function: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments (expected = 1)")
				}

				if args[0].Type() != ARRAY_OBJECT {
					return newError("argument to `first` must be array")
				}

				array := args[0].(*Array)
				length := len(array.Elements)

				if length > 0 {
					tailElements := make([]Object, length-1, length-1)
					copy(tailElements, array.Elements[1:length])
					return &Array{Elements: tailElements}
				} else {
					return nil
				}
			},
		},
	},
	{
		"push",
		&BuiltIn{
			Function: func(args ...Object) Object {
				if len(args) != 2 {
					return newError("wrong number of arguments (expected = 2)")
				}

				if args[0].Type() != ARRAY_OBJECT {
					return newError("argument to `first` must be array")
				}

				array := args[0].(*Array)
				length := len(array.Elements)

				newElements := make([]Object, length+1, length+1)
				copy(newElements, array.Elements)
				newElements[length] = args[1]
				return &Array{Elements: newElements}
			},
		},
	},
	{
		"print",
		&BuiltIn{
			Function: func(args ...Object) Object {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}
				return nil
			},
		},
	},
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func GetBuiltin(name string) *BuiltIn {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin
		}
	}
	return nil
}
