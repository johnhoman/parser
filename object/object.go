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

type Object interface {
	Type() Type
	Inspect() string
}

type addend interface{ Add(Object) Object }
type term interface{ Sub(Object) Object }
type multiplier interface{ Mul(Object) Object }
type dividend interface{ Div(Object) Object }

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

func (i *Integer) Sub(other Object) Object {
	if o, ok := other.(*Integer); !ok {
		return nil
	} else {
		return &Integer{Value: i.Value - o.Value}
	}
}

func (i *Integer) Mul(other Object) Object {
	if o, ok := other.(*Integer); !ok {
		return nil
	} else {
		return &Integer{Value: i.Value * o.Value}
	}
}

func (i *Integer) Div(other Object) Object {
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

func Add(obj1, obj2 Object) Object {
	if _, ok := obj1.(addend); !ok {
		return nil
	}
	if ans := obj1.(addend).Add(obj2); ans != nil {
		return ans
	}
	return nil
}

func Sub(obj1, obj2 Object) Object {
	term, ok := obj1.(term)
	if !ok {
		return nil
	}
	if diff := term.Sub(obj2); diff != nil {
		return diff
	}
	return nil
}

func Mul(obj1, obj2 Object) Object {
	multiplier, ok := obj1.(multiplier)
	if !ok {
		return nil
	}
	if product := multiplier.Mul(obj2); product != nil {
		return product
	}
	return nil
}

func Div(obj1, obj2 Object) Object {
	dividend, ok := obj1.(dividend)
	if !ok {
		return nil
	}
	if product := dividend.Div(obj2); product != nil {
		return product
	}
	return nil
}
