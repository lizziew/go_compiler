package parser

import (
	"fmt"
	"github.com/fatih/color"
	"go_interpreter/ast"
	"go_interpreter/lexer"
	"go_interpreter/token"
	"strconv"
)

var PRINT_PARSE = false

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
	p.registerPrefix(token.LPAREN, p.parseGrouped)
	p.registerPrefix(token.IF, p.parseIf)
	p.registerPrefix(token.FUNCTION, p.parseFunction)
	p.registerPrefix(token.STRING, p.parseString)
	p.registerPrefix(token.LSQUARE, p.parseArray)
	p.registerPrefix(token.LBRACE, p.parseHash)

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
	p.registerInfix(token.LPAREN, p.parseCall)
	p.registerInfix(token.LSQUARE, p.parseIndex)

	return p
}

func (p *Parser) GetNextToken() {
	p.currentToken = p.nextToken
	p.nextToken = p.l.NextToken()

	if PRINT_PARSE {
		color.Red("Current token: %s", p.currentToken)
	}
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
	INDEX                  // 8: array[index]
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
	token.LPAREN:   CALL,
	token.LSQUARE:  INDEX,
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
	if PRINT_PARSE {
		color.Cyan("CALL parser.ParseProgram()")
	}

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
	if PRINT_PARSE {
		color.Cyan("  CALL parser.parseStatement()")
	}

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
	if PRINT_PARSE {
		color.Cyan("    CALL parser.parseLetStatement()")
	}
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
	if p.nextToken.Type == token.SEMICOLON {
		p.GetNextToken()
	}

	if PRINT_PARSE {
		color.Blue("    RET parser.parseLetStatement():%s", statement.String())
	}
	return statement
}

// e.g. "return 5;"
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	if PRINT_PARSE {
		color.Cyan("    CALL parser.parseReturnStatement()")
	}
	// "return"
	statement := &ast.ReturnStatement{Token: p.currentToken}
	p.GetNextToken()

	// e.g. "5"
	statement.Value = p.parseExpression(LOWEST)

	// ";"
	if p.nextToken.Type == token.SEMICOLON {
		p.GetNextToken()
	}

	if PRINT_PARSE {
		color.Blue("    RET parser.parseReturnStatement():%s", statement.String())
	}
	return statement
}

// Parse expression statements e.g. "5 + foo"
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	if PRINT_PARSE {
		color.Cyan("    CALL parser.parseExpressionStatement()")
	}
	// e.g. "5"
	statement := &ast.ExpressionStatement{Token: p.currentToken}

	// Pass the lowest precedence since we didn't parse anything yet
	statement.Expression = p.parseExpression(LOWEST)

	// Optional semicolon
	if p.nextToken.Type == token.SEMICOLON {
		p.GetNextToken()
	}

	if PRINT_PARSE {
		color.Blue("    RET parser.parseExpressionStatement():%s", statement.String())
	}
	return statement
}

// Parse block statement
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	if PRINT_PARSE {
		color.Cyan("      CALL parseBlockStatement()")
	}
	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = []ast.Statement{}

	p.GetNextToken()
	for p.currentToken.Type != token.RBRACE && p.currentToken.Type != token.EOF {
		statement := p.parseStatement()
		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}

		p.GetNextToken()
	}

	if PRINT_PARSE {
		color.Blue("      RET parser.parseBlockStatement(): %s", block.String())
	}
	return block
}

