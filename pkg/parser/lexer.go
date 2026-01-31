package parser

import (
	"strings"
	"unicode"
)

// TokenType represents the type of token
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenKeyword
	TokenIdentifier
	TokenNumber
	TokenString
	TokenComma
	TokenSemicolon
	TokenLParen
	TokenRParen
	TokenStar
	TokenEquals
	TokenNotEquals
	TokenLessThan
	TokenGreaterThan
	TokenLessEquals
	TokenGreaterEquals
)

// Token represents a lexical token
type Token struct {
	Type    TokenType
	Literal string
}

// Lexer tokenizes SQL input
type Lexer struct {
	input   string
	pos     int
	readPos int
	ch      byte
}

// NewLexer creates a new lexer
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// NextToken returns the next token
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	var tok Token

	switch l.ch {
	case 0:
		tok = Token{Type: TokenEOF, Literal: ""}
	case ',':
		tok = Token{Type: TokenComma, Literal: ","}
		l.readChar()
	case ';':
		tok = Token{Type: TokenSemicolon, Literal: ";"}
		l.readChar()
	case '(':
		tok = Token{Type: TokenLParen, Literal: "("}
		l.readChar()
	case ')':
		tok = Token{Type: TokenRParen, Literal: ")"}
		l.readChar()
	case '*':
		tok = Token{Type: TokenStar, Literal: "*"}
		l.readChar()
	case '=':
		tok = Token{Type: TokenEquals, Literal: "="}
		l.readChar()
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TokenLessEquals, Literal: "<="}
		} else if l.peekChar() == '>' {
			l.readChar()
			tok = Token{Type: TokenNotEquals, Literal: "<>"}
		} else {
			tok = Token{Type: TokenLessThan, Literal: "<"}
		}
		l.readChar()
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TokenGreaterEquals, Literal: ">="}
		} else {
			tok = Token{Type: TokenGreaterThan, Literal: ">"}
		}
		l.readChar()
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TokenNotEquals, Literal: "!="}
			l.readChar()
		}
	case '\'':
		tok = Token{Type: TokenString, Literal: l.readString()}
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			tokType := TokenIdentifier
			if isKeyword(literal) {
				tokType = TokenKeyword
			}
			tok = Token{Type: tokType, Literal: literal}
		} else if isDigit(l.ch) {
			tok = Token{Type: TokenNumber, Literal: l.readNumber()}
		} else {
			tok = Token{Type: TokenEOF, Literal: ""}
		}
	}

	return tok
}

func (l *Lexer) readIdentifier() string {
	pos := l.pos
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readNumber() string {
	pos := l.pos
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readString() string {
	l.readChar() // skip opening quote
	pos := l.pos
	for l.ch != '\'' && l.ch != 0 {
		l.readChar()
	}
	str := l.input[pos:l.pos]
	l.readChar() // skip closing quote
	return str
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch))
}

func isDigit(ch byte) bool {
	return unicode.IsDigit(rune(ch))
}

var keywords = map[string]bool{
	"SELECT": true, "FROM": true, "WHERE": true, "INSERT": true,
	"INTO": true, "VALUES": true, "CREATE": true, "TABLE": true,
	"AND": true, "OR": true, "NOT": true, "NULL": true,
	"TRUE": true, "FALSE": true, "INTEGER": true, "INT": true,
	"TEXT": true, "VARCHAR": true, "STRING": true, "BOOLEAN": true,
	"BOOL": true, "DELETE": true, "UPDATE": true, "SET": true,
	"BEGIN": true, "COMMIT": true, "ROLLBACK": true, "TRANSACTION": true,
}

func isKeyword(s string) bool {
	return keywords[strings.ToUpper(s)]
}
