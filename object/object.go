package object

type Type string

const (
	TypeString   Type = "str"
	TypeInteger  Type = "int"
	TypeBoolean  Type = "bool"
	TypeNull     Type = "NULL"
	TypeError    Type = "ERROR"
	TypeReturn   Type = "RETURN_VALUE"
	TypeFunction Type = "FUNCTION"
	TypeBuiltin  Type = "BUILTIN"
	TypeList     Type = "List"
)

func (t Type) String() string { return string(t) }

type Object interface {
	Type() Type
	Inspect() string
}

type Null struct{}

func (n *Null) Inspect() string { return "null" }
func (n *Null) Type() Type      { return TypeNull }
