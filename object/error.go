package object

type Error struct{ Message string }

func (e *Error) Type() Type      { return TypeError }
func (e *Error) Inspect() string { return e.Message }

var _ Object = &Error{}