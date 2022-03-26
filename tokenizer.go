package parser

import (
    "errors"
    "strconv"
)

// Lexical analysis

var SyntaxError = errors.New("syntax error")

type Program struct {
    Type string
    Body interface{}
}

type Parser interface {
    // Parse parses a string into an abstract syntax tree (AST)
    Parse(s string) Program

    // Program
    //   ; Literal
    //   ;
    Program() Program

    // Literal
    //   ; StringLiteral
    //   ; NumericLiteral
    //   ;
    Literal() Literal

    // StringLiteral
    //   ; string
    StringLiteral() *StringLiteral

    // IntLiteral
    //   ; int
    //   ;
    IntLiteral() *IntLiteral
}

type parser struct {
    tokenizer Tokenizer
    lookAhead Token
}

// Parse the string s into an abstract syntax tree
func (p *parser) Parse(s string) Program {
    p.tokenizer = NewTokenizer(s)
    p.lookAhead = p.tokenizer.NextToken()
    return p.Program()
}

func (p *parser) Program() Program {
    return Program{
        Type: "Program",
        Body: repr(p.Literal()),
    }
}

func (p *parser) eat(tokenType LiteralType) (Token, error) {
    token := p.lookAhead
    if token.IsEmpty() {
        return token, SyntaxError
    }
    p.lookAhead = p.tokenizer.NextToken()
    return token, nil
}

func (p *parser) IntLiteral() *IntLiteral {
    token, _ := p.eat(IntLiteralType)
    i, _ := strconv.Atoi(token.Value)
    return &IntLiteral{i}
}

func (p *parser) StringLiteral() *StringLiteral {
    token, _ := p.eat(StringLiteralType)
    return &StringLiteral{token.Value}
}

func (p *parser) Literal() Literal {
    switch LiteralType(p.lookAhead.Type) {
    case StringLiteralType: { return p.StringLiteral() }
    case IntLiteralType: { return p.IntLiteral() }
    default: { return nil }
    }
}

var _ Parser = &parser{}

func New() *parser {
    return &parser{}
}

type Token struct {
    Type string
    Value string
}

func (tok *Token) IsEmpty() bool {
    return tok.Type == ""
}

type String string

func (s String) Len() int {
    return len(s)
}

func (s String) Slice(start int) String {
    return s[start:]
}

type Tokenizer interface {
    NextToken() Token
}

type tokenizer struct {
    String
    cursor int
}

func (tok *tokenizer) hasMoreTokens() bool { return tok.cursor < tok.String.Len() }

func (tok *tokenizer) NextToken() Token {
    if !tok.hasMoreTokens() {
        return Token{}
    }
    // Numbers
    str := tok.String.Slice(tok.cursor)
    if '0' <= str[0] && str[0] <= '9' {
        k := 0
        integer := make([]byte, 0)
        for k < str.Len() && '0' <= str[k] && str[k] <= '9' {
            integer = append(integer, str[k])
            k++
        }
        tok.cursor += k
        return Token{
            Type: string(IntLiteralType),
            Value: string(integer),
        }
    }
    if str[0] == '"' {
        k := 1
        start := k
        for k < str.Len() && str[k] != '"' {
            k++
        }
        tok.cursor += k + 1
        return Token{
            Type: string(StringLiteralType),
            Value: string(str[start:k]),
        }
    }
    return Token{}
}

func NewTokenizer(s string) *tokenizer {
    return &tokenizer{String(s), 0}
}