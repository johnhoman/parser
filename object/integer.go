package object

import "fmt"

type Integer struct{ Value int64 }

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() Type      { return TypeInteger }

func (i *Integer) add(o *Integer) *Integer {
	return &Integer{Value: i.Value + o.Value}
}

func (i *Integer) Add(other Object) Object {
	o, ok := other.(*Integer)
	if !ok {
		return nil
	}
	return i.add(o)
}

func (i *Integer) sub(other *Integer) *Integer {
	return &Integer{Value: i.Value - other.Value}
}

func (i *Integer) Sub(other Object) Object {
	o, ok := other.(*Integer)
	if !ok {
		return nil
	}
	return i.sub(o)
}

func (i *Integer) mul(other *Integer) *Integer {
	return &Integer{Value: i.Value * other.Value}
}

func (i *Integer) Mul(other Object) Object {
	o, ok := other.(*Integer)
	if !ok {
		return nil
	}
	return i.mul(o)
}

func (i *Integer) div(other *Integer) *Integer {
	return &Integer{Value: i.Value / other.Value}
}

func (i *Integer) Div(other Object) Object {
	o, ok := other.(*Integer)
	if !ok {
		return nil
	}
	return i.div(o)
}