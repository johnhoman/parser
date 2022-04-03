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
		expected      string
		expectedValue int64
	}{
		{"x", 5},
		{"y", 10},
		{"foobar", 838383},
	}

	for k, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			stmt := program.Statements[k]
			require.Equal(t, stmt.TokenLiteral(), "let")
			require.IsType(t, &ast.LetStatement{}, stmt)
			let := stmt.(*ast.LetStatement)
			require.Equal(t, let.Name.Value, test.expected)
			require.Equal(t, let.Name.TokenLiteral(), test.expected)
			require.IsType(t, &ast.IntegerLiteral{}, let.Value)
			lit := let.Value.(*ast.IntegerLiteral)
			require.Equal(t, test.expectedValue, lit.Value)
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
	checkErrors(t, p.Errors())
	require.Len(t, program.Statements, 3)

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

func TestParser_StringLiteralExpression(t *testing.T) {
	input := `"this is a string"`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	for _, err := range p.Errors() {
		fmt.Println(err)
	}
	require.Len(t, program.Statements, 1)
	require.IsType(t, &ast.ExpressionStatement{}, program.Statements[0])
	statement := program.Statements[0].(*ast.ExpressionStatement)
	require.Equal(t, statement.Token.Type, token.String)
	require.Equal(t, statement.Token.Literal, "this is a string")
	require.IsType(t, &ast.StringLiteral{}, statement.Expression)
	literal := statement.Expression.(*ast.StringLiteral)
	require.Equal(t, literal.Value, "this is a string")
	require.Equal(t, literal.TokenLiteral(), "this is a string")
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
		{"true == true", "(true == true)"},
		{"true != false", "(true != false)"},
		{"5 > 3 == true;3 < 5 != false", "((5 > 3) == true)((3 < 5) != false)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
		{"-(5 + 5)", "(-(5 + 5))"},
		{"!(true == true)", "(!(true == true))"},
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		require.Len(t, p.Errors(), 0)
		require.Equal(t, tt.expected, program.String())
	}
}

func TestParser_IfStatement(t *testing.T) {
	tests := []struct {
		input string
	}{
		{`if (1 > 2) { x }`},
		{`if (1 > 2) { x } else { y }`},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkErrors(t, p.Errors())
		require.IsType(t, &ast.ExpressionStatement{}, program.Statements[0])
		expression := program.Statements[0].(*ast.ExpressionStatement)
		require.Equal(t, expression.TokenLiteral(), "if")
	}
}

func TestParser_FunctionLiteral(t *testing.T) {
	tests := []struct {
		input               string
		expectedIdentifiers []string
		expectedBody        string
	}{
		{
			`fn(a, b, c) { a + b + c; }`,
			[]string{"a", "b", "c"},
			"((a + b) + c)",
		},
		{
			`fn() { 1 + 2 * 3; }`,
			[]string{},
			"(1 + (2 * 3))",
		},
		{
			`fn() {};`,
			[]string{},
			"",
		},
		{
			`
fn() {
    let x = 10;
    return x;
};`,
			[]string{},
			"let x = 10;return x;",
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkErrors(t, p.Errors())
		require.IsType(t, &ast.ExpressionStatement{}, program.Statements[0])
		expression := program.Statements[0].(*ast.ExpressionStatement)
		require.IsType(t, &ast.FunctionLiteralExpression{}, expression.Expression)
		functionLit := expression.Expression.(*ast.FunctionLiteralExpression)
		require.Equal(t, "fn", functionLit.TokenLiteral())

		identifiers := make([]string, 0, len(functionLit.Parameters))
		for _, ident := range functionLit.Parameters {
			identifiers = append(identifiers, ident.Value)
		}
		require.Equal(t, tt.expectedIdentifiers, identifiers)
		require.Equal(t, tt.expectedBody, functionLit.Body.String())
	}
}

func TestParser_CallFunction(t *testing.T) {
	tests := []struct {
		input               string
		expectedFunction    string
		expectedExpressions []string
	}{
		{
			`fn(a, b, c) { a + b + c; }(1, 2, 3)`,
			`fn(a, b, c) { ((a + b) + c) }`,
			[]string{"1", "2", "3"},
		},
		{
			`add(1, 2, 3)`,
			"add",
			[]string{"1", "2", "3"},
		},
		{
			`sub()`,
			"sub",
			[]string{},
		},
		{
			`fn() {}()`,
			`fn() {  }`,
			[]string{},
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkErrors(t, p.Errors())
		require.IsType(t, &ast.ExpressionStatement{}, program.Statements[0])
		expression := program.Statements[0].(*ast.ExpressionStatement)
		require.IsType(t, &ast.CallExpression{}, expression.Expression)
		callExpression := expression.Expression.(*ast.CallExpression)

		args := make([]string, 0, len(callExpression.Arguments))
		for _, arg := range callExpression.Arguments {
			args = append(args, arg.String())
		}
		require.Equal(t, tt.expectedExpressions, args)
		require.Equal(t, tt.expectedFunction, callExpression.Function.String())
	}
}
func checkErrors(t *testing.T, errors []string) {
	for _, err := range errors {
		fmt.Println(fmt.Errorf(err))
	}
	require.Len(t, errors, 0)
}
