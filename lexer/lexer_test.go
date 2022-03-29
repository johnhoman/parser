package lexer_test

import (
	"github.com/stretchr/testify/require"
	"mitchlang/lexer"
	"mitchlang/token"
	"testing"
)

func TestLexer_NextToken(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.Assign, "="},
		{token.Plus, "+"},
		{token.LParen, "("},
		{token.RParen, ")"},
		{token.LBrace, "{"},
		{token.RBrace, "}"},
		{token.Comma, ","},
		{token.SemiColon, ";"},
	}
	lex := lexer.New(input)

	for _, test := range tests {
		tok := lex.NextToken()
		t.Run(test.expectedType.String(), func(t *testing.T) {
			require.Equal(t, test.expectedType, tok.Type)
			require.Equal(t, test.expectedLiteral, tok.Literal)
		})
	}
}
