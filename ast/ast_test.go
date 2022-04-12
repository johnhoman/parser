package ast

import (
	"github.com/stretchr/testify/require"
	"mitchlang/token"
	"testing"
)

func TestProgram_String(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.NewFromString(token.Let, "let"),
				Name: &Identifier{
					Token: token.NewFromString(token.Ident, "myVar"),
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.NewFromString(token.Ident, "anotherVar"),
					Value: "anotherVar",
				},
			},
		},
	}
	require.Equal(t, "let myVar = anotherVar;", program.String())
}

func TestProgram_List(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&ExpressionStatement{
				Expression: &ListExpression{
					Token: token.New(token.LBracket, '['),
					Items: []Expression{
						&IntegerLiteral{
							Token: token.New(token.Int, '1'),
							Value: 1,
						},
						&IntegerLiteral{
							Token: token.New(token.Int, '2'),
							Value: 2,
						},
						&StringLiteral{
							Token: token.NewFromString(token.String, "10"),
							Value: "10",
						},
						&Boolean{
							Token: token.NewIdentifier("true"),
							Value: true,
						},
					},
				},
			},
		},
	}
	require.Equal(t, `[1, 2, "10", true]`, program.String())
}


func TestProgram_IndexExpression(t *testing.T) {
	program := Program{
		Statements: []Statement{
			&ExpressionStatement{
				Expression: &IndexExpression{
					Left: &ListExpression{
						Token: token.New(token.LBracket, '['),
						Items: []Expression{
							&IntegerLiteral{
								Token: token.New(token.Int, '1'),
								Value: 1,
							},
							&IntegerLiteral{
								Token: token.New(token.Int, '2'),
								Value: 2,
							},
							&StringLiteral{
								Token: token.NewFromString(token.String, "10"),
								Value: "10",
							},
							&Boolean{
								Token: token.NewIdentifier("true"),
								Value: true,
							},
						},
					},
					Index: &IntegerLiteral{
						Token: token.NewFromString(token.Int, "22"),
						Value: 22,
					},
				},
			},
		},
	}
	require.Equal(t, `([1, 2, "10", true][22])`, program.String())
}