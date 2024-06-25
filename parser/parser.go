// parser/parser.go

package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

const (
	// Operator precedences

	// The iota keyword gives the following constants incrementing numbers as values; The blank
	// identifier _ takes the zero value and the following constants get assigned the values 1 to 7
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // < or >
	SUM         // +
	PRODUCT     // *
	PREFIX      // -x or !x
	CALL        // myFunction(x)
)

type Parser struct {
	// The parser implementation

	l *lexer.Lexer

	// Slice of strings to hold error messages
	errors []string

	// These act like the two pointers that the lexer has, but instead of pointing to chars in the
	// input, they point to tokens
	curToken  token.Token
	peekToken token.Token

	// Used to check if the appropriate map (prefix or infix) has a parsing function associated with
	// curToken.Type
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression

	// The argument being passed is the "left side" of the infix operator, e.g. the 5 in `5 + 6`
	infixParseFn func(ast.Expression) ast.Expression
)

func New(l *lexer.Lexer) *Parser {
	// Creates a new parser

	p := &Parser{l: l, errors: []string{}}

	// Initialize the prefix parse function map and register a parsing function
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	// Returns parser errors to check if any were encountered

	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	// Adds a new error to the parser when the next token is not as expected

	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	// Advances curToken and peekToken

	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	// Constructs the root node of the AST, iterates over every token in the input, and returns the
	// root node when finished

	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	// Parses a statement based on its corresponding token

	// The only two pure statement types in monkey are `let` and `return` statements, so if they
	// aren't encountered, the statement must be an expression
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	// Parses an expression based on its operator precedence

	// Check if there is a parsing function associated with the current token type in the prefix
	// position
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		return nil
	}

	leftExp := prefix()

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	// Returns an identifier with the current token and the current token literal

	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	// Constructs an *ast.LetStatement node with a LET token
	// let <identifer> = <expression>;

	stmt := &ast.LetStatement{Token: p.curToken}

	// Ensure the identifier exists
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// Ensure the assignment operator exists
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: 06/19/24 - For now, we're skipping the expressions until we encounter a semicolon

	// Ensure the line ends
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	// Constructs an *ast.ReturnStatement node with a RETURN token
	// return <expression>;

	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// TODO: 06/22/24 - For now, we're skipping the expressions until we encounter a semicolon

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// Constructs an *ast.ExpressionStatement node with an expression statement

	stmt := &ast.ExpressionStatement{Token: p.curToken}

	// Parse the expression starting with the lowest operator precedence
	stmt.Expression = p.parseExpression(LOWEST)

	// Check for an optional semicolon
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	// COnstructs an *ast.IntegerLiteral node with an integer literal

	lit := &ast.IntegerLiteral{Token: p.curToken}

	// Convert the integer literal string into an int64
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	// Checks if the current token is of type `t`

	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	// Checks if the next token is of type `t`

	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	// Checks the type of peekToken and advances if the type is as expected, otherwise logs a
	// peekError

	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	// Adds a function entry to the prefix map

	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	// Adds a function entry to the infix map

	p.infixParseFns[tokenType] = fn
}
