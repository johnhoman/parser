package parser

import (
    "fmt"
    "regexp"
    "strconv"
)

var specs = map[string]LiteralType{
    `(^\d+)`: IntLiteralType,
    `^"([^"]*)"`: StringLiteralType,
    `^\s+`: WhitespaceLiteral,
    `^//.*`: CommentLiteral,
    `^/\*[\s\S]*?\*/`: CommentLiteral,
}

// Lexical analysis

type SyntaxError struct {
    message string
}

func (err *SyntaxError) Error() string {
    return err.message
}

func NewSyntaxError(message string) error {
    return &SyntaxError{message}
}

type ExpressionStatement struct {}
type StatementList []ExpressionStatement


type Program struct {
    Type string
    Body StatementList
}

type Parser interface {
    // Parse parses a string into an abstract syntax tree (AST)
    Parse(s string) (Program, error)

    // Program
    //   ; StatementList
    //   ;
    Program() Program

    // StatementList
    //   : Statement
    //   | StatementList Statement
    //   ;
    StatementList() StatementList

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
func (p *parser) Parse(s string) (Program, error) {
    p.tokenizer = NewTokenizer(s)
    var err error
    p.lookAhead, err = p.tokenizer.NextToken()
    if err != nil {
        return Program{}, err
    }
    return p.Program(), nil
}

func (p *parser) Program() Program {
    return Program{
        Type: "Program",
        Body: p.Literal(),
    }
}

func (p *parser) eat(tokenType LiteralType) (Token, error) {
    token := p.lookAhead
    if token.IsEmpty() {
        return token, NewSyntaxError("")
    }
    var err error
    p.lookAhead, err = p.tokenizer.NextToken()
    if err != nil {
        return Token{}, err
    }
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
    NextToken() (Token, error)
}

type tokenizer struct {
    String
    cursor int
}

func (tok *tokenizer) hasMoreTokens() bool { return tok.cursor < tok.String.Len() }

func (tok *tokenizer) NextToken() (Token, error) {
    if !tok.hasMoreTokens() {
        return Token{}, nil
    }
    for pattern, literalType := range specs {
        str := string(tok.String.Slice(tok.cursor))

        re := regexp.MustCompile(pattern)
        if re.MatchString(str) {
            match := re.FindStringSubmatch(str)
            tok.cursor += len(match[0])
            if literalType == WhitespaceLiteral {
                return tok.NextToken()
            }
            if literalType == CommentLiteral {
                return tok.NextToken()
            }
            return Token{Type: string(literalType), Value: match[1]}, nil
        }
    }
    return Token{}, NewSyntaxError(fmt.Sprintf(`Unexpected token: %c`, tok.String[0]))
}

func NewTokenizer(s string) *tokenizer {
    return &tokenizer{String(s), 0}
}