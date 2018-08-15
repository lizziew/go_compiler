package ast

import (
	"bytes"
	"go_interpreter/token"
	"strings"
)

// Node of AST
type Node interface {
	TokenLiteral() string // for debugging
	String() string       // for debugging
}

// Statement type for Node
type Statement interface {
	Node
	statementNode()
}

// Expression type for Node
type Expression interface {
	Node
	expressionNode()
}

// Program Node (AST root)
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// Let Statement Node
type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

// Return Statement Node
type ReturnStatement struct {
	Token token.Token // token.RETURN
	Value Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.Value != nil {
		out.WriteString(rs.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

// Expression Statement Node
// Wrapper: Statement consists of 1 Expression
// e.g. "x + 5;" is valid
type ExpressionStatement struct {
	Token      token.Token // first token of expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

// Block Statement Node
type BlockStatement struct {
	Token      token.Token // { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// Identifier Expression Node
type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

// Integer Literal Expression Node
type IntegerLiteral struct {
	Token token.Token // token.INT
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

// Prefix Expression Node
type Prefix struct {
	Token    token.Token // prefix token e.g. "!"
	Operator string
	Value    Expression
}

func (p *Prefix) expressionNode() {}

func (p *Prefix) TokenLiteral() string {
	return p.Token.Literal
}

func (p *Prefix) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Value.String())
	out.WriteString(")")

	return out.String()
}

// Infix Expression Node
type Infix struct {
	Token    token.Token // infix token e.g. "+"
	Left     Expression
	Operator string
	Right    Expression
}

func (i *Infix) expressionNode() {}

func (i *Infix) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Infix) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")

	return out.String()
}

// Boolean Expression Node
type Boolean struct {
	Token token.Token // token.TRUE, token.FALSE
	Value bool
}

func (b *Boolean) expressionNode() {}

func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return b.Token.Literal
}

// If Expression Node
type If struct {
	Token       token.Token // token.IF
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i *If) expressionNode() {}

func (i *If) TokenLiteral() string {
	return i.Token.Literal
}

func (i *If) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(i.Condition.String())
	out.WriteString(" ")
	out.WriteString(i.Consequence.String())

	if i.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(i.Alternative.String())
	}

	return out.String()
}

// Function Expression Node
type Function struct {
	Token      token.Token // token.FUNCTION
	Parameters []*Identifier
	Body       *BlockStatement
}

func (f *Function) expressionNode() {}

func (f *Function) TokenLiteral() string {
	return f.Token.Literal
}

func (f *Function) String() string {
	var out bytes.Buffer

	parameters := []string{}
	for _, p := range f.Parameters {
		parameters = append(parameters, p.String())
	}

	out.WriteString(f.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(parameters, ", "))
	out.WriteString(")")

	out.WriteString(f.Body.String())

	return out.String()
}

// Call Expression Node
type Call struct {
	Token     token.Token // token.LPAREN
	Function  Expression  // Identifier or Function Node
	Arguments []Expression
}

func (c *Call) expressionNode() {}

func (c *Call) TokenLiteral() string {
	return c.Token.Literal
}

func (c *Call) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range c.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(c.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// String Expression Node
type String struct {
	Token token.Token
	Value string
}

func (s *String) expressionNode() {}

func (s *String) TokenLiteral() string {
	return s.Token.Literal
}

func (s *String) String() string {
	return s.Token.Literal
}