// Parse expressions e.g. "5 + foo"
func (p *Parser) parseExpression(precedence int) ast.Expression {
	if PRINT_PARSE {
		color.Cyan("      CALL parser.parseExpression(%v)\n", precedence)
	}
	prefixFunc := p.prefixMap[p.currentToken.Type]
	if prefixFunc == nil {
		p.reportMissingPrefixFunctionError(p.currentToken.Type)
		return nil
	}

	if PRINT_PARSE {
		color.Yellow("      EXEC leftExpression: %s %s", p.currentToken.Literal, p.currentToken.Type)
	}
	leftExpression := prefixFunc()

	// Tries to find infixFunc for tokens until finds token with lower precedence
	for (p.nextToken.Type != token.SEMICOLON) && precedence < p.getNextPrecedence() {
		infixFunc := p.infixMap[p.nextToken.Type]
		if infixFunc == nil {
			if PRINT_PARSE {
				color.Blue("      RET parser.parseExpression(): %s", leftExpression.String())
			}
			return leftExpression
		}

		p.GetNextToken()

		if PRINT_PARSE {
			color.Yellow("      EXEC is infix function")
		}
		leftExpression = infixFunc(leftExpression)
	}

	if leftExpression != nil {
		if PRINT_PARSE {
			color.Blue("      RET parser.parseExpression(): %s", leftExpression.String())
		}
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
	if PRINT_PARSE {
		color.Cyan("      CALL p.parsePrefix()")
	}
	// e.g. "-"
	expression := &ast.Prefix{Token: p.currentToken, Operator: p.currentToken.Literal}

	p.GetNextToken()

	// e.g. "add(1, 2)"
	expression.Value = p.parseExpression(PREFIX)

	if PRINT_PARSE {
		color.Blue("      RET p.parsePrefix(): %s", expression.String())
	}
	return expression
}

// Parse infix expressions e.g. "2+foo"
func (p *Parser) parseInfix(left ast.Expression) ast.Expression {
	if PRINT_PARSE {
		color.Cyan("      CALL p.parseInfix()")
	}
	// e.g. "2" and "+"
	expression := &ast.Infix{Token: p.currentToken, Operator: p.currentToken.Literal, Left: left}

	// e.g. "foo"
	precedence := p.getCurrentPrecedence()
	p.GetNextToken()
	expression.Right = p.parseExpression(precedence)

	if expression != nil {
		if PRINT_PARSE {
			color.Blue("      RET p.parseInfix(): %s", expression.String())
		}
	}
	return expression
}

// Parse boolean expressions e.g. "true"
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currentToken, Value: p.currentToken.Type == token.TRUE}
}

// Parse grouped expressions e.g. "(5+5)*2"
func (p *Parser) parseGrouped() ast.Expression {
	if PRINT_PARSE {
		color.Cyan("      CALL p.parseGrouped()")
	}

	// "("
	p.GetNextToken()

	expression := p.parseExpression(LOWEST)

	// ")"
	if !p.GetExpectNextToken(token.RPAREN) {
		return nil
	} else {
		if PRINT_PARSE {
			color.Blue("      RET p.parseGrouped():", expression.String())
		}
		return expression
	}
}

// Parse if expressions e.g. "if (4 < 5) { x } else { y }"
func (p *Parser) parseIf() ast.Expression {
	if PRINT_PARSE {
		color.Cyan("      CALL p.parseIf()")
	}
	// "if"
	expression := &ast.If{Token: p.currentToken}

	// "("
	if !p.GetExpectNextToken(token.LPAREN) {
		return nil
	}

	// e.g. "4 < 5"
	p.GetNextToken()
	expression.Condition = p.parseExpression(LOWEST)

	// ")"
	if !p.GetExpectNextToken(token.RPAREN) {
		return nil
	}

	// "{"
	if !p.GetExpectNextToken(token.LBRACE) {
		return nil
	}

	// e.g. "x"
	expression.Consequence = p.parseBlockStatement()

	// else
	if p.nextToken.Type == token.ELSE {
		p.GetNextToken()

		// {
		if !p.GetExpectNextToken(token.LBRACE) {
			return nil
		}

		// e.g. "y"
		expression.Alternative = p.parseBlockStatement()
	}

	if PRINT_PARSE {
		color.Blue("      RET p.parseIf(): %s", expression.String())
	}
	return expression
}

