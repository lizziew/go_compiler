package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go_interpreter/ast"
	"go_interpreter/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foo = bar;", "foo", "bar"},
	}

	for _, test := range tests {
		l := lexer.BuildLexer(test.input)
		p := BuildParser(l)
		prog := p.ParseProgram()

		checkParserErrors(t, p)

		assert.Equal(t, 1, len(prog.Statements), "Expected number of statements")

		statement := prog.Statements[0]
		testLetStatement(t, statement, test.expectedIdentifier)

		value := statement.(*ast.LetStatement).Value
		testLiteral(t, value, test.expectedValue)
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foo;", "foo"},
	}

	for _, test := range tests {
		l := lexer.BuildLexer(test.input)
		p := BuildParser(l)
		prog := p.ParseProgram()

		checkParserErrors(t, p)

		assert.Equal(t, 1, len(prog.Statements), "Expected number of statements")

		statement := prog.Statements[0]
		rs, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("Expected type of Statement: ReturnStatement, actual: %T", statement)
		}

		assert.Equal(t, "return", rs.TokenLiteral(), "Expected ReturnStatement literal")
		testLiteral(t, rs.Value, test.expectedValue)
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foo;"

	l := lexer.BuildLexer(input)
	p := BuildParser(l)
	prog := p.ParseProgram()

	checkParserErrors(t, p)

	assert.Equal(t, 1, len(prog.Statements), "Expected number of statements")

	statement, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected Statement type: ExpressionStatement, actual: %T", prog.Statements[0])
	}

	ident, ok := statement.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected expression type: Identifier, actual: %T", statement.Expression)
	}
	assert.Equal(t, ident.Value, "foo", "Expected value of ident")
	assert.Equal(t, ident.TokenLiteral(), "foo", "Expected TokenLiteral() of foo")
}

func TestIntegerValueExpression(t *testing.T) {
	input := "5;"

	l := lexer.BuildLexer(input)
	p := BuildParser(l)
	prog := p.ParseProgram()

	checkParserErrors(t, p)

	assert.Equal(t, 1, len(prog.Statements), "Expected number of statements")

	statement, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected Statement type: ExpressionStatement, actual: %T", prog.Statements[0])
	}

	literal, ok := statement.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expected expression type: IntegerLiteral, actual: %T", statement.Expression)
	}

	assert.Equal(t, literal.Value, int64(5), "Expected value of literal")
	assert.Equal(t, literal.TokenLiteral(), "5", "Expected TokenLiteral() of literal")
}

func TestPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!foo", "!", "foo"},
		{"-foo", "-", "foo"},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, test := range prefixTests {
		l := lexer.BuildLexer(test.input)
		p := BuildParser(l)
		prog := p.ParseProgram()

		checkParserErrors(t, p)

		assert.Equal(t, 1, len(prog.Statements), "Expected number of statements")

		statement, ok := prog.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expected Statement type: ExpressionStatement, actual: %T", prog.Statements[0])
		}

		prefix, ok := statement.Expression.(*ast.Prefix)
		if !ok {
			t.Fatalf("Expected expression type: Prefix, actual: %T", statement.Expression)
		}

		assert.Equal(t, test.operator, prefix.Operator, "Expected prefix operator")

		testLiteral(t, prefix.Value, test.value)
	}
}

func TestInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input    string
		left     interface{}
		operator string
		right    interface{}
	}{
		{"1 + 2;", 1, "+", 2},
		{"1 - 2;", 1, "-", 2},
		{"1 * 2;", 1, "*", 2},
		{"1 / 2;", 1, "/", 2},
		{"1 > 2;", 1, ">", 2},
		{"1 < 2;", 1, "<", 2},
		{"1 == 2;", 1, "==", 2},
		{"1 != 2;", 1, "!=", 2},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, test := range infixTests {
		l := lexer.BuildLexer(test.input)
		p := BuildParser(l)
		prog := p.ParseProgram()

		checkParserErrors(t, p)

		assert.Equal(t, 1, len(prog.Statements), "Expected number of statements")

		statement, ok := prog.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expected Statement type: ExpressionStatement, actual:%T", prog.Statements[0])
		}

		infix, ok := statement.Expression.(*ast.Infix)
		if !ok {
			t.Fatalf("Expected expression type: Infix, actual:%T", statement.Expression)
		}

		testLiteral(t, infix.Left, test.left)
		assert.Equal(t, test.operator, infix.Operator, "Expected infix operator")
		testLiteral(t, infix.Right, test.right)
	}
}

