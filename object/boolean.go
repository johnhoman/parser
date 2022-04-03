package object

import "strconv"

var (
	True  = &Boolean{Value: true}
	False = &Boolean{Value: false}
)

type Boolean struct{ Value bool }

func (b *Boolean) Inspect() string { return strconv.FormatBool(b.Value) }
func (b *Boolean) Type() Type      { return TypeBoolean }
func (b *Boolean) eq(o *Boolean) *Boolean {
	if b.Value == o.Value {
		return True
	}
	return False
}
func (b *Boolean) Eq(other Object) Object {
	o, ok := other.(*Boolean)
	if !ok {
		return nil
	}
	return b.eq(o)
}

func (b *Boolean) Lt(other Object) Object {
	// This is wrong -- should raise an error
	return NewTypeError("boolean does not support < operator")
}
