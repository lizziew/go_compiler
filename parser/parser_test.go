package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go_interpreter/ast"
	"go_interpreter/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `let x = 5;
						let y = 10;
						let foobar = 838383;`

	l := lexer.BuildLexer(input)
	p := BuildParser(l)
	prog := p.ParseProgram()

	checkParserErrors(t, p)

	if prog == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	assert.Equal(t, 3, len(prog.Statements), "Expected number of statements")

	expected := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, test := range expected {
		statement := prog.Statements[i]
		if !testLetStatement(t, statement, test.expectedIdentifier) {
			return
		}
	}
}

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

func testLetStatement(t *testing.T, statement ast.Statement, name string) bool {
	assert.Equal(t, statement.TokenLiteral(), "let", "LetStatement")

	expectedLetStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Fatalf("expected type of Statement: LetStatement, actual: %T", statement)
		return false
	}

	assert.Equal(t, expectedLetStatement.Name.Value, name, "Name of identifier")

	assert.Equal(t, expectedLetStatement.Name.TokenLiteral(), name, "Name of identifier")

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `return 5;
						return 10;
						return 993322;`

	l := lexer.BuildLexer(input)
	p := BuildParser(l)
	prog := p.ParseProgram()

	checkParserErrors(t, p)

	if prog == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	assert.Equal(t, 3, len(prog.Statements), "Expected number of statements")

	for _, statement := range prog.Statements {
		rs, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("Expected type of Statement: ReturnStatement, actual: %T", statement)
		}

		assert.Equal(t, "return", rs.TokenLiteral(), "Expected ReturnStatement literal")
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
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
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

		testIntegerLiteral(t, prefix.Value, test.integerValue)
	}
}

func testIntegerLiteral(t *testing.T, expression ast.Expression, value int64) {
	il, ok := expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Expected Expression type: IntegerLiteral, actual: %T", il)
	}

	assert.Equal(t, value, il.Value, "Expected value")
	assert.Equal(t, fmt.Sprintf("%d", value), il.TokenLiteral(), "Expected token literal")
}

func TestInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input    string
		left     int64
		operator string
		right    int64
	}{
		{"1 + 2;", 1, "+", 2},
		{"1 - 2;", 1, "-", 2},
		{"1 * 2;", 1, "*", 2},
		{"1 / 2;", 1, "/", 2},
		{"1 > 2;", 1, ">", 2},
		{"1 < 2;", 1, "<", 2},
		{"1 == 2;", 1, "==", 2},
		{"1 != 2;", 1, "!=", 2},
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

		testIntegerLiteral(t, infix.Left, test.left)
		assert.Equal(t, test.operator, infix.Operator, "Expected infix operator")
		testIntegerLiteral(t, infix.Right, test.right)
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
	}

	for _, test := range tests {
		l := lexer.BuildLexer(test.input)
		p := BuildParser(l)
		prog := p.ParseProgram()

		checkParserErrors(t, p)

		assert.Equal(t, test.expected, prog.String(), "Expected precedence")
	}
}
