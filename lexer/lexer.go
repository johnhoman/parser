package lexer

import (
	"mitchlang/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
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

func (l *Lexer) NextToken() *token.Token {

	var tok *token.Token
	switch l.ch {
	case '+':
		{
			tok = token.New(token.Plus, string(l.ch))
		}
	case '=':
		{
			tok = token.New(token.Assign, string(l.ch))
		}
	case ';':
		{
			tok = token.New(token.SemiColon, string(l.ch))
		}
	case ',':
		{
			tok = token.New(token.Comma, string(l.ch))
		}
	case '(':
		{
			tok = token.New(token.LParen, string(l.ch))
		}
	case ')':
		{
			tok = token.New(token.RParen, string(l.ch))
		}
	case '{':
		{
			tok = token.New(token.LBrace, string(l.ch))
		}
	case '}':
		{
			tok = token.New(token.RBrace, string(l.ch))
		}
	case 0:
		{
			tok = token.New(token.EOF, "")
		}
	}
    l.readChar()
    return tok
}

func New(input string) *Lexer {
	lex := &Lexer{input: input}
	lex.readChar()
	return lex
}
