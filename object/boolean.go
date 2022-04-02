package object

import "strconv"

type Boolean struct{ Value bool }

func (b *Boolean) Inspect() string { return strconv.FormatBool(b.Value) }
func (b *Boolean) Type() Type      { return TypeBoolean }
