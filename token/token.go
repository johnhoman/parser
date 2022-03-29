package token

type Type string

func (t Type) String() string { return string(t) }

type Token struct {
	Type Type
	Literal string
}

const (
	Illegal   Type = "ILLEGAL"
	EOF       Type = "EOF"
	Ident     Type = "IDENT"
	Int       Type = "INT"
	Assign    Type = "ASSIGN"
	Plus      Type = "PLUS"

	Comma     Type = "COMMA"
	SemiColon Type = "SEMI_COLON"

	LParen    Type = "L_PAREN"
	RParen    Type = "R_PAREN"
	LBrace    Type = "L_BRACE"
	RBrace    Type = "R_BRACE"

	Function  Type = "FUNCTION"
	Let       Type = "LET"
)
