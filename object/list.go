package object

import (
	"bytes"
	"strings"
)

type List struct {
	Values []Object
}

func (l *List) Type() Type {
	return TypeList
}

func (l *List) Inspect() string {
	values := make([]string, 0, len(l.Values))

	for k := range l.Values {
		values = append(values, l.Values[k].Inspect())
	}
	out := new(bytes.Buffer)
	out.WriteString("[")
	out.WriteString(strings.Join(values, ", "))
	out.WriteString("]")
	return out.String()
}

var _ Object = &List{}
