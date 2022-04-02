package object

type addend interface{ Add(Object) Object }
type term interface{ Sub(Object) Object }
type multiplier interface{ Mul(Object) Object }
type dividend interface{ Div(Object) Object }

type comparable interface {
	Eq(Object) Object
	Lt(Object) Object
}

type BinaryOpFunc func(ob1, ob2 Object) Object

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

func Eq(obj1, obj2 Object) Object {
	ob, ok := obj1.(comparable)
	if !ok {
		return nil
	}
	if result := ob.Eq(obj2); result != nil {
		return result
	}
	return nil
}

func NotEq(obj1, obj2 Object) Object {
	eq := Eq(obj1, obj2)
	if eq == nil {
		return nil
	}
	if eq == True {
		return False
	}
	return True
}

func Lt(obj1, obj2 Object) Object {
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

func Gt(obj1, obj2 Object) Object {
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
