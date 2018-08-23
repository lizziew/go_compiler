package compiler

import (
	"github.com/stretchr/testify/assert"
	"go_interpreter/ast"
	"go_interpreter/bytecode"
	"go_interpreter/lexer"
	"go_interpreter/object"
	"go_interpreter/parser"
	"strconv"
	"testing"
)

type testCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []bytecode.Instructions
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []testCase{
		{
			"1 + 2",
			[]interface{}{1, 2},
			[]bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			"3 * 2",
			[]interface{}{3, 2},
			[]bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpMul),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			"3 / 2",
			[]interface{}{3, 2},
			[]bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpDiv),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			"3 - 2",
			[]interface{}{3, 2},
			[]bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpSub),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			"-5",
			[]interface{}{5},
			[]bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpMinus),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	testCompiler(t, tests)
}

func TestBoolean(t *testing.T) {
	tests := []testCase{
		{
			"true",
			[]interface{}{},
			[]bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			"false",
			[]interface{}{},
			[]bytecode.Instructions{
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			"1 > 2",
			[]interface{}{1, 2},
			[]bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpGreater),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			"1 < 2",
			[]interface{}{2, 1},
			[]bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpGreater),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			"1 == 2",
			[]interface{}{1, 2},
			[]bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpEqual),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			"1 != 2",
			[]interface{}{1, 2},
			[]bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpNotEqual),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			"!true",
			[]interface{}{},
			[]bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpBang),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	testCompiler(t, tests)
}

func TestConditional(t *testing.T) {
	tests := []testCase{
		{
			"if (true) { 10 }; 3333;",
			[]interface{}{10, 3333},
			[]bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),              // 0000
				bytecode.Make(bytecode.OpJumpNotTruthy, 10), // 0001
				bytecode.Make(bytecode.OpConstant, 0),       // 0004
				bytecode.Make(bytecode.OpJump, 11),          // 0007
				bytecode.Make(bytecode.OpNull),              // 0010
				bytecode.Make(bytecode.OpPop),               // 0011
				bytecode.Make(bytecode.OpConstant, 1),       // 0012
				bytecode.Make(bytecode.OpPop),               // 0015
			},
		},
		{
			"if (true) { 10 } else { 20 }; 3333; ",
			[]interface{}{10, 20, 3333},
			[]bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),              // 0000
				bytecode.Make(bytecode.OpJumpNotTruthy, 10), // 0001
				bytecode.Make(bytecode.OpConstant, 0),       // 0004
				bytecode.Make(bytecode.OpJump, 13),          // 0007 (Skip executing the alternative)
				bytecode.Make(bytecode.OpConstant, 1),       // 0010
				bytecode.Make(bytecode.OpPop),               // 0013
				bytecode.Make(bytecode.OpConstant, 2),       // 0014
				bytecode.Make(bytecode.OpPop),               // 0017
			},
		},
	}

	testCompiler(t, tests)
}

func TestGlobalLet(t *testing.T) {
	tests := []testCase{
		{
			"let one = 1; let two = 2;",
			[]interface{}{1, 2},
			[]bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpSetGlobal, 1),
			},
		},
		{
			"let one = 1; one;",
			[]interface{}{1},
			[]bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			"let one = 1; let two = one; two;",
			[]interface{}{1},
			[]bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpSetGlobal, 1),
				bytecode.Make(bytecode.OpGetGlobal, 1),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	testCompiler(t, tests)
}

func TestString(t *testing.T) {
	tests := []testCase{
		{
			`"foo"`,
			[]interface{}{"foo"},
			[]bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			`"foo" + "bar"`,
			[]interface{}{"foo", "bar"},
			[]bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	testCompiler(t, tests)
}

func TestArray(t *testing.T) {
	tests := []testCase{
		{
			"[]",
			[]interface{}{},
			[]bytecode.Instructions{
				bytecode.Make(bytecode.OpArray, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			"[1,2]",
			[]interface{}{1, 2},
			[]bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpArray, 2),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	testCompiler(t, tests)
}

// Helper method to parse input string
func parse(input string) *ast.Program {
	l := lexer.BuildLexer(input)
	p := parser.BuildParser(l)
	return p.ParseProgram()
}

// Helper method to test compiler
func testCompiler(t *testing.T, tests []testCase) {
	for _, test := range tests {
		program := parse(test.input)

		compiler := BuildCompiler()
		err := compiler.Compile(program)
		if err != nil {
			t.Fatalf("Compiler error: %s", err)
		}

		bytecode := compiler.Bytecode()
		testInstructions(t, test.expectedInstructions, bytecode.Instructions)
		testConstants(t, test.expectedConstants, bytecode.Constants)
	}
}

// Helper method to test instructions
func testInstructions(t *testing.T, expectedList []bytecode.Instructions, actual bytecode.Instructions) {
	expected := joinInstructions(expectedList)

	assert.Equal(t, len(actual), len(expected))

	for i, b := range expected {
		assert.Equal(t, actual[i], b, strconv.Itoa(i))
	}
}

// Helper method to join instructions (needed because input is a slice of slice of bytes)
func joinInstructions(input []bytecode.Instructions) bytecode.Instructions {
	result := bytecode.Instructions{}

	for _, b := range input {
		result = append(result, b...)
	}

	return result
}

// Helper method to test constants
func testConstants(t *testing.T, expected []interface{}, actual []object.Object) {
	assert.Equal(t, len(expected), len(actual))

	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			testIntegerObject(t, int64(constant), actual[i])
		case string:
			testStringObject(t, constant, actual[i])
		}
	}
}

// Helper method to test integer objects
func testIntegerObject(t *testing.T, expected int64, actual object.Object) {
	result, ok := actual.(*object.Integer)
	if !ok {
		t.Fatalf("Object is not integer")
	}

	assert.Equal(t, result.Value, expected)
}

// Helper method to test string objects
func testStringObject(t *testing.T, expected string, actual object.Object) {
	result, ok := actual.(*object.String)
	if !ok {
		t.Fatalf("Object is not string")
	}

	assert.Equal(t, result.Value, expected)
}
