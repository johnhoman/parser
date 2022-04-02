package object

import (
	"fmt"
	"strconv"
)

type Type string

const (
	TypeInteger Type = "INTEGER"
	TypeBoolean Type = "BOOLEAN"
	TypeNull    Type = "NULL"
)

func BinaryOp(operator string) func(Object, Object) Object {
    switch operator {
    case "+":
        return add
    case "-":
        return sub
    case "*":
        return mul
    case "/":
        return div
    }
    return nil
}

func add(obj1 Object, obj2 Object) Object {
    return obj1.(Numeric).Add(obj2.(Object))
}

func sub(obj1 Object, obj2 Object) Object {
    return obj1.(Numeric).Subtract(obj2.(Object))
}

func mul(obj1 Object, obj2 Object) Object {
    return obj1.(Numeric).Multiply(obj2.(Object))
}

func div(obj1 Object, obj2 Object) Object {
    return obj1.(Numeric).Divide(obj2.(Object))
}

type Object interface {
	Type() Type
	Inspect() string
}

type Numeric interface {
	Add(Object) Object
    Subtract(Object) Object
    Multiply(Object) Object
    Divide(Object) Object
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() Type      { return TypeInteger }
func (i *Integer) Add(other Object) Object {
	if o, ok := other.(*Integer); !ok {
		return nil
	} else {
		return &Integer{Value: i.Value + o.Value}
	}
}

func (i *Integer) Subtract(other Object) Object {
    if o, ok := other.(*Integer); !ok {
        return nil
    } else {
        return &Integer{Value: i.Value - o.Value}
    }
}

func (i *Integer) Multiply(other Object) Object {
    if o, ok := other.(*Integer); !ok {
        return nil
    } else {
        return &Integer{Value: i.Value * o.Value}
    }
}

func (i *Integer) Divide(other Object) Object {
    if o, ok := other.(*Integer); !ok {
        return nil
    } else {
        return &Integer{Value: i.Value / o.Value}
    }
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string { return strconv.FormatBool(b.Value) }
func (b *Boolean) Type() Type      { return TypeBoolean }

type Null struct{}

func (n *Null) Inspect() string { return "null" }
func (n *Null) Type() Type      { return TypeNull }
