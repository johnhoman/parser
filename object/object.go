package object

type Type string

const (
	TypeInteger Type = "INTEGER"
	TypeBoolean Type = "BOOLEAN"
	TypeNull    Type = "NULL"
	TypeError   Type = "ERROR"
)

type Object interface {
	Type() Type
	Inspect() string
}

type Null struct{}

func (n *Null) Inspect() string { return "null" }
func (n *Null) Type() Type      { return TypeNull }
