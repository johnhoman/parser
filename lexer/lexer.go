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

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) readWhen(pred func(byte) bool) string {
	if l.position >= len(l.input) {
		// EOF
		return string(byte(0))
	}
	position := l.position

	for pred(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	s := l.readWhen(func(b byte) bool {
		return b != '"' && b != 0
	})
	return s
}

func (l *Lexer) readIdentifier() string {
	return l.readWhen(isLetter)
}

func (l *Lexer) readInt() string {
	return l.readWhen(isDigit)
}

func (l *Lexer) readWhitespace() string {
	return l.readWhen(isWhitespace)
}

func (l *Lexer) NextToken() *token.Token {

	_ = l.readWhitespace()

	var tok *token.Token
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			tok = token.New(token.Eq, l.ch, l.peekChar())
			l.readChar()
		} else {
			tok = token.New(token.Assign, l.ch)
		}
	case '+':
		tok = token.New(token.Plus, l.ch)
	case '-':
		tok = token.New(token.Minus, l.ch)
	case '!':
		if l.peekChar() == '=' {
			tok = token.New(token.NotEq, l.ch, l.peekChar())
			l.readChar()
		} else {
			tok = token.New(token.Bang, l.ch)
		}
	case '*':
		tok = token.New(token.Asterisk, l.ch)
	case '/':
		tok = token.New(token.Slash, l.ch)
	case '<':
		tok = token.New(token.LT, l.ch)
	case '>':
		tok = token.New(token.GT, l.ch)
	case ',':
		tok = token.New(token.Comma, l.ch)
	case ';':
		tok = token.New(token.SemiColon, l.ch)
	case '(':
		tok = token.New(token.LParen, l.ch)
	case ')':
		tok = token.New(token.RParen, l.ch)
	case '{':
		tok = token.New(token.LBrace, l.ch)
	case '}':
		tok = token.New(token.RBrace, l.ch)
	case '[':
		tok = token.New(token.LBracket, l.ch)
	case ']':
		tok = token.New(token.RBracket, l.ch)
	case '"':
		l.readChar()
		tok = token.NewFromString(token.String, l.readString())
	case 0:
		tok = token.New(token.EOF)
	default:
		{
			if isLetter(l.ch) {
				tok = token.NewIdentifier(l.readIdentifier())
			} else if isDigit(l.ch) {
				tok = token.NewFromString(token.Int, l.readInt())
			} else {
				tok = token.New(token.Illegal, l.ch)
			}
			return tok
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

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isAsciiLower(ch byte) bool { return 'a' <= ch && ch <= 'z' }
func isAsciiUpper(ch byte) bool { return 'A' <= ch && ch <= 'Z' }
func isLetter(ch byte) bool {
	return isAsciiLower(ch) || isAsciiUpper(ch) || ch == '_'
}
func isDigit(ch byte) bool { return '0' <= ch && ch <= '9' }
