package token

import "fmt"

type Type string

func (t Type) String() string { return string(t) }

const (
	Illegal Type = "ILLEGAL"
	EOF     Type = "EOF"
	Ident   Type = "IDENT"
	Int     Type = "INT"
	Assign  Type = "ASSIGN"
	Plus    Type = "PLUS"

	Comma     Type = "COMMA"
	SemiColon Type = "SEMI_COLON"

	LParen Type = "L_PAREN"
	RParen Type = "R_PAREN"
	LBrace Type = "L_BRACE"
	RBrace Type = "R_BRACE"

	Function Type = "FUNCTION"
	Let      Type = "LET"
)

type Token struct {
	Type    Type
	Literal string
}

func (t *Token) String() string {
	if t.Literal != "" {
		return fmt.Sprintf(`%s("%s")`, t.Type, t.Literal)
	}
	return t.Type.String()
}

func New(t Type, literal string) *Token {
	return &Token{Type: t, Literal: literal}
}
