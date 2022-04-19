package eval

import (
	"fmt"
	"mitchlang/ast"
	"mitchlang/object"
	"strings"
)

var (
	NullSingleton = &object.Null{}
)

var builtins = map[string]*object.Builtin{
	"len":  {Fn: object.BuiltinLen},
	"add":  {Fn: object.BuiltinAdd},
	"exit": {Fn: object.BuiltinExit},
	"list": {Fn: object.BuiltinList},
	"print": {Fn: object.BuiltinPrintln},
}

func Eval(node ast.Node, env *object.Env) object.Object {
	switch n := node.(type) {
	case *ast.Program:
		return evalStatements(n.Statements, env)
	case *ast.BlockStatement:
		return evalBlockStatements(n.Statements, env)
	case *ast.IfExpression:
		condition := Eval(n.Condition, env)
		if isError(condition) {
			return condition
		}
		if condition == object.True {
			return Eval(n.Consequence, env)
		} else {
			if n.Alternative != nil {
				return Eval(n.Alternative, env)
			}
			return NullSingleton
		}
	case *ast.InfixExpression:
		left := Eval(n.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixIntegerExpression(n.Operator, left, right)
	case *ast.ExpressionStatement:
		return Eval(n.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}
	case *ast.StringLiteral:
		return &object.String{Value: n.Value}
	case *ast.PrefixExpression:
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixOperator(n.Operator, right)
	case *ast.Boolean:
		if n.Value {
			return object.True
		}
		return object.False
	case *ast.ReturnStatement:
		return &object.ReturnValue{Value: Eval(n.ReturnValue, env)}
	case *ast.LetStatement:
		obj := Eval(n.Value, env)
		if isError(obj) {
			return obj
		}
		env.Set(n.Name.Value, obj)
		return obj
	case *ast.Identifier:
		if obj, ok := env.Get(n.Value); ok {
			return obj
		}
		if obj, ok := builtins[n.Value]; ok {
			return obj
		}

		return &object.Error{
			Message: fmt.Sprintf("identifier not found: %s", n.Value),
		}
	case *ast.FunctionLiteralExpression:
		obj := &object.Function{
			Parameters: n.Parameters,
			Body:       n.Body,
			Env:        env,
		}
		return obj
	case *ast.CallExpression:
		obj := Eval(n.Function, env)
		if isError(obj) {
			return obj
		}
		args := make([]object.Object, 0, len(n.Arguments))
		for _, exp := range n.Arguments {
			out := Eval(exp, env)
			if isError(out) {
				return out
			}
			args = append(args, out)
		}
		switch fn := obj.(type) {
		case *object.Builtin:
			return fn.Fn(args...)
		case *object.Function:
			functionEnv := fn.Env.Push()
			for k := range fn.Parameters {
				functionEnv.Set(fn.Parameters[k].Value, args[k])
			}
			// Need to remove the variables from the environment
			out := Eval(fn.Body, functionEnv)
			if rv, ok := out.(*object.ReturnValue); ok { return rv.Value }
			return out
		default:
			return &object.Error{Message: fmt.Sprintf("not a function %s", fn.Type())}
		}
	case *ast.ListExpression:
		items := make([]object.Object, 0, len(n.Items))
		for k := range n.Items {
			item := Eval(n.Items[k], env)
			if isError(item) {
				return item
			}
			items = append(items, item)
		}
		return &object.List{Values: items}
	case *ast.IndexExpression:
		items := Eval(n.Left, env)
		if isError(items) {
			return items
		}
		rank := Eval(n.Index, env)
		if isError(rank) {
			return rank
		}
		integer, ok := rank.(*object.Integer)
		if !ok {
			return object.NewTypeError("expected integer, got %s", rank.Type())
		}
		index := int(integer.Value)
		length := int(object.BuiltinLen(items).(*object.Integer).Value)
		if index < 0 {
			index = length + index
		}

		if index >= length || index < 0 {
			typeString := strings.ToLower(items.Type().String())
			return &object.Error{
				ErrorType: "IndexError",
				Message: fmt.Sprintf("%s index out of range", typeString),
			}
		}
		switch ob := items.(type) {
		case *object.String:
			return &object.String{Value: string(ob.Value[index])}
		case *object.List:
			return ob.Values[index]
		default:
			return object.NewTypeError("expected list or string, got %s", ob.Type())
		}

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

func evalStatements(statements []ast.Statement, env *object.Env) object.Object {
	var result object.Object
	for _, statement := range statements {
		result = Eval(statement, env)
		if _, ok := result.(*object.Error); ok {
			return result
		}
		if rv, ok := result.(*object.ReturnValue); ok {
			return rv.Value
		}
	}
	return result
}

func evalBlockStatements(statements []ast.Statement, env *object.Env) object.Object {
	var result object.Object
	for _, statement := range statements {
		result = Eval(statement, env)
		switch result.(type) {
		case *object.Error:
			return result
		case *object.ReturnValue:
			return result
		}
	}
	// explicitly return value or null
	return NullSingleton
}

func isError(obj object.Object) bool {
	if _, ok := obj.(*object.Error); ok {
		return true
	}
	return false
}
