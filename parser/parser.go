package parser

import (
	"fmt"
	"go_interpreter/ast"
	"go_interpreter/lexer"
	"go_interpreter/token"
	"log"
	"os"
	"strconv"
)

var Trace *log.Logger

// Top down operator precedence parser builds AST out of tokens
type Parser struct {
	l *lexer.Lexer // corresponding lexer

	currentToken token.Token // points to current token
	nextToken    token.Token // points to next token

	errors []string // errors when parsing

	prefixMap map[token.TokenType]parsePrefix // parse prefix expressions
	infixMap  map[token.TokenType]parseInfix  // parse infix expressions
}

func BuildParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// Set currentToken and nextToken
	p.GetNextToken()
	p.GetNextToken()

	// Prefix: Map tokens --> parsing functions
	p.prefixMap = make(map[token.TokenType]parsePrefix)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefix)
	p.registerPrefix(token.MINUS, p.parsePrefix)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)

	// Infix: Map tokens --> parsing functions
	p.infixMap = make(map[token.TokenType]parseInfix)
	p.registerInfix(token.PLUS, p.parseInfix)
	p.registerInfix(token.MINUS, p.parseInfix)
	p.registerInfix(token.SLASH, p.parseInfix)
	p.registerInfix(token.ASTERISK, p.parseInfix)
	p.registerInfix(token.EQ, p.parseInfix)
	p.registerInfix(token.NOT_EQ, p.parseInfix)
	p.registerInfix(token.LT, p.parseInfix)
	p.registerInfix(token.GT, p.parseInfix)

	// Logger (for debugging): os.Stdout or ioutil.Discard
	Trace = log.New(os.Stdout, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)

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
		p.reportExpectedTokenError(t)
		return false
	}
}

// Report errors
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) reportExpectedTokenError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token: %s, actual: %s", t, p.nextToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) reportMissingPrefixFunctionError(t token.TokenType) {
	msg := fmt.Sprintf("missing prefix function for %s", t)
	p.errors = append(p.errors, msg)
}

// Parse prefix and infix expressions
type (
	parsePrefix func() ast.Expression
	parseInfix  func(ast.Expression) ast.Expression // input is "left side" of infix operator
)

func (p *Parser) registerPrefix(t token.TokenType, f parsePrefix) {
	p.prefixMap[t] = f
}

func (p *Parser) registerInfix(t token.TokenType, f parseInfix) {
	p.infixMap[t] = f
}

const (
	_           int = iota // 0
	LOWEST                 // 1
	EQUALS                 // 2: ==
	LESSGREATER            // 3: <,>
	SUM                    // 4: +
	PRODUCT                // 5: *
	PREFIX                 // 6: -foo, !foo
	CALL                   // 7: foo(bar)
)

// Maps token types --> precedences
var precedencesMap = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

func (p *Parser) getCurrentPrecedence() int {
	precedence, ok := precedencesMap[p.currentToken.Type]
	if ok {
		return precedence
	} else {
		return LOWEST
	}
}

func (p *Parser) getNextPrecedence() int {
	precedence, ok := precedencesMap[p.nextToken.Type]
	if ok {
		return precedence
	} else {
		return LOWEST
	}
}

func (p *Parser) ParseProgram() *ast.Program {
	Trace.Println("parser.ParseProgram()")

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
	Trace.Println("  parser.parseStatement()")

	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
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

	// e.g. "5"
	p.GetNextToken()
	statement.Value = p.parseExpression(LOWEST)

	// ";"
	for p.currentToken.Type != token.SEMICOLON {
		p.GetNextToken()
	}

	Trace.Println("    parser.parseLetStatement():", statement.String())
	return statement
}

// e.g. "return 5;"
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	// "return"
	statement := &ast.ReturnStatement{Token: p.currentToken}
	p.GetNextToken()

	// e.g. "5"
	statement.Value = p.parseExpression(LOWEST)

	// ";"
	for p.currentToken.Type != token.SEMICOLON {
		p.GetNextToken()
	}

	Trace.Println("    parser.parseReturnStatement():", statement.String())
	return statement
}

// Parse expression statements e.g. "5 + foo"
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// e.g. "5"
	statement := &ast.ExpressionStatement{Token: p.currentToken}

	// Pass the lowest precedence since we didn't parse anything yet
	statement.Expression = p.parseExpression(LOWEST)

	// Optional semicolon
	if p.nextToken.Type == token.SEMICOLON {
		p.GetNextToken()
	}

	Trace.Println("    parser.parseExpressionStatement():", statement.String())
	return statement
}

// Parse expressions e.g. "5 + foo"
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefixFunc := p.prefixMap[p.currentToken.Type]
	if prefixFunc == nil {
		p.reportMissingPrefixFunctionError(p.currentToken.Type)
		return nil
	}

	Trace.Printf("      parser.parseExpression(%v)\n", precedence)
	Trace.Println("        leftExpression:", p.currentToken.Literal, p.currentToken.Type)
	leftExpression := prefixFunc()

	// Tries to find infixFunc for tokens until finds token with lower precedence
	for (p.nextToken.Type != token.SEMICOLON) && precedence < p.getNextPrecedence() {
		infixFunc := p.infixMap[p.nextToken.Type]
		if infixFunc == nil {
			Trace.Println("        Was prefix function")
			return leftExpression
		}

		p.GetNextToken()

		Trace.Println("        Is infix function")
		leftExpression = infixFunc(leftExpression)
	}

	return leftExpression
}

// Parse identifier expressions e.g. "foo"
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{p.currentToken, p.currentToken.Literal}
}

// Parse integer literal expressions e.g. "5"
func (p *Parser) parseIntegerLiteral() ast.Expression {
	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("couldn't parse %q as integer", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	return &ast.IntegerLiteral{p.currentToken, value}
}

// Parse prefix expressions e.g. "-add(1, 2)"
func (p *Parser) parsePrefix() ast.Expression {
	// e.g. "-"
	expression := &ast.Prefix{Token: p.currentToken, Operator: p.currentToken.Literal}

	p.GetNextToken()

	// e.g. "add(1, 2)"
	expression.Value = p.parseExpression(PREFIX)

	Trace.Println("        p.parsePrefix():", expression.String())
	return expression
}

// Parse infix expressions e.g. "2+foo"
func (p *Parser) parseInfix(left ast.Expression) ast.Expression {
	// e.g. "2" and "+"
	expression := &ast.Infix{Token: p.currentToken, Operator: p.currentToken.Literal, Left: left}

	// e.g. "foo"
	precedence := p.getCurrentPrecedence()
	p.GetNextToken()
	expression.Right = p.parseExpression(precedence)

	Trace.Println("        p.parseInfix():", expression.String())
	return expression
}

// Parse boolean expressions e.g. "true"
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currentToken, Value: p.currentToken.Type == token.TRUE}
}
