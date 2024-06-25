// ast/ast.go

package ast

import (
	"bytes"
	"monkey/token"
)

type Node interface {
	TokenLiteral() string
	String() string
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

func (p *Program) String() string {
	// Creates a buffer and writes the return value of each statement's String() method to it

	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type LetStatement struct {
	// Holds the LET token, the identifier, and the expression
	// let <name> = <value>; <=> let <identifer> = <expression>;
	// let x = 5; => holds: LET, Identifier(IDENT, "x"), and 5

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

func (ls *LetStatement) String() string {
	// Returns "let <name> = <value>;" as a string

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

type ReturnStatement struct {
	// Holds the RETURN token and the expression
	// return <expression>;

	Token       token.Token // The token.RETURN token
	ReturnValue Expression
}

// Implements the Statement interface
func (rs *ReturnStatement) statementNode() {}

func (rs *ReturnStatement) TokenLiteral() string {
	// Implements the Node interface

	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	// Returns "return <value>;" as a string

	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	// Holds the first token of an expression and the expression itself

	Token      token.Token // The first token of the expression
	Expression Expression
}

// Implements the Statement interface
func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) TokenLiteral() string {
	// Implements the Node interface

	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	// Returns the entire expression as a string

	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

type Identifier struct {
	// Holds the identifier of a binding
	// let x = 5; => holds: IDENT and "x"

	Token token.Token // The token.IDENT token
	Value string
}

// Implements the Expression interface
func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	// Implements the Node interface

	return i.Token.Literal
}

func (i *Identifier) String() string {
	// Returns the identifier as a string

	return i.Value
}

type IntegerLiteral struct {
	// Holds an integer literal
	// 5; => holds: INT and 5

	Token token.Token
	Value int64
}

// Implements the Expression interface
func (il *IntegerLiteral) expressionNode() {}

func (il *IntegerLiteral) TokenLiteral() string {
	// Implements the Node interface

	return il.Token.Literal
}

func (il *IntegerLiteral) String() string {
	// Returns the integer literal as a string

	return il.Token.Literal
}
