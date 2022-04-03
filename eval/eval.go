package eval

import (
	"fmt"
	"mitchlang/ast"
	"mitchlang/object"
)

var (
	NullSingleton = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch n := node.(type) {
	case *ast.Program:
		return evalStatements(n.Statements)
	case *ast.BlockStatement:
		return evalBlockStatements(n.Statements)
	case *ast.IfExpression:
		condition := Eval(n.Condition)
		if condition == object.True {
			return Eval(n.Consequence)
		} else {
			if n.Alternative != nil {
				return Eval(n.Alternative)
			}
			return NullSingleton
		}
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
		return evalPrefixOperator(n.Operator, right)
	case *ast.Boolean:
		if n.Value {
			return object.True
		}
		return object.False
	case *ast.ReturnStatement:
		return &object.ReturnValue{Value: Eval(n.ReturnValue)}
	}
	return nil
}

func evalBangOperator(right object.Object) object.Object {
	switch right {
	case object.True:
		return object.False
	case object.False:
		return object.True
	default:
		return NullSingleton
	}
}

func evalMinusPrefixOperator(right object.Object) object.Object {
	switch right := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: 0 - right.Value}
	default:
		return &object.Error{
		    Message: fmt.Sprintf("unknown operator: -%s", right.Type()),
		}
	}
}

func evalPrefixOperator(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperator(right)
	case "-":
		return evalMinusPrefixOperator(right)
	default:
		return nil
	}
}

func evalInfixIntegerExpression(
	operator string,
	left object.Object,
	right object.Object,
) object.Object {
	var binaryFunc object.BinaryOpFunc
	switch operator {
	case "+":
		binaryFunc = object.Add
	case "-":
		binaryFunc = object.Sub
	case "*":
		binaryFunc = object.Mul
	case "/":
		binaryFunc = object.Div
	case "==":
		binaryFunc = object.Eq
	case "!=":
		binaryFunc = object.NotEq
	case "<":
		binaryFunc = object.Lt
	case ">":
		binaryFunc = object.Gt
	default:
		return NullSingleton
	}
	if binaryFunc != nil {
		val := binaryFunc(left, right)
		if e, ok := val.(*object.Error); ok {
			return e
		}
		if val != nil {
			return val
		}
	}
	return nil
}

func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object
	for _, statement := range statements {
		result = Eval(statement)
		if _, ok := result.(*object.Error); ok {
			return result
		}
		if rv, ok := result.(*object.ReturnValue); ok {
			return rv.Value
		}
	}
	return result
}

func evalBlockStatements(statements []ast.Statement) object.Object {
	var result object.Object
	for _, statement := range statements {
		result = Eval(statement)
		switch result.(type) {
		case *object.Error:
			return result
		case *object.ReturnValue:
			return result
		}
	}
	return result
}
