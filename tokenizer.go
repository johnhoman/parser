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
    `^;`: ExpressionTerm,
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

type Expression struct {
    Literal Literal
}

type ExpressionStatement struct {
    Expression Expression
}

type Statement struct {
    ExpressionStatement ExpressionStatement
}

type StatementList struct {
    Statements []Statement
    StatementList *StatementList
}

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
    Program() (Program, error)

    // StatementList
    //   : Statement
    //   | StatementList Statement
    //   ;
    StatementList() (StatementList, error)

    // Statement
    //   : ExpressionStatement
    //   ;
    Statement() (Statement, error)

    // ExpressionStatement
    //   : Expression ;
    //   ;
    ExpressionStatement() (ExpressionStatement, error)

    // Expression
    //   : Literal
    //   ;
    Expression() (Expression, error)

    // Literal
    //   : StringLiteral
    //   | NumericLiteral
    //   ;
    Literal() (Literal, error)

    // StringLiteral
    //   ; string
    StringLiteral() (*StringLiteral, error)

    // IntLiteral
    //   ; int
    //   ;
    IntLiteral() (*IntLiteral, error)
}

type parser struct {
    tokenizer *tokenizer
    lookAhead Token
}

// Parse the string s into an abstract syntax tree
func (p *parser) Parse(s string) (Program, error) {
    p.tokenizer = NewTokenizer(s)
    p.lookAhead, _ = p.tokenizer.NextToken()
    return p.Program()
}

func (p *parser) Expression() (Expression, error) {
    lit, err := p.Literal()
    if err != nil {
        return Expression{}, err
    }
    _, err = p.eat(ExpressionTerm)
    if err != nil {
        return Expression{}, err
    }
    return Expression{Literal: lit}, nil
}

func (p *parser) ExpressionStatement() (ExpressionStatement, error) {
    expression, err := p.Expression()
    if err != nil {
        return ExpressionStatement{}, err
    }
    return ExpressionStatement{Expression: expression}, nil
}

func (p *parser) Statement() (Statement, error) {
    expressionStatement, err := p.ExpressionStatement()
    if err != nil {
        return Statement{}, err
    }
    return Statement{ExpressionStatement: expressionStatement}, nil
}

func (p *parser) StatementList() (StatementList, error) {
    statements := StatementList{}

    for !p.lookAhead.IsEmpty() {
        statement, err := p.Statement()
        if err != nil {
            return StatementList{}, err
        }
        statements.Statements = append(statements.Statements, statement)
    }
    return statements, nil
}

func (p *parser) Program() (Program, error) {
    statements, err := p.StatementList()
    if err != nil {
        return Program{}, err
    }
    return Program{Type: "Program", Body: statements}, nil
}

func (p *parser) eat(tokenType LiteralType) (Token, error) {
    token := p.lookAhead
    if token.IsEmpty() {
        return token, NewSyntaxError("EOF")
    }
    if token.Type != string(tokenType) {
        return token, NewSyntaxError(fmt.Sprintf("unexpected token: '%s'", token.Value))
    }
    p.lookAhead, _ = p.tokenizer.NextToken()
    return token, nil
}

func (p *parser) IntLiteral() (*IntLiteral, error) {
    token, err := p.eat(IntLiteralType)
    if err != nil {
        return nil, err
    }
    i, _ := strconv.Atoi(token.Value)
    return &IntLiteral{i}, nil
}

func (p *parser) StringLiteral() (*StringLiteral, error) {
    token, err := p.eat(StringLiteralType)
    if err != nil {
        return nil, err
    }
    return &StringLiteral{token.Value}, nil
}

func (p *parser) Literal() (Literal, error) {
    switch LiteralType(p.lookAhead.Type) {
    case StringLiteralType: { return p.StringLiteral() }
    case IntLiteralType: { return p.IntLiteral() }
    default: { return nil, NewSyntaxError(fmt.Sprintf("Invalid literal type %s", p.lookAhead.Type)) }
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
            if literalType == ExpressionTerm {
                return Token{Type: string(literalType)}, nil
            }
            return Token{Type: string(literalType), Value: match[1]}, nil
        }
    }
    return Token{}, NewSyntaxError(fmt.Sprintf(`Unexpected token: %c`, tok.String[0]))
}

func NewTokenizer(s string) *tokenizer {
    return &tokenizer{String(s), 0}
}