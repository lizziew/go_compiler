package parser

import (
	"fmt"
	"go_interpreter/ast"
	"go_interpreter/lexer"
	"go_interpreter/token"
)

// Top down operator precedence parser builds AST out of tokens
type Parser struct {
	l *lexer.Lexer // corresponding lexer

	currentToken token.Token // points to current token
	nextToken    token.Token // points to next token

	errors []string // errors when parsing
}

func BuildParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// Set currentToken and nextToken
	p.GetNextToken()
	p.GetNextToken()

	return p
}

func (p *Parser) GetNextToken() {
	p.currentToken = p.nextToken
	p.nextToken = p.l.NextToken()
}

func (p *Parser) GetExpectNextToken(t token.TokenType) bool {
	if p.nextToken.Type == t {
		p.GetNextToken()
		return true
	} else {
		p.reportError(t)
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) reportError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token: %s, actual: %s", t, p.nextToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) ParseProgram() *ast.Program {
	// Construct root Node of AST
	prog := &ast.Program{}
	prog.Statements = []ast.Statement{}

	// Iterates over tokens in input until EOF
	for p.currentToken.Type != token.EOF {
		statement := p.parseStatement()

		if statement != nil {
			prog.Statements = append(prog.Statements, statement)
		}

		p.GetNextToken()
	}

	return prog
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

// e.g. "let x = 5;"
func (p *Parser) parseLetStatement() *ast.LetStatement {
	// "let"
	statement := &ast.LetStatement{Token: p.currentToken}

	// e.g. "x"
	if !p.GetExpectNextToken(token.IDENT) {
		return nil
	}
	statement.Name = &ast.Identifier{p.currentToken, p.currentToken.Literal}

	// "="
	if !p.GetExpectNextToken(token.ASSIGN) {
		return nil
	}

	// TODO: skip expression for now, e.g. "5"
	for p.currentToken.Type != token.SEMICOLON {
		p.GetNextToken()
	}

	return statement
}

// e.g. "return 5;"
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	// "return"
	statement := &ast.ReturnStatement{Token: p.currentToken}

	// TODO: skip expression for now, e.g. "5"
	p.GetNextToken()
	for p.currentToken.Type != token.SEMICOLON {
		p.GetNextToken()
	}

	return statement
}
