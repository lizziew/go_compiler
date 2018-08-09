package ast

import (
	"go_interpreter/token"
)

// Node of AST
type Node interface {
	TokenLiteral() string // for debugging
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

// Return Statement Node
type ReturnStatement struct {
	Token token.Token // token.RETURN
	Value Expression
}

func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
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
