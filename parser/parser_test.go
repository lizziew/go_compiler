package parser

import (
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
	} else if len(prog.Statements) != 3 {
		t.Fatalf("Expected number of statements = 3, actual = %d", len(prog.Statements))
	}

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
	} else if len(prog.Statements) != 3 {
		t.Fatalf("Expected number of statements = 3, actual = %d", len(prog.Statements))
	}

	for _, statement := range prog.Statements {
		rs, ok := statement.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("Expected type of Statement: ReturnStatement, actual: %T", statement)
		}

		assert.Equal(t, "return", rs.TokenLiteral(), "Expected ReturnStatement literal")
	}
}
