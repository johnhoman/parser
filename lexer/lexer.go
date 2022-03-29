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

func (l *Lexer) NextToken() token.Token {

    tok := token.Token{}
    switch l.ch {
    case 0: { tok.Type = token.EOF }
    }
}

func New(input string) *Lexer {
    lex := &Lexer{input: input}
    lex.readChar()
    return lex
}