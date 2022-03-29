package token

import "fmt"

type Type string

func (t Type) String()       string { return string(t) }
func (t Type) Repr()         string { return t.String() }

const (
	Illegal  Type = "ILLEGAL"
	EOF      Type = "EOF"
	Ident    Type = "IDENTIFIER"
	Int      Type = "INTEGER"
	Assign   Type = "ASSIGN"
	Plus     Type = "PLUS"
	Minus    Type = "MINUS"
	Bang     Type = "BANG"
	Asterisk Type = "ASTERISK"
	Slash    Type = "SLASH"
	LT       Type = "LESS_THAN"
	GT       Type = "GREATER_THAN"

	Eq       Type = "EQUAL"
	NotEq    Type = "NOT_EQUAL"

	Comma     Type = "COMMA"
	SemiColon Type = "SEMI_COLON"

	LParen Type = "LEFT_PARENTHESIS"
	RParen Type = "RIGHT_PARENTHESIS"
	LBrace Type = "LEFT_BRACE"
	RBrace Type = "RIGHT_BRACE"

	Function Type = "FUNCTION"
	Let      Type = "LET"
	True     Type = "TRUE"
	False    Type = "FALSE"
	If       Type = "IF"
	Else     Type = "ELSE"
	Return   Type = "RETURN"
)

type Token struct {
	Type    Type
	Literal string
}

func (t *Token) IsType(tt Type) bool { return t.Type == tt }

func (t *Token) String() string {
	return fmt.Sprintf(`%s("%s")`, t.Type, t.Literal)
}

func NewFromString(t Type, literal string) *Token {
	return &Token{Type: t, Literal: literal}
}

func NewIdentifier(ident string) *Token {
	identifier := lookupIdent(ident)
	return NewFromString(identifier, ident)
}

func New(t Type, literal ...byte) *Token {
	return NewFromString(t, string(literal))
}

var keywords = map[string]Type{
	"fn":  Function,
	"let": Let,
	"true": True,
	"false": False,
	"if": If,
	"else": Else,
	"return": Return,
}

func lookupIdent(ident string) Type {
	if tokenType, ok := keywords[ident]; ok {
		// found keyword
		return tokenType
	}
	return Ident
}
