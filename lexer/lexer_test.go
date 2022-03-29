package lexer_test

import (
	"github.com/stretchr/testify/require"
	"mitchlang/lexer"
	"mitchlang/token"
	"testing"
)

func TestLexer_NextToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
    x + y;
};

let result = add(five, ten);
!-/*5;
5 < 10 > 5

if (5 < 10) {
    return true;
} else {
    return false;
}

10 == 10
10 != 9
`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.Let, "let"},
		{token.Ident, "five"},
		{token.Assign, "="},
		{token.Int, "5"},
		{token.SemiColon, ";"},
		{token.Let, "let"},
		{token.Ident, "ten"},
		{token.Assign, "="},
		{token.Int, "10"},
		{token.SemiColon, ";"},
		{token.Let, "let"},
		{token.Ident, "add"},
		{token.Assign, "="},
		{token.Function, "fn"},
		{token.LParen, "("},
		{token.Ident, "x"},
		{token.Comma, ","},
		{token.Ident, "y"},
		{token.RParen, ")"},
		{token.LBrace, "{"},
		{token.Ident, "x"},
		{token.Plus, "+"},
		{token.Ident, "y"},
		{token.SemiColon, ";"},
		{token.RBrace, "}"},
		{token.SemiColon, ";"},
		{token.Let, "let"},
		{token.Ident, "result"},
		{token.Assign, "="},
		{token.Ident, "add"},
		{token.LParen, "("},
		{token.Ident, "five"},
		{token.Comma, ","},
		{token.Ident, "ten"},
		{token.RParen, ")"},
		{token.SemiColon, ";"},
		// !-/*5;
		{token.Bang, "!"},
		{token.Minus, "-"},
		{token.Slash, "/"},
		{token.Asterisk, "*"},
		{token.Int, "5"},
		{token.SemiColon, ";"},
		// 5 < 10 > 5
		{token.Int, "5"},
		{token.LT, "<"},
		{token.Int, "10"},
		{token.GT, ">"},
		{token.Int, "5"},
		// if (5 < 10) {
		// 	return true;
		// } else {
		// 	return false;
		// }
		{token.If, "if"},
		{token.LParen, "("},
		{token.Int, "5"},
		{token.LT, "<"},
		{token.Int, "10"},
		{token.RParen, ")"},
		{token.LBrace, "{"},
		{token.Return, "return"},
		{token.True, "true"},
		{token.SemiColon, ";"},
		{token.RBrace, "}"},
		{token.Else, "else"},
		{token.LBrace, "{"},
		{token.Return, "return"},
		{token.False, "false"},
		{token.SemiColon, ";"},
		{token.RBrace, "}"},
		// 10 == 10
		{token.Int, "10"},
		{token.Eq, "=="},
		{token.Int, "10"},
		// 10 != 9
		{token.Int, "10"},
		{token.NotEq, "!="},
		{token.Int, "9"},
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
