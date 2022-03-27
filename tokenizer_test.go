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
    expectedAst := parser.StatementList{
        Statements: []parser.Statement{{
            ExpressionStatement: &parser.ExpressionStatement{
                Expression: parser.Expression{
                    Literal: parser.Literal{Value: 42},
                },
            }},
        },
    }
	require.Equal(t, ast.Body, expectedAst)
}

func TestParser_IntMultiple(t *testing.T) {
	expected1 := 42
	expected2 := 43
	ast, err := parser.New().Parse(fmt.Sprintf("%d;\n%d;", expected1, expected2))
    require.Nil(t, err)
    expectedAst := parser.StatementList{
        Statements: []parser.Statement{
            {
                ExpressionStatement: &parser.ExpressionStatement{
                    Expression: parser.Expression{
                        Literal: parser.Literal{Value: expected1},
                    },
                },
            },
            {
                ExpressionStatement: &parser.ExpressionStatement{
                    Expression: parser.Expression{
                        Literal: parser.Literal{Value: expected2},
                    },
                },
            },
        },
    }
    require.Equal(t, ast.Body, expectedAst)
}

func TestParser_StringLiteral(t *testing.T) {
	expected := "string literal"
	ast, err := parser.New().Parse(fmt.Sprintf(`"%s";`, expected))
	require.Nil(t, err)
    expectedAst := parser.StatementList{
        Statements: []parser.Statement{{
            ExpressionStatement: &parser.ExpressionStatement{
                Expression: parser.Expression{
                    Literal: parser.Literal{Value: "string literal"},
                },
            }},
        },
    }
    require.Equal(t, ast.Body, expectedAst)
}

func TestParser_Whitespace(t *testing.T) {
	expected := 42
	ast, err := parser.New().Parse(fmt.Sprintf(`     %d;`, expected))
	require.Nil(t, err)
    expectedAst := parser.StatementList{
        Statements: []parser.Statement{{
            ExpressionStatement: &parser.ExpressionStatement{
                Expression: parser.Expression{
                    Literal: parser.Literal{Value: 42},
                },
            }},
        },
    }
    require.Equal(t, ast.Body, expectedAst)
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
    expectedAst := parser.StatementList{
        Statements: []parser.Statement{{
            ExpressionStatement: &parser.ExpressionStatement{
                Expression: parser.Expression{
                    Literal: parser.Literal{Value: 42},
                },
            }},
        },
    }
    require.Equal(t, ast.Body, expectedAst)
}

func TestParser_BlockStatement(t *testing.T) {
    expected := 42
    doc := fmt.Sprintf(`{
  %d;
}`, expected)
    ast, err := parser.New().Parse(doc)
    require.Nil(t, err)
    expectedAst := parser.StatementList{
        Statements: []parser.Statement{{
            BlockStatement: &parser.BlockStatement{
                StatementList: parser.StatementList{Statements: []parser.Statement{{
                    ExpressionStatement: &parser.ExpressionStatement{
                        Expression: parser.Expression{
                            Literal: parser.Literal{Value: expected},
                        },
                    },
                }}},
            },
        }},
    }
    require.Equal(t, ast.Body, expectedAst)
}

func TestParser_BlockStatementNested(t *testing.T) {
    expected := 42
    doc := fmt.Sprintf(`{
  {42;}
  {}
}`)
    ast, err := parser.New().Parse(doc)
    require.Nil(t, err)
    expectedAst := parser.StatementList{
        Statements: []parser.Statement{{
            BlockStatement: &parser.BlockStatement{
                StatementList: parser.StatementList{Statements: []parser.Statement{{
                    BlockStatement: &parser.BlockStatement{
                        StatementList: parser.StatementList{Statements: []parser.Statement{{
                            ExpressionStatement: &parser.ExpressionStatement{
                                Expression: parser.Expression{
                                    Literal: parser.Literal{Value: expected},
                                },
                            },
                        }}},
                    },
                },{
                    BlockStatement: &parser.BlockStatement{
                        StatementList: parser.StatementList{},
                    },
                }}},
            },
        }},
    }
    require.Equal(t, ast.Body, expectedAst)
}
