package parser

import (
	"fmt"
	"mitchlang/token"
	"testing"

	"github.com/stretchr/testify/require"

	"mitchlang/ast"
	"mitchlang/lexer"
)

func TestParser_LetStatement(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	require.NotNil(t, program)
	require.Len(t, program.Statements, 3)
	require.Len(t, p.Errors(), 0, fmt.Sprintf("%#v", p.Errors()))

	tests := []struct {
		expected string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for k, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			stmt := program.Statements[k]
			require.Equal(t, stmt.TokenLiteral(), "let")
			require.IsType(t, &ast.LetStatement{}, stmt)
			let := stmt.(*ast.LetStatement)
			require.Equal(t, let.Name.Value, test.expected)
			require.Equal(t, let.Name.TokenLiteral(), test.expected)
		})
	}
}

func TestParser_LetStatementError(t *testing.T) {
	input := `
let x 5;
let = 10;
let 83838383;
let x = 10;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	require.NotNil(t, program)
	require.Greater(t, len(p.Errors()), 3)
}

func TestParser_ReturnStatement(t *testing.T) {
	input := `
return 3;
return add(15);
return 10;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	require.NotNil(t, program)
	require.Len(t, program.Statements, 3)
	require.Len(t, p.Errors(), 0, fmt.Sprintf("%#v", p.Errors()))

	tests := []struct {
		name string
	}{
		{"return int literal"},
		{"return function expression"},
		{"return int literal"},
	}

	for k, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stmt := program.Statements[k]
			require.Equal(t, stmt.TokenLiteral(), "return")
			require.IsType(t, &ast.ReturnStatement{}, stmt)
		})
	}
}

func TestParser_IdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	require.Len(t, program.Statements, 1)
	require.Len(t, p.Errors(), 0)
	require.IsType(t, &ast.ExpressionStatement{}, program.Statements[0])
	statement := program.Statements[0].(*ast.ExpressionStatement)
	require.IsType(t, &ast.Identifier{}, statement.Expression)
	ident := statement.Expression.(*ast.Identifier)
	require.Equal(t, ident.Value, "foobar")
	require.Equal(t, ident.TokenLiteral(), "foobar")
}

func TestParser_IntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	require.Len(t, program.Statements, 1)
	require.IsType(t, &ast.ExpressionStatement{}, program.Statements[0])
	statement := program.Statements[0].(*ast.ExpressionStatement)
	require.Equal(t, statement.Token.Type, token.Int)
	require.Equal(t, statement.Token.Literal, "5")
	require.IsType(t, &ast.IntegerLiteral{}, statement.Expression)
	literal := statement.Expression.(*ast.IntegerLiteral)
	require.Equal(t, literal.Value, int64(5))
	require.Equal(t, literal.TokenLiteral(), "5")
}

func TestParser_PrefixExpressions(t *testing.T) {
	tests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		require.Len(t, p.Errors(), 0)

		require.Len(t, program.Statements, 1)
		require.IsType(t, &ast.ExpressionStatement{}, program.Statements[0])
		statement := program.Statements[0].(*ast.ExpressionStatement)
		require.IsType(t, &ast.PrefixExpression{}, statement.Expression)
		expression := statement.Expression.(*ast.PrefixExpression)
		require.Equal(t, tt.operator, expression.Operator)
		require.IsType(t, &ast.IntegerLiteral{}, expression.Right)
		integer := expression.Right.(*ast.IntegerLiteral)
		require.Equal(t, tt.integerValue, integer.Value)
		require.Equal(t, integer.TokenLiteral(), fmt.Sprintf("%d", integer.Value))
	}
}

func TestParser_InfixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		left     int64
		right    int64
	}{
		{"5 + 5;", "+", 5, 5},
		{"5 - 5;", "-", 5, 5},
		{"5 * 5;", "*", 5, 5},
		{"5 / 5;", "/", 5, 5},
		{"5 > 5;", ">", 5, 5},
		{"5 < 5;", "<", 5, 5},
		{"5 == 5;", "==", 5, 5},
		{"5 != 5;", "!=", 5, 5},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		require.Len(t, p.Errors(), 0)

		require.Len(t, program.Statements, 1)
		require.IsType(t, &ast.ExpressionStatement{}, program.Statements[0])
		statement := program.Statements[0].(*ast.ExpressionStatement)
		require.IsType(t, &ast.InfixExpression{}, statement.Expression)
		expression := statement.Expression.(*ast.InfixExpression)
		require.Equal(t, expression.String(), fmt.Sprintf("(%d %s %d)", tt.left, tt.operator, tt.right))
		require.Equal(t, tt.operator, expression.Operator)
		require.IsType(t, &ast.IntegerLiteral{}, expression.Left)
		integerLeft := expression.Left.(*ast.IntegerLiteral)
		require.IsType(t, &ast.IntegerLiteral{}, expression.Right)
		require.Equal(t, tt.left, integerLeft.Value)
		integerRight := expression.Right.(*ast.IntegerLiteral)
		require.Equal(t, tt.right, integerRight.Value)
	}
}

func TestParser_OperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		require.Equal(t, tt.expected, program.String())
	}
}
