package lexer

import (
	"monkey/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte

	context token.TokenType
	Line    int
	Col     int
	column  int
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.context = token.EOF
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	l.ch = l.peekChar()
	l.position = l.readPosition
	l.readPosition += 1

	if l.ch == '\n' {
		l.Line++
		l.Col = l.column
		l.column = 0
	} else {
		l.column++
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) preEqual(pre token.TokenType, alone token.TokenType) token.Token {
	if l.peekChar() == '=' {
		ch := l.ch
		l.readChar()
		return token.Token{Type: pre, Literal: string(ch) + string(l.ch)}
	} else {
		return newToken(alone, l.ch)
	}
}

func (l *Lexer) readTemplateIdent() token.Token {
	var tok token.Token

	if l.context == token.TEMPLATE {
		l.readChar()
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.IDENT

			if l.ch == '`' {
				l.readChar()
				l.context = token.EOF
			}

			return tok
		}
	}
	return newToken(token.ILLEGAL, l.ch)
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	if l.context == token.EOF {
		l.skipWhitespace()
	}

	if l.context == token.TEMPLATE {
		if l.ch == '$' {
			return l.readTemplateIdent()
		} else {
			tok.Type = token.TEMPLATE
			tok.Literal = l.readTemplate(0, string(l.ch))
		}
	} else {
		switch l.ch {
		case '=':
			tok = l.preEqual(token.EQ, token.ASSIGN)
		case '!':
			tok = l.preEqual(token.NE, token.BANG)
		case '<':
			tok = l.preEqual(token.LE, token.LT)
		case '>':
			tok = l.preEqual(token.GE, token.GT)
		case '[':
			tok = newToken(token.LBRACKET, l.ch)
		case ']':
			tok = newToken(token.RBRACKET, l.ch)
		case '}':
			tok = newToken(token.RBRACE, l.ch)
		case '{':
			tok = newToken(token.LBRACE, l.ch)
		case ')':
			tok = newToken(token.RPAREN, l.ch)
		case '(':
			tok = newToken(token.LPAREN, l.ch)
		case ';':
			tok = newToken(token.SEMICOLON, l.ch)
		case ',':
			tok = newToken(token.COMMA, l.ch)
		case '+':
			tok = newToken(token.PLUS, l.ch)
		case '-':
			tok = newToken(token.MINUS, l.ch)
		case '*':
			tok = newToken(token.ASTERISK, l.ch)
		case '/':
			tok = newToken(token.SLASH, l.ch)
		case '%':
			tok = newToken(token.PERCENT, l.ch)
		case ':':
			tok = newToken(token.COLON, l.ch)
		case '&':
			tok = newToken(token.AND, l.ch)
		case '|':
			tok = newToken(token.OR, l.ch)
		case '"':
			tok.Type = token.STRING
			tok.Literal = l.readString()
		case '`':
			if l.context != token.TEMPLATE {
				l.context = token.TEMPLATE
				tok.Type = token.TEMPLATE
				tok.Literal = l.readTemplate(1, "")
			}
		case 0:
			tok.Literal = ""
			tok.Type = token.EOF
		default:
			if isLetter(l.ch) {
				tok.Literal = l.readIdentifier()
				tok.Type = token.LookupIdent(tok.Literal)
				return tok
			} else if isNumber(l.ch) {
				tok.Literal = l.readNumber()
				tok.Type = token.INT
				return tok
			} else {
				tok = newToken(token.ILLEGAL, l.ch)
			}
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isNumber(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	l.readChar()
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readTemplate(offset int, skip string) string {
	if l.peekChar() == '$' {
		return skip
	}

	position := l.position + offset
	l.readChar()

	for l.ch != '`' && l.peekChar() != '$' && l.ch != 0 {
		l.readChar()
	}

	if l.ch == '`' {
		l.context = token.EOF
		return l.input[position:l.position]
	}

	return l.input[position : l.position+1]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isNumber(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipComments() {
	// Handle comments lol
	if l.ch == '/' && l.peekChar() == '/' {
		for l.ch != '\n' && l.ch != 0 {
			l.readChar()
		}
	}
}

func (l *Lexer) skipWhitespace() {
	l.skipComments()
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
		l.skipComments()
	}
}

func isLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func isNumber(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
