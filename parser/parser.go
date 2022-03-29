package parser

import (
    "fmt"
    "mitchlang/ast"
    "mitchlang/lexer"
    "mitchlang/token"
)

type Parser struct {
    l *lexer.Lexer
    current *token.Token
    next    *token.Token
    errors  []string
}

func (p *Parser) nextToken() {
    p.current = p.next
    p.next    = p.l.NextToken()
}

func (p *Parser) expectNext(tokenType token.Type) bool {
    if !p.next.IsType(tokenType) {
        p.errors = append(
            p.errors,
            fmt.Sprintf(
                "expected next token to be %s, got %s instead",
                token.Assign,
                p.next.Type,
            ),
        )
        return false
    }
    p.nextToken()
    return true
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
    statement := &ast.LetStatement{Token: p.current}

    // let x
    if !p.expectNext(token.Ident) {
        return nil
    }

    statement.Name = &ast.Identifier{Token: p.current, Value: p.current.Literal}

    // let x =
    if !p.expectNext(token.Assign) {
        return nil
    }

    for !p.current.IsType(token.SemiColon) {
        p.nextToken()
    }
    return statement
}

func (p *Parser) parseStatement() ast.Statement {
    switch p.current.Type {
    case token.Let:
        if stmt := p.parseLetStatement(); stmt != nil {
            return stmt
        }
        return nil
    default:
        return nil
    }
}

func (p *Parser) ParseProgram() *ast.Program {
    program := &ast.Program{}
    program.Statements = []ast.Statement{}

    for !p.current.IsType(token.EOF) {
        statement := p.parseStatement()
        if statement != nil {
            program.Statements = append(program.Statements, statement)
        }
        p.nextToken()
    }
    return program
}

func (p *Parser) Errors() []string {
    return p.errors
}

func New(l *lexer.Lexer) *Parser {
    p := &Parser{l: l}
    p.nextToken()
    p.nextToken()
    return p
}