func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
	}

	for _, test := range tests {
		l := lexer.BuildLexer(test.input)
		p := BuildParser(l)
		prog := p.ParseProgram()

		checkParserErrors(t, p)

		assert.Equal(t, test.expected, prog.String(), "Expected precedence")
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.BuildLexer(input)
	p := BuildParser(l)
	prog := p.ParseProgram()

	checkParserErrors(t, p)

	assert.Equal(t, 1, len(prog.Statements), "Expected number of statements")

	statement, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected Statement type: ExpressionStatement, actual: %T", prog.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.If)
	if !ok {
		t.Fatalf("Expected Expression type: If, actual: %T", statement.Expression)
	}

	// Test condition
	testInfix(t, expression.Condition, "x", "<", "y")

	// Test consequence
	assert.Equal(t, len(expression.Consequence.Statements), 1,
		"Expected number of Consequence statements")

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected Statement Type: ExpressionStatement, actual: %T",
			expression.Consequence.Statements[0])
	}

	testIdentifier(t, consequence.Expression, "x")

	// Test alternative
	if expression.Alternative != nil {
		t.Fatalf("Alternative should be empty")
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.BuildLexer(input)
	p := BuildParser(l)
	prog := p.ParseProgram()

	checkParserErrors(t, p)

	// Test condition
	assert.Equal(t, 1, len(prog.Statements), "Expected number of statements")
	statement, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			prog.Statements[0])
	}
	expression, ok := statement.Expression.(*ast.If)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", statement.Expression)
	}
	testInfix(t, expression.Condition, "x", "<", "y")

	// Test consequence
	assert.Equal(t, 1, len(expression.Consequence.Statements), "Expected number of statements")
	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			expression.Consequence.Statements[0])
	}
	testIdentifier(t, consequence.Expression, "x")

	// Test alternative
	assert.Equal(t, 1, len(expression.Alternative.Statements), "Expected number of statements")
	alternative, ok := expression.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			expression.Alternative.Statements[0])
	}
	testIdentifier(t, alternative.Expression, "y")
}

func TestFunctionExpression(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	l := lexer.BuildLexer(input)
	p := BuildParser(l)
	prog := p.ParseProgram()

	checkParserErrors(t, p)

	assert.Equal(t, 1, len(prog.Statements), "Expected number of statements")
	statement, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected Statement type: ExpressionStatement, actual: %T", prog.Statements[0])
	}
	function, ok := statement.Expression.(*ast.Function)
	if !ok {
		t.Fatalf("Expected Expression type: Function, actual: %T", statement.Expression)
	}
	assert.Equal(t, 2, len(function.Parameters), "Expected number of parameters")
	testLiteral(t, function.Parameters[0], "x")
	testLiteral(t, function.Parameters[1], "y")

	assert.Equal(t, 1, len(function.Body.Statements), "Expected number of body statements")
	bodyStatement, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Expected function body statement type: ExpressionStatement")
	}

	testInfix(t, bodyStatement.Expression, "x", "+", "y")
}

// Helper method for checking parser errors
func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	for _, msg := range errors {
		t.Errorf("Parser error: %q", msg)
	}

	t.Fatalf("Parser had a total of %d errors", len(errors))
}

// Helper method for let statements
func testLetStatement(t *testing.T, statement ast.Statement, name string) {
	assert.Equal(t, statement.TokenLiteral(), "let", "LetStatement")

	expectedLetStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Fatalf("expected type of Statement: LetStatement, actual: %T", statement)
	}

	assert.Equal(t, expectedLetStatement.Name.Value, name, "Name of identifier")
	assert.Equal(t, expectedLetStatement.Name.TokenLiteral(), name, "Name of identifier")
}

// Helper method for infix expression
func testInfix(t *testing.T, expression ast.Expression, left interface{},
	operator string, right interface{}) {
	infixExpression, ok := expression.(*ast.Infix)
	if !ok {
		t.Fatalf("Expected Expression type: ast.Infix, actual: %T", expression)
	}

	testLiteral(t, infixExpression.Left, left)
	assert.Equal(t, infixExpression.Operator, operator)
	testLiteral(t, infixExpression.Right, right)
}

// Helper method for identifier and integer literal expressions
func testLiteral(t *testing.T, expression ast.Expression, value interface{}) {
	switch typedValue := value.(type) {
	case int:
		testIntegerLiteral(t, expression, int64(typedValue))
		return
	case int64:
		testIntegerLiteral(t, expression, typedValue)
		return
	case string:
		testIdentifier(t, expression, typedValue)
		return
	case bool:
		testBooleanLiteral(t, expression, typedValue)
		return
	}

	t.Fatalf("Type of literal expression is not identifier or integer literal: %T", value)
}

// Helper method for boolean expression
func testBooleanLiteral(t *testing.T, expression ast.Expression, value bool) {
	b, ok := expression.(*ast.Boolean)
	if !ok {
		t.Fatalf("Expected Expression type: ast.Boolean, actual: %T", expression)
	}

	assert.Equal(t, b.Value, value, "Expected value of bool")
	assert.Equal(t, b.TokenLiteral(), fmt.Sprintf("%t", value), "Expected token literal of bool")
}

// Helper method for identifier expression
func testIdentifier(t *testing.T, expression ast.Expression, value string) {
	i, ok := expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Expected Expression type: ast.Identifier, actual: %T", expression)
	}

	assert.Equal(t, value, i.Value, "Expected value of identifier")
	assert.Equal(t, value, i.TokenLiteral(), "Expected token literal of identifier")
}

// Helper method for integer literal expression
func testIntegerLiteral(t *testing.T, expression ast.Expression, value int64) {
	il, ok := expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expected Expression type: IntegerLiteral, actual: %T", expression)
	}

	assert.Equal(t, value, il.Value, "Expected value")
	assert.Equal(t, fmt.Sprintf("%d", value), il.TokenLiteral(), "Expected token literal")
}
