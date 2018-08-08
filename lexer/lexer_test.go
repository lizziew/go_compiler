package lexer

import (
	"go_interpreter/token"
	"testing"
)

func TestDelimitersOperators(t *testing.T) {
	input := "=+(){},;"

	expectedTokens := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	testLexer(t, input, expectedTokens)
}

func TestAddFunction(t *testing.T) {
	input := `let five = 5;
					 let ten = 10;

					 let add = fn(x, y) {
						x+y;
					 }

					 let result = add(five, ten);`

	expectedTokens := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	testLexer(t, input, expectedTokens)
}

func TestSingleCharacterTokens(t *testing.T) {
	input := `!-/*5;
						5 < 10 > 5`

	expectedTokens := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.EOF, ""},
	}

	testLexer(t, input, expectedTokens)
}

func TestMultiCharacterTokens(t *testing.T) {
	input := `if (5 < 10) {
							return true; 
						} else {
							return false;
						}`

	expectedTokens := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	testLexer(t, input, expectedTokens)
}

func TestDoubleCharacterTokens(t *testing.T) {
	input := `10 == 10;
						10 != 9;`

	expectedTokens := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	testLexer(t, input, expectedTokens)
}

func testLexer(t *testing.T, input string, expectedTokens []struct {
	expectedType    token.TokenType
	expectedLiteral string
}) {
	l := BuildLexer(input)

	for i, expectedToken := range expectedTokens {
		actualToken := l.NextToken()

		if actualToken.Type != expectedToken.expectedType {
			t.Fatalf("Wrong TokenType at %d: expected=%q, actual=%q",
				i, expectedToken.expectedType, actualToken.Type)
		}

		if actualToken.Literal != expectedToken.expectedLiteral {
			t.Fatalf("Wrong Literal at %d: expected=%q, actual=%q",
				i, expectedToken.expectedLiteral, actualToken.Literal)
		}
	}
}
