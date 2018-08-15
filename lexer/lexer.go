package lexer

import (
	"go_interpreter/token"
)

// Converts source code to tokens
type Lexer struct {
	input           string
	currentPosition int  // position that lexer points to in input
	nextPosition    int  // next position after current position
	currentChar     byte // character at current position
}

func BuildLexer(input string) *Lexer {
	lexer := &Lexer{input: input}

	// Initialize currentPosition, nextPosition, currentChar
	lexer.advanceCharacter()

	return lexer
}

// Read next character and advance lexer
func (l *Lexer) advanceCharacter() {
	if l.nextPosition >= len(l.input) {
		l.currentChar = 0 // ASCII code for null character
	} else {
		l.currentChar = l.input[l.nextPosition]
	}

	l.currentPosition = l.nextPosition
	l.nextPosition += 1
}

// Read next token and advance lexer
func (l *Lexer) advanceToken(constraint func(byte) bool) string {
	startPosition := l.currentPosition

	for constraint(l.currentChar) {
		l.advanceCharacter()
	}

	return l.input[startPosition:l.currentPosition]
}

// Read next character, but without advancing lexer
func (l *Lexer) peekCharacter() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.nextPosition]
	}
}

// Skip whitespace in between tokens
func (l *Lexer) skipWhitespace() {
	for l.currentChar == ' ' || l.currentChar == '\t' ||
		l.currentChar == '\n' || l.currentChar == '\r' {
		l.advanceCharacter()
	}
}

// Get next token
func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	var t token.Token

	switch l.currentChar {
	case '=':
		if l.peekCharacter() == '=' {
			l.advanceCharacter()
			t = token.Token{token.EQ, string("=" + string(l.currentChar))}
		} else {
			t = token.Token{token.ASSIGN, string(l.currentChar)}
		}
	case '!':
		if l.peekCharacter() == '=' {
			l.advanceCharacter()
			t = token.Token{token.NOT_EQ, string("!" + string(l.currentChar))}
		} else {
			t = token.Token{token.BANG, string(l.currentChar)}
		}
	case ';':
		t = token.Token{token.SEMICOLON, string(l.currentChar)}
	case '(':
		t = token.Token{token.LPAREN, string(l.currentChar)}
	case ')':
		t = token.Token{token.RPAREN, string(l.currentChar)}
	case ',':
		t = token.Token{token.COMMA, string(l.currentChar)}
	case '+':
		t = token.Token{token.PLUS, string(l.currentChar)}
	case '{':
		t = token.Token{token.LBRACE, string(l.currentChar)}
	case '}':
		t = token.Token{token.RBRACE, string(l.currentChar)}
	case '-':
		t = token.Token{token.MINUS, string(l.currentChar)}
	case '/':
		t = token.Token{token.SLASH, string(l.currentChar)}
	case '*':
		t = token.Token{token.ASTERISK, string(l.currentChar)}
	case '<':
		t = token.Token{token.LT, string(l.currentChar)}
	case '>':
		t = token.Token{token.GT, string(l.currentChar)}
	case '"':
		t = token.Token{token.STRING, l.readString()}
	case 0:
		t = token.Token{token.EOF, ""}
	default:
		if isLetter(l.currentChar) {
			t.Literal = l.advanceToken(isLetter)
			t.Type = token.GetIdentifier(t.Literal)
			return t
		} else if isDigit(l.currentChar) {
			t.Literal = l.advanceToken(isDigit)
			t.Type = token.INT
			return t
		} else {
			t = token.Token{token.ILLEGAL, string(l.currentChar)}
		}
	}

	l.advanceCharacter()
	return t
}

// Helper function
func (l *Lexer) readString() string {
	startPosition := l.currentPosition + 1

	for {
		l.advanceCharacter()
		if l.currentChar == '"' || l.currentChar == 0 {
			break
		}
	}

	return l.input[startPosition:l.currentPosition]
}

// Helper function
func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

// Helper function
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
