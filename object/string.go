package object

import "bytes"

type String struct {
	Value string
}

func (s *String) add(o *String) *String {
	return &String{Value: s.Value + o.Value}
}

func (s *String) Add(other Object) Object {
	s2, ok := other.(*String)
	if !ok {
		return nil
	}
	return s.add(s2)
}

func (s *String) length() *Integer {
	l := len(s.Value)
	return &Integer{Value: int64(l)}
}

func (s *String) Len() Object { return s.length() }

func (s *String) list() *List {
	values := make([]Object, 0, len(s.Value))
	for _, b := range s.Value {
		values = append(values, &String{Value: string(b)})
	}
	return &List{Values: values}
}

func (s *String) List() Object { return s.list() }

func (s *String) Type() Type  { return TypeString }

func (s *String) Inspect() string {
	out := new(bytes.Buffer)
	out.WriteByte('"')
	out.WriteString(s.Value)
	out.WriteByte('"')
	return out.String()
}

var _ Object = &String{}
var _ addend = &String{}