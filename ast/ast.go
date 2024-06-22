// ast/ast.go

package ast

import "monkey/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	// Root node of every AST

	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	// Returns the literal value of the token the node is associated with

	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	// Holds the LET token, the identifier, and the expression
	// let <name> = <value> <=> let <identifer> = <expression>
	// let x = 5 => holds: LET, Identifier(IDENT, "x"), and 5

	Token token.Token // The token.LET token
	Name  *Identifier
	Value Expression
}

// Implements the Statement interface
func (ls *LetStatement) statementNode() {}

func (ls *LetStatement) TokenLiteral() string {
	// Implements the Node interface
	return ls.Token.Literal
}

type Identifier struct {
	// Holds the identifier of a binding
	// let x = 5 => holds: IDENT and "x"

	Token token.Token // The token.IDENT token
	Value string
}

// Implements the Expression interface
func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	// Implements the Node interface
	return i.Token.Literal
}
