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
