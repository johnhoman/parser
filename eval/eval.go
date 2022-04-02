package eval

import (
	"mitchlang/ast"
	"mitchlang/object"
)

var (
	TrueSingleton  = &object.Boolean{Value: true}
	FalseSingleton = &object.Boolean{Value: false}
	NullSingleton  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch n := node.(type) {
	case *ast.Program:
		var result object.Object
		for _, statement := range n.Statements {
			result = Eval(statement)
		}
		return result
	case *ast.InfixExpression:
		left := Eval(n.Left)
		right := Eval(n.Right)
		return evalInfixIntegerExpression(n.Operator, left, right)
	case *ast.ExpressionStatement:
		return Eval(n.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}
	case *ast.PrefixExpression:
		right := Eval(n.Right)
		return parsePrefixOperator(n.Operator, right)
	case *ast.Boolean:
		if n.Value {
			return TrueSingleton
		}
		return FalseSingleton
	}
	return nil
}

func parseBangOperator(right object.Object) object.Object {
	switch right {
	case TrueSingleton:
		return FalseSingleton
	case FalseSingleton:
		return TrueSingleton
	default:
		return nil
	}
}

func parseMinusPrefixOperator(right object.Object) object.Object {
	switch right := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: 0 - right.Value}
	}
	return nil
}

func parsePrefixOperator(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return parseBangOperator(right)
	case "-":
		return parseMinusPrefixOperator(right)
	default:
		return nil
	}
}

func evalInfixIntegerExpression(
	operator string,
	left object.Object,
	right object.Object,
) object.Object {
	switch operator {
	case "+":
		if sum := object.Add(left, right); sum != nil {
			return sum
		}
		// TODO: error
		return nil
	case "-":
		if sum := object.Sub(left, right); sum != nil {
			return sum
		}
		// TODO: error
		return nil
	case "*":
		if product := object.Mul(left, right); product != nil {
			return product
		}
		// TODO: error
		return nil
	case "/":
		if quo := object.Div(left, right); quo != nil {
			return quo
		}
		// TODO: error
		return nil
	}
	return nil
}
