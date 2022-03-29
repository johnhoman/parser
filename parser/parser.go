package parser

import (
	"fmt"
	"mitchlang/ast"
	"mitchlang/lexer"
	"mitchlang/token"
	"strconv"
)

const (
	_ int = iota
	Lowest
	Equals
	LessGreater
	Sum
	Product
	Prefix
	Call
)

type (
	prefixFunc func() ast.Expression
	infixFunc  func(ast.Expression) ast.Expression
)

type Parser struct {
	l       *lexer.Lexer
	current *token.Token
	next    *token.Token
	errors  []string

	prefixFuncs map[token.Type]prefixFunc
	infixFuncs  map[token.Type]infixFunc
}

func (p *Parser) registerPrefix(tType token.Type, fn prefixFunc) {
	p.prefixFuncs[tType] = fn
}

func (p *Parser) registerInfix(tType token.Type, fn infixFunc) {
	p.infixFuncs[tType] = fn
}

func (p *Parser) nextToken() {
	p.current = p.next
	p.next = p.l.NextToken()
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

func (p *Parser) parserReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: p.current}

	for !p.current.IsType(token.SemiColon) {
		p.nextToken()
	}
	return statement
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

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixFuncs[p.current.Type]
	if prefix == nil {
		p.errors = append(
			p.errors,
			fmt.Sprintf("no prefix parse function for %s", p.current),
		)
		return nil
	}
	left := prefix()
	return left
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{Token: p.current}

	stmt.Expression = p.parseExpression(Lowest)

	if p.next.IsType(token.SemiColon) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.current.Type {
	case token.Let:
		if stmt := p.parseLetStatement(); stmt != nil {
			return stmt
		}
		return nil
	case token.Return:
		if stmt := p.parserReturnStatement(); stmt != nil {
			return stmt
		}
		return nil
	default:
		if stmt := p.parseExpressionStatement(); stmt != nil {
			return stmt
		}
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

func (p *Parser) parseIdentifier() ast.Expression {
	if !p.current.IsType(token.Ident) {
		return nil
	}
	return &ast.Identifier{Token: p.current, Value: p.current.Literal}
}

func (p *Parser) parseInteger() ast.Expression {
	if !p.current.IsType(token.Int) {
		return nil
	}
	v, err := strconv.ParseInt(p.current.Literal, 10, 64)
	if err != nil {
		p.errors = append(
			p.errors,
			fmt.Sprintf("coun tnot parse %q as integer", p.current.Literal),
		)
		return nil
	}
	return &ast.IntegerLiteral{Token: p.current, Value: v}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token: p.current,
		Operator: p.current.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(Prefix)
	return expression
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()

	p.prefixFuncs = make(map[token.Type]prefixFunc)
	p.registerPrefix(token.Ident, p.parseIdentifier)
	p.registerPrefix(token.Int, p.parseInteger)
	p.registerPrefix(token.Bang, p.parsePrefixExpression)
	p.registerPrefix(token.Minus, p.parsePrefixExpression)
	return p
}