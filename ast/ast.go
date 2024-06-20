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
	// root node of every AST

	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	// returns the literal value of the token the node is associated with

	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	// holds the LET token, the identifier, and the expression
	// let <name> = <value> <=> let <identifer> = <expression>
	// let x = 5 => holds: LET, Identifier(IDENT, "x"), and 5

	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

// implements Statement interface
func (ls *LetStatement) statementNode() {}

// implements Node interface
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type Identifier struct {
	// holds the identifier of a binding
	// let x = 5 => holds: IDENT and "x"

	Token token.Token // the token.IDENT token
	Value string
}

// implements Expression interface
func (i *Identifier) expressionNode() {}

// implements Node interface
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