// Parse function expressions e.g. "fn(x, y) { x + y; }"
func (p *Parser) parseFunction() ast.Expression {
	if PRINT_PARSE {
		color.Cyan("      CALL p.parseFunction()")
	}
	// "fn"
	f := &ast.Function{Token: p.currentToken}

	// "("
	if !p.GetExpectNextToken(token.LPAREN) {
		return nil
	}

	// e.g. "x, y"
	f.Parameters = p.parseFunctionParameters()

	// e.g. "{"
	if !p.GetExpectNextToken(token.LBRACE) {
		return nil
	}

	f.Body = p.parseBlockStatement()

	if PRINT_PARSE {
		color.Blue("      RET p.parseFunction(): %s", f.String())
	}
	return f
}

// Helper method to parse function parameters
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	if PRINT_PARSE {
		color.Cyan("      CALL p.parseFunctionParameters()")
	}
	identifiers := []*ast.Identifier{}

	// Empty list of parameters: already ")"
	if p.nextToken.Type == token.RPAREN {
		p.GetNextToken()
		return identifiers
	}

	// First parameter
	p.GetNextToken()
	i := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
	identifiers = append(identifiers, i)

	// Other parameters
	for p.nextToken.Type == token.COMMA {
		p.GetNextToken()
		p.GetNextToken()

		i := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
		identifiers = append(identifiers, i)
	}

	// ")"
	if !p.GetExpectNextToken(token.RPAREN) {
		return nil
	}

	if PRINT_PARSE {
		color.Blue("      RET p.parseFunctionParameters()")
	}
	return identifiers
}

// Parse call expressions e.g. "add(1, 2);"
func (p *Parser) parseCall(function ast.Expression) ast.Expression {
	if PRINT_PARSE {
		color.Cyan("      CALL parseCall()")
	}

	c := &ast.Call{Token: p.currentToken, Function: function}
	c.Arguments = p.parseExpressionList(token.RPAREN)

	if PRINT_PARSE {
		color.Blue("      RET parseCall(): %s", c.String())
	}
	return c
}

// Parse string expressions
func (p *Parser) parseString() ast.Expression {
	return &ast.String{Token: p.currentToken, Value: p.currentToken.Literal}
}

// Parse array expressions
func (p *Parser) parseArray() ast.Expression {
	return &ast.Array{p.currentToken, p.parseExpressionList(token.RSQUARE)}
}

// Helper method to parse expression list
func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	// "]"
	if p.nextToken.Type == end {
		p.GetNextToken()
		return list
	}

	p.GetNextToken()

	list = append(list, p.parseExpression(LOWEST))

	for p.nextToken.Type == token.COMMA {
		p.GetNextToken()
		p.GetNextToken()

		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.GetExpectNextToken(end) {
		return nil
	} else {
		return list
	}
}

// Parse index expressions
func (p *Parser) parseIndex(array ast.Expression) ast.Expression {
	i := &ast.Index{Token: p.currentToken, Array: array}

	p.GetNextToken()

	i.Index = p.parseExpression(LOWEST)

	if !p.GetExpectNextToken(token.RSQUARE) {
		return nil
	} else {
		return i
	}
}

// Parse hash expressions
func (p *Parser) parseHash() ast.Expression {
	hash := &ast.Hash{Token: p.currentToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for p.nextToken.Type != token.RBRACE {
		// Get key
		p.GetNextToken()
		key := p.parseExpression(LOWEST)

		// ":"
		if !p.GetExpectNextToken(token.COLON) {
			return nil
		}

		// Get value
		p.GetNextToken()
		value := p.parseExpression(LOWEST)

		hash.Pairs[key] = value

		if p.nextToken.Type != token.RBRACE && !p.GetExpectNextToken(token.COMMA) {
			return nil
		}
	}

	if !p.GetExpectNextToken(token.RBRACE) {
		return nil
	} else {
		return hash
	}
}
