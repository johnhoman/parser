package parser_test

import (
    "fmt"
    "testing"

    "github.com/stretchr/testify/require"

    "parser"
)

func TestParser_IntLiteral(t *testing.T) {
    expected := 42
    ast, err := parser.New().Parse(fmt.Sprintf(`%d;`, expected))
    require.Nil(t, err)
    lit := ast.Body.Statements[0].ExpressionStatement.Expression.Literal
    require.Equal(t, lit.Type(), parser.IntLiteralType)
    require.Equal(t, lit.Value(), expected)
}

func TestParser_IntMultiple(t *testing.T) {
    expected1 := 42
    expected2 := 43
    ast, err := parser.New().Parse(fmt.Sprintf("%d;\n%d;", expected1, expected2))
    require.Nil(t, err)
    for k, exp := range []int{expected1, expected2} {
        lit := ast.Body.Statements[k].ExpressionStatement.Expression.Literal
        require.Equal(t, lit.Type(), parser.IntLiteralType)
        require.Equal(t, lit.Value(), exp)
    }
}

func TestParser_StringLiteral(t *testing.T) {
    expected := "string literal"
    ast, err := parser.New().Parse(fmt.Sprintf(`"%s";`, expected))
    require.Nil(t, err)
    lit := ast.Body.Statements[0].ExpressionStatement.Expression.Literal
    require.Equal(t, lit.Type(), parser.StringLiteralType)
    require.Equal(t, lit.Value(), expected)
}

func TestParser_Whitespace(t *testing.T) {
    expected := 42
    ast, err := parser.New().Parse(fmt.Sprintf(`     %d;`, expected))
    require.Nil(t, err)
    lit := ast.Body.Statements[0].ExpressionStatement.Expression.Literal
    require.Equal(t, lit.Type(), parser.IntLiteralType)
    require.Equal(t, lit.Value(), expected)
}

func TestParser_Comment(t *testing.T) {
    expected := 42
    doc := fmt.Sprintf(
`// this is a comment
       %d;`,
       expected,
    )
    ast, err := parser.New().Parse(doc)
    require.Nil(t, err)
    lit := ast.Body.Statements[0].ExpressionStatement.Expression.Literal
    require.Equal(t, lit.Type(), parser.IntLiteralType)
    require.Equal(t, lit.Value(), expected)
}
