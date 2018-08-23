package vm

import (
	"github.com/stretchr/testify/assert"
	"go_interpreter/ast"
	"go_interpreter/compiler"
	"go_interpreter/lexer"
	"go_interpreter/object"
	"go_interpreter/parser"
	"testing"
)

type testCase struct {
	input    string
	expected interface{}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []testCase{
		{"1", 1},
		{"2", 2},
		{"1+2", 3},
		{"3-5", -2},
		{"8*9", 72},
		{"4/3", 1},
		{"(3 + 9)*2", 24},
		{"2 * (3 + 9)", 24},
		{"3 + 9 * 2", 21},
		{"-5", -5},
		{"-3 + 9", 6},
		{"(15/-3) + 7", 2},
	}

	testVM(t, tests)
}

func TestBoolean(t *testing.T) {
	tests := []testCase{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 == 1", true},
		{"2 != 3", true},
		{"true == true", true},
		{"true != false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!(if (false) { 5; })", true},
	}

	testVM(t, tests)
}

func TestConditional(t *testing.T) {
	tests := []testCase{
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 > 2) { 10 }", Null},
		{"if ((if (false) { 10 })) { 10 } else { 20 }", 20},
	}

	testVM(t, tests)
}

func TestGlobalLet(t *testing.T) {
	tests := []testCase{
		{"let one = 1; one", 1},
		{"let one = 1; let two = 2; one + two", 3},
		{"let one = 1; let two = one + one; one + two", 3},
	}

	testVM(t, tests)
}

func TestString(t *testing.T) {
	tests := []testCase{
		{`"foo"`, "foo"},
		{`"foo" + "bar"`, "foobar"},
	}

	testVM(t, tests)
}

func testVM(t *testing.T, tests []testCase) {
	for _, test := range tests {
		prog := parse(test.input)

		c := compiler.BuildCompiler()
		err := c.Compile(prog)
		if err != nil {
			t.Fatalf("Compiler error: %s", err)
		}

		vm := BuildVM(c.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("VM error: %s", err)
		}

		lastPopped := vm.LastPopped()
		testExpectedObject(t, test.expected, lastPopped)
	}
}

func parse(input string) *ast.Program {
	l := lexer.BuildLexer(input)
	p := parser.BuildParser(l)
	return p.ParseProgram()
}

func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	switch expected := expected.(type) {
	case int:
		testIntegerObject(t, int64(expected), actual)
	case bool:
		testBooleanObject(t, bool(expected), actual)
	case string:
		testStringObject(t, string(expected), actual)
	case *object.Null:
		if actual != Null {
			t.Fatalf("Expected null, but actual is not")
		}
	}
}

func testIntegerObject(t *testing.T, expected int64, actual object.Object) {
	result, ok := actual.(*object.Integer)
	if !ok {
		t.Fatalf("Object is not an integer")
	}

	assert.Equal(t, result.Value, expected)
}

func testBooleanObject(t *testing.T, expected bool, actual object.Object) {
	result, ok := actual.(*object.Boolean)
	if !ok {
		t.Fatalf("Object is not a boolean")
	}

	assert.Equal(t, result.Value, expected)
}

func testStringObject(t *testing.T, expected string, actual object.Object) {
	result, ok := actual.(*object.String)
	if !ok {
		t.Fatalf("Object is not a string")
	}

	assert.Equal(t, result.Value, expected)
}
