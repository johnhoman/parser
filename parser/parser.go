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
	Equals      // ==
	LessGreater // > or <
	Sum         // +
	Product     // *
	Prefix      // -X or !X
	Call
)

var (
	precedences = map[token.Type]int{
		token.Eq:       Equals,
		token.NotEq:    Equals,
		token.LT:       LessGreater,
		token.GT:       LessGreater,
		token.Plus:     Sum,
		token.Minus:    Sum,
		token.Slash:    Product,
		token.Asterisk: Product,
		token.LParen:   Call,
	}
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

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.next.Type]; ok {
		return p
	}
	return Lowest
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.current.Type]; ok {
		return p
	}
	return Lowest
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
				tokenType,
				p.next.Type,
			),
		)
		return false
	}
	p.nextToken()
	return true
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: p.current}
	if !p.current.IsType(token.Return) {
		// This would be some internal error because it should never
		// get here
		return nil
	}
	p.nextToken()

	statement.ReturnValue = p.parseExpression(Lowest)
	if p.next.IsType(token.SemiColon) {
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

	p.nextToken()
	statement.Value = p.parseExpression(Lowest)
	// TODO: should this be required?
	if p.next.IsType(token.SemiColon) {
		p.nextToken()
	}
	return statement
}

func (p *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	call := &ast.CallExpression{Token: p.current, Function: left}
	call.Arguments = []ast.Expression{}
	p.nextToken()
	for !p.current.IsType(token.RParen) {
		arg := p.parseExpression(Lowest)
		if arg != nil {
			call.Arguments = append(call.Arguments, arg)
		}
		if p.next.IsType(token.Comma) {
			p.nextToken()
		}
		p.nextToken()
	}
	return call
}

// parseExpression parses an expression
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// let <identifier> = <prefix-operator | expression> <infix-operator> <expression>;
	// prefix can be anything
	prefix := p.prefixFuncs[p.current.Type]
	if prefix == nil {
		p.errors = append(
			p.errors,
			fmt.Sprintf("no prefix parse function for %s", p.current),
		)
		return nil
	}
	left := prefix()
	if left == nil {
		return nil
	}
	for !p.next.IsType(token.SemiColon) && precedence < p.peekPrecedence() {
		infix := p.infixFuncs[p.next.Type]
		if infix == nil {
			return left
		}
		p.nextToken()
		left = infix(left)
		if left == nil {
			return nil
		}
	}
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

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.current}
	block.Statements = []ast.Statement{}

	p.nextToken()
	for !p.current.IsType(token.RBrace) && !p.current.IsType(token.EOF) {
		statement := p.parseStatement()
		if statement != nil {
			block.Statements = append(block.Statements, statement)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.current.Type {
	case token.Let:
		if stmt := p.parseLetStatement(); stmt != nil {
			return stmt
		}
		return nil
	case token.Return:
		if stmt := p.parseReturnStatement(); stmt != nil {
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

func (p *Parser) parseString() ast.Expression {
	if !p.current.IsType(token.String) {
		return nil
	}
	return &ast.StringLiteral{Token: p.current, Value: p.current.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	if !p.current.IsType(token.True) && !p.current.IsType(token.False) {
		p.errors = append(
			p.errors,
			fmt.Sprintf("unexpected token %#v", p.current),
		)
		return nil
	}
	if p.current.IsType(token.True) {
		return &ast.Boolean{Token: p.current, Value: true}
	}
	return &ast.Boolean{Token: p.current, Value: false}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(Lowest)
	if exp == nil {
		return nil
	}
	if !p.expectNext(token.RParen) {
		return nil
	}
	return exp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.current,
		Operator: p.current.Literal,
	}
	p.nextToken()
	right := p.parseExpression(Prefix)
	if right == nil {
		return nil
	}
	expression.Right = right
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.current,
		Operator: p.current.Literal,
		Left:     left,
	}
	precedence := p.currentPrecedence()
	p.nextToken()
	right := p.parseExpression(precedence)
	if right == nil {
		return nil
	}
	expression.Right = right
	return expression
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.current}

	if !p.expectNext(token.LParen) {
		return nil
	}
	p.nextToken()
	expression.Condition = p.parseExpression(Lowest)
	if !p.expectNext(token.RParen) {
		return nil
	}
	if !p.expectNext(token.LBrace) {
		return nil
	}
	expression.Consequence = p.parseBlockStatement()

	if p.next.IsType(token.Else) {
		p.nextToken()
		if !p.expectNext(token.LBrace) {
			return nil
		}
		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseFunctionLiteralExpression() ast.Expression {
	expression := &ast.FunctionLiteralExpression{Token: p.current}
	expression.Parameters = []*ast.Identifier{}

	if !p.expectNext(token.LParen) {
		return nil
	}
	p.nextToken()
	for !p.current.IsType(token.RParen) {
		ident := &ast.Identifier{Token: p.current, Value: p.current.Literal}
		expression.Parameters = append(expression.Parameters, ident)
		if p.next.IsType(token.Comma) {
			p.nextToken()
		}
		p.nextToken()
	}
	if !p.expectNext(token.LBrace) {
		return nil
	}
	expression.Body = p.parseBlockStatement()
	return expression
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()

	p.prefixFuncs = make(map[token.Type]prefixFunc)
	p.registerPrefix(token.Ident, p.parseIdentifier)
	p.registerPrefix(token.Int, p.parseInteger)
	p.registerPrefix(token.String, p.parseString)
	p.registerPrefix(token.Bang, p.parsePrefixExpression)
	p.registerPrefix(token.Minus, p.parsePrefixExpression)
	p.registerPrefix(token.True, p.parseBoolean)
	p.registerPrefix(token.False, p.parseBoolean)
	p.registerPrefix(token.LParen, p.parseGroupedExpression)
	p.registerPrefix(token.If, p.parseIfExpression)
	p.registerPrefix(token.Function, p.parseFunctionLiteralExpression)

	p.infixFuncs = make(map[token.Type]infixFunc)
	p.registerInfix(token.Plus, p.parseInfixExpression)
	p.registerInfix(token.Minus, p.parseInfixExpression)
	p.registerInfix(token.Slash, p.parseInfixExpression)
	p.registerInfix(token.Asterisk, p.parseInfixExpression)
	p.registerInfix(token.Eq, p.parseInfixExpression)
	p.registerInfix(token.NotEq, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LParen, p.parseCallExpression)
	return p
}
