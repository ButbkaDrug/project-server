package parser

import (
	"strings"
	"unicode"
)

type system int

const (
	INT    = "INT"
	FLOAT  = "FLOAT"
	STRING = "STRING"
	EOF    = "EOF"

	LINUX system = iota
	WIN
	UNKNOWN
)

type TokenType string

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()

	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

func (l *Lexer) readString() string {
	position := l.position

	for unicode.IsLetter(rune(l.ch)) || l.ch == ' ' {
		l.readChar()
	}

	return l.input[position:l.position]
}

// Will read any chars until encounters stop char
func (l *Lexer) readAny(stop byte) string {
	position := l.position

	for l.ch != stop && l.ch != 0 {

		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position

	for unicode.IsNumber(rune(l.ch)) || l.ch == '.' {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) isFloat(s string) bool {
	return strings.Index(s, ".") != -1
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '"':
		l.readChar()
		str := l.readAny('"')
		tok = Token{Type: STRING, Value: str}
	case '{':
		str := l.readAny('}')
		tok = Token{Type: STRING, Value: str}
	case 0:
		return Token{Type: EOF, Value: ""}
	default:
		if unicode.IsDigit(rune(l.ch)) {
			n := l.readNumber()
			tok = Token{Value: n}

			if l.isFloat(n) {
				tok.Type = FLOAT
			} else {
				tok.Type = INT
			}
			return tok
		} else if unicode.IsLetter(rune(l.ch)) {
			str := l.readString()
			return Token{Type: STRING, Value: str}
		}

		return Token{Type: EOF, Value: "0"}
	}

	l.readChar()

	return tok
}
