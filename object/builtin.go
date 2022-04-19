package object

import (
	"fmt"
	"io"
	"os"
)

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() Type      { return TypeBuiltin }
func (b *Builtin) Inspect() string { return "BUILTIN_FUNCTION" }

var _ Object = &Builtin{}

type iterable interface{ Len() Object }

func BuiltinLen(args ...Object) Object {
	if len(args) > 1 {
		return NewTypeError("expected 1 position argument but received %d", len(args))
	}
	obj := args[0]
	it, ok := obj.(iterable)
	if !ok {
		return NewTypeError("object is not iterable: %s", obj.Type())
	}
	return it.Len()
}

func BuiltinAdd(args ...Object) Object {
	if len(args) != 2 {
		return NewTypeError("expected 2 position arguments but received %d", len(args))
	}
	one, two := args[0], args[1]
	rv := Add(one, two)
	if rv == nil {
		return nil
	}
	return rv
}

func BuiltinExit(args ...Object) Object {
	if len(args) > 1 {
		return NewTypeError("expected 1 positional argument but received %d", len(args))
	}
	code := 0
	if len(args) == 1 {
		one := args[0]
		i, ok := one.(*Integer)
		if !ok {
			return NewTypeError(
				"expected positional argument 1 to be type %s but received type %s",
				TypeInteger,
				one.Type(),
			)
		}
		code = int(i.Value)
	}
	os.Exit(code)
	return NullValue
}

func BuiltinList(args ...Object) Object {
	if len(args) > 1 {
		return NewTypeError("expected 1 positional argument but received %d", len(args))
	}
	if _, ok := args[0].(interface{ List() Object }); !ok {
		return NewTypeError("object %s is not iterable", args[0].Type())
	}
	it := args[0].(interface{ List() Object })
	return it.List()
}

func builtinPrint(args ...Object) Object {
	if len(args) > 1 {
		return NewTypeError("expected 1 positional arguments but received %d", len(args))
	}
	if _, ok := args[0].(fmt.Stringer); ok {
		stringer := args[0].(fmt.Stringer)
		_, _ = os.Stdout.WriteString(stringer.String())
	} else {
		s := args[0].Inspect()
		_, _ = os.Stdout.WriteString(s)
	}
	return NullValue
}

func BuiltinPrintln(args ...Object) Object {
	obj := builtinPrint(args...)
	if obj != NullValue {
		return obj
	}
	_, _ = io.WriteString(os.Stdout, "\n")
	return NullValue
}