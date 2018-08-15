package evaluator

import (
	"github.com/stretchr/testify/assert"
	"go_interpreter/lexer"
	"go_interpreter/object"
	"go_interpreter/parser"
	"testing"
)

// Testing integer expressions e.g. "5;"
func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5*2", 10},
		{"3+2*5", 13},
		{"-4*6", -24},
		{"6/7", 0},
		{"10/5 + 2", 4},
	}

	for _, test := range tests {
		result := testEval(test.input)
		testInteger(t, result, test.expected)
	}
}

// Testing boolean expressions e.g. "true;"
func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 == 1", true},
		{"2 != 3", true},
		{"true == true", true},
		{"true != false", true},
		{"(1 < 2) == true", true},
		{"(1 > 2) == false", true},
	}

	for _, test := range tests {
		result := testEval(test.input)
		testBoolean(t, result, test.expected)
	}
}

// Testing prefix expressions e.g. "!5"
func TestEvalPrefixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, test := range tests {
		result := testEval(test.input)
		testBoolean(t, result, test.expected)
	}
}

// Testing if else expressions
func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
	}

	for _, test := range tests {
		result := testEval(test.input)

		expectedInteger, ok := test.expected.(int)
		if ok {
			testInteger(t, result, int64(expectedInteger))
		} else {
			testNull(t, result)
		}
	}
}

// Testing return statements
func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"9; return 2*5; 8;", 10},
		{"if (10 > 1) { if (10 > 1) { return 10; } return 1; }", 10},
	}

	for _, test := range tests {
		result := testEval(test.input)
		testInteger(t, result, test.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`if (10 > 1) {  if (10 > 1) { return true + false; } return 1; }`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar", "identifier not found: foobar",
		},
	}

	for _, test := range tests {
		result := testEval(test.input)

		errObj, ok := result.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", result, result)
			continue
		}

		assert.Equal(t, test.expectedMessage, errObj.Message, test.input)
	}
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5*8; a;", 40},
		{"let a = 5; let b = a; let c = a + b + 2; c;", 12},
	}

	for _, test := range tests {
		testInteger(t, testEval(test.input), test.expected)
	}
}

func TestFunctionDefinition(t *testing.T) {
	input := "fn(x) {x + 2;};"

	result := testEval(input)

	f, ok := result.(*object.Function)
	if !ok {
		t.Fatalf("Object is not function")
	}

	assert.Equal(t, 1, len(f.Parameters), "Expected number of parameters")
	assert.Equal(t, "x", f.Parameters[0].String(), "Expected parameter")
	assert.Equal(t, "(x + 2)", f.Body.String(), "Expected body")
}

func TestFunctionCall(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			"let res = fn(x) { x; }; res(5);",
			5,
		},
		{
			"let add = fn(x, y) { x + y;}; add(5+2, add(2, 3));",
			12,
		},
	}

	for _, test := range tests {
		testInteger(t, testEval(test.input), test.expected)
	}
}

func TestClosure(t *testing.T) {
	input := "let a = fn(x) { fn(y) {x+y}}; let b = a(2); b(3);"
	testInteger(t, testEval(input), 5)
}

func TestString(t *testing.T) {
	input := `"Hello world!"`
	result := testEval(input)
	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("Object isn't string")
	}

	assert.Equal(t, str.Value, "Hello world!", "Expected value of string")
}

func TestStringConcatenation(t *testing.T) {
	input := `"foo" + " " + "bar"`
	result := testEval(input)
	str, ok := result.(*object.String)
	if !ok {
		t.Fatalf("Object isn't string")
	}

	assert.Equal(t, str.Value, "foo bar", "Expected value of concatenated string")
}

// Helper method for calling eval
func testEval(input string) object.Object {
	l := lexer.BuildLexer(input)
	p := parser.BuildParser(l)

	prog := p.ParseProgram()
	env := object.BuildEnvironment()
	return Eval(prog, env)
}

// Helper method for checking integer objects
func testInteger(t *testing.T, obj object.Object, expected int64) {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Fatalf("Expected object type: Integer, actual: %T", obj)
	}

	assert.Equal(t, result.Value, expected, "Expected value")
}

// Helper method for checking boolean objects
func testBoolean(t *testing.T, obj object.Object, expected bool) {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Fatalf("Expected object type: Boolean, actual: %T", obj)
	}

	assert.Equal(t, result.Value, expected, "Expected value")
}

// Helper method for checking null objects
func testNull(t *testing.T, obj object.Object) {
	assert.Equal(t, obj, NULL, "Expected Null")
}
