package parser

import (
    "errors"
    "strconv"
)

// Lexical analysis

var SyntaxError = errors.New("syntax error")

type NumericLiteral struct {
    Type string
    Value int
}

type Program struct {
    Type string
    Body interface{}
}

type Parser interface {
    // Parse parses a string into an abstract syntax tree (AST)
    Parse(s string) Program

    // Program
    //   ; NumericLiteral
    //   ;
    Program() Program
    // NumericLiteral
    //   ; int
    //   ;
    NumericLiteral() NumericLiteral
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
        Body: p.NumericLiteral(),
    }
}

func (p *parser) eat(tokenType string) (Token, error) {
    token := p.lookAhead
    if token.IsEmpty() {
        return token, SyntaxError
    }
    p.lookAhead = p.tokenizer.NextToken()
    return token, nil
}

func (p *parser) NumericLiteral() NumericLiteral {
    token, _ := p.eat("NUMBER")
    i, _ := strconv.Atoi(token.Value)
    return NumericLiteral{"NumericLiteral", i}
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
    k := 0
    numbers := make([]byte, 0)
    for k < str.Len() {
        if '0' <= str[k] && str[k] <= '9' {
            numbers = append(numbers, str[k])
            k++
        } else {
            break
        }
    }
    return Token{
        Type: "NUMBER",
        Value: string(numbers),
    }
}

func NewTokenizer(s string) *tokenizer {
    return &tokenizer{String(s), 0}
}