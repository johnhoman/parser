package object

type Type string

const (
	TypeInteger Type = "int"
	TypeBoolean Type = "bool"
	TypeNull    Type = "NULL"
	TypeError   Type = "ERROR"
	TypeReturn  Type = "RETURN_VALUE"
)

type Object interface {
	Type() Type
	Inspect() string
}

type Null struct{}

func (n *Null) Inspect() string { return "null" }
func (n *Null) Type() Type      { return TypeNull }
