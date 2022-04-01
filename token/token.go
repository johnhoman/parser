package token

type Type string

func (t Type) String() string { return string(t) }
func (t Type) Repr() string   { return t.String() }

const (
	Illegal  Type = "ILLEGAL"
	EOF      Type = "EOF"
	Ident    Type = "IDENTIFIER"
	Int      Type = "INTEGER"
	Assign   Type = "="
	Plus     Type = "+"
	Minus    Type = "-"
	Bang     Type = "!"
	Asterisk Type = "*"
	Slash    Type = "/"
	LT       Type = "<"
	GT       Type = ">"

	Eq    Type = "="
	NotEq Type = "!="

	Comma     Type = ","
	SemiColon Type = ";"

	LParen   Type = "("
	RParen   Type = ")"
	LBrace   Type = "{"
	RBrace   Type = "}"

	Function Type = "fn"
	Let      Type = "let"
	True     Type = "true"
	False    Type = "false"
	If       Type = "if"
	Else     Type = "else"
	Return   Type = "return"
)

type Token struct {
	Type    Type
	Literal string
}

func (t *Token) IsType(tt Type) bool { return t.Type == tt }

func (t *Token) String() string {
	return t.Literal
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
	"fn":     Function,
	"let":    Let,
	"true":   True,
	"false":  False,
	"if":     If,
	"else":   Else,
	"return": Return,
}

func lookupIdent(ident string) Type {
	if tokenType, ok := keywords[ident]; ok {
		// found keyword
		return tokenType
	}
	return Ident
}
