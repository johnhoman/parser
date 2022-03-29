package parser

import (
    "fmt"
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

    tests := []struct{
        expected string
    } {
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
`

    l := lexer.New(input)
    p := New(l)

    program := p.ParseProgram()
    require.NotNil(t, program)
    require.Len(t, program.Statements, 0)
    require.Len(t, p.Errors(), 1, fmt.Sprintf("%#v", p.Errors()))
}
