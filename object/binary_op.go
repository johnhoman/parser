package object

import (
	"reflect"
)

type addend interface{ Add(Object) Object }
type term interface{ Sub(Object) Object }
type multiplier interface{ Mul(Object) Object }
type dividend interface{ Div(Object) Object }

type comparable interface {
	Eq(Object) Object
	Lt(Object) Object
}

type BinaryOpFunc func(ob1, ob2 Object) Object


func strict(opFunc BinaryOpFunc, v interface{}, op string) BinaryOpFunc {
	return func(ob1, ob2 Object) Object {

		t1 := ob1.Type()
		t2 := ob2.Type()

		inf := reflect.TypeOf(v).Elem()
		if !reflect.TypeOf(ob1).Implements(inf) {
			return NewTypeError("invalid operation: %s %s %s", t1, op, t2)
		}
		if t1 != t2 {
			return NewTypeError("type mismatch: %s %s %s", t1, op, t2)
		}
		return opFunc(ob1, ob2)
	}
}

func add(obj1, obj2 Object) Object {
	if _, ok := obj1.(addend); !ok {
		// type error - obj1 does not support add
		return nil
	}
	if ans := obj1.(addend).Add(obj2); ans != nil {
		return ans
	}
	return nil
}

// Add - adds two objects and returns the result. Only
// objects of the same type can be added.
var Add = strict(add, (*addend)(nil), "+")

func sub(obj1, obj2 Object) Object {
	term := obj1.(term)
	if diff := term.Sub(obj2); diff != nil {
		return diff
	}
	return nil
}

var Sub = strict(sub, (*term)(nil), "-")

func mul(obj1, obj2 Object) Object {
	multiplier, ok := obj1.(multiplier)
	if !ok {
		return nil
	}
	if product := multiplier.Mul(obj2); product != nil {
		return product
	}
	return nil
}

var Mul = strict(mul, (*multiplier)(nil), "*")

func div(obj1, obj2 Object) Object {
	dividend, ok := obj1.(dividend)
	if !ok {
		return nil
	}
	if product := dividend.Div(obj2); product != nil {
		return product
	}
	return nil
}

var Div = strict(div, (*dividend)(nil), "/")

func eq(obj1, obj2 Object) Object {
	ob, ok := obj1.(comparable)
	if !ok {
		return nil
	}
	if result := ob.Eq(obj2); result != nil {
		return result
	}
	return nil
}

var Eq = strict(eq, (*comparable)(nil), "==")

func notEq(obj1, obj2 Object) Object {
	eq := Eq(obj1, obj2)
	if eq == nil {
		return nil
	}
	if eq == True {
		return False
	}
	return True
}

var NotEq = strict(notEq, (*comparable)(nil), "!=")

func lt(obj1, obj2 Object) Object {
	ob, ok := obj1.(comparable)
	if !ok {
		return nil
	}
	res := ob.Lt(obj2)
	if res != nil {
		return res
	}
	return nil
}

var Lt = strict(lt, (*comparable)(nil), "<")

func gt(obj1, obj2 Object) Object {
	lt := Lt(obj1, obj2)
	eq := Eq(obj1, obj2)
	if lt == nil || eq == nil {
		return nil
	}
	if lt == True || eq == True {
		return False
	}
	return True
}

var Gt = strict(gt, (*comparable)(nil), ">")
