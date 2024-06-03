// lexer/lexer.go

package lexer

import "monkey/token"

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	// create a new Lexer and read the first char

	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	// give the next char and advance the cursor position

	if l.readPosition >= len(l.input) {
		// ASCII code for NULL is 0
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	// advance the current position
	l.position = l.readPosition
	// advance the reading position
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	// read the current char and return its corresponding token after advancing the cursor

	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			// save l.ch in a local variable before calling l.readChar() again so we don't lose the
			// current character
			ch := l.ch
			l.readChar()
			// concatenate the current assignment operator `=` and the subsequent `=`
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			// save l.ch in a local variable before calling l.readChar() again so we don't lose the
			// current character
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			// concatenate the current bang operator `!` and the subsequent `=`
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	// create a new token

	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	// read in an identifier and advance the lexer's position until encountering a non-letter char

	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	// check if the char falls within the ASCII code tables for valid letters, the code tables from
	// a-z and A-Z are sequential

	// `a`: 01100001, `z`: 01111010
	// `A`: 01000001, `Z`: 01011010
	// `_`: 01011111
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) skipWhitespace() {
	// skips spaces, tabs, newlines, and carriage returns

	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readNumber() string {
	// read in a number and advance the lexer's position until encountering a non-digit char

	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	// check if the char falls within the ASCII code tables for valid numbers, the code tables from
	// 0-9 are sequential

	// `0`: 00110000, `9`: 00111001
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) peekChar() byte {
	// look ahead by one char and return it; similar to readChar() without incrementing the cursor

	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}
