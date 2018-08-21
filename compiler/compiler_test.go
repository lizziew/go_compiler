package compiler

import (
	"github.com/stretchr/testify/assert"
	"go_interpreter/ast"
	"go_interpreter/bytecode"
	"go_interpreter/lexer"
	"go_interpreter/object"
	"go_interpreter/parser"
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
		assert.Equal(t, actual[i], b)
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
