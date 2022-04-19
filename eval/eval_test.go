package eval

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	"mitchlang/lexer"
	"mitchlang/object"
	"mitchlang/parser"
)

func testParseInput(in string) object.Object {
	l := lexer.New(in)
	p := parser.New(l)
	env := object.NewEnv()
	for _, err := range p.Errors() {
		fmt.Println(err)
	}
	return Eval(p.ParseProgram(), env)
}

func TestEval_IntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-10", -10},
		{"5 + 5 + 5 + 5", 20},
		{"2 * 2 * 2 * 2", 16},
		{"-50 + 100 - 50 + 100", 100},
		{"5 + 2 * 10", 25},
	}

	for _, subtest := range tests {
		obj := testParseInput(subtest.input)
		require.IsType(t, &object.Integer{}, obj)
		integer := obj.(*object.Integer)
		require.Equal(t, subtest.expected, integer.Value)
	}
}

func TestEval_StringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"5"`, "5"},
		{`"foo bar"`, "foo bar"},
		{`"foo" + "bar"`, "foobar"},
	}

	for _, subtest := range tests {
		obj := testParseInput(subtest.input)
		require.IsType(t, &object.String{}, obj)
		s := obj.(*object.String)
		require.Equal(t, subtest.expected, s.Value)
	}
}

func TestEval_BooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"5 == 5", true},
		{"5 == 10", false},
		{"-5 == 12", false},
		{"5 != 5", false},
		{"5 != 10", true},
		{"-5 != 12", true},
		{"5 > 5", false},
		{"5 < 10", true},
		{"true == true", true},
		{"true == false", false},
		{"false == true", false},
		{"false == false", true},
		{"true != true", false},
		{"true != false", true},
		{"false != true", true},
		{"false != false", false},
	}

	for _, subtest := range tests {
		obj := testParseInput(subtest.input)
		require.IsType(t, &object.Boolean{}, obj)
		boolean := obj.(*object.Boolean)
		require.Equal(t, subtest.expected, boolean.Value)
	}
}

func TestEval_BangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!!true", false},
		{"!!!false", true},
	}

	for _, subtest := range tests {
		obj := testParseInput(subtest.input)
		require.IsType(t, &object.Boolean{}, obj)
		boolean := obj.(*object.Boolean)
		require.Equal(t, subtest.expected, boolean.Value)
	}
}

func TestEval_MinusPrefixOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
	}

	for _, subtest := range tests {
		evaluated := testParseInput(subtest.input)
		require.IsType(t, &object.Integer{}, evaluated)
		integer := evaluated.(*object.Integer)
		require.Equal(t, subtest.expected, integer.Value)
	}
}

func TestEval_IfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { return 10; }", 10},
		{"if (false) { return 10; } else { return 20; }", 20},
		{"if (10 == 10) { return 10; }", 10},
		{"if (10 != 10) { return 10; } else { return 20; }", 20},
		{"if (1 < 10) { return 10; }", 10},
		{"if (1 > 10) { return 10; } else { return 20; }", 20},
		{"if (1 > 10) { return 10; }", nil},
	}

	for _, subtest := range tests {
		t.Run(subtest.input, func(t *testing.T) {
			evaluated := testParseInput(subtest.input)
			if subtest.expected == nil {
				require.IsType(t, &object.Null{}, evaluated)
			} else {
				require.IsType(t, &object.Integer{}, evaluated)
				integer := evaluated.(*object.Integer)
				require.Equal(t, int64(subtest.expected.(int)), integer.Value)
			}
		})
	}
}

func TestEval_ReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`
if (10 > 1) {
    if (50 > 10) {
        return 10;
    }
    return 1;
}
`,
			10,
		},
		{
			`
if (10 > 1) {
    if (5 > 10) {
        return 10;
    }
    return 1;
}
`,
			1,
		},
		{
			`
if (10 > 1) {
    if (5 > 10) {
        return 10;
    } 
}
return 5;
`,
			5,
		},
	}

	for _, subtest := range tests {
		t.Run(subtest.input, func(t *testing.T) {
			evaluated := testParseInput(subtest.input)
			require.IsType(t, &object.Integer{}, evaluated)
			integer := evaluated.(*object.Integer)
			require.Equal(t, int64(subtest.expected.(int)), integer.Value)
		})
	}
}

func TestEval_ErrorHandling(t *testing.T) {
	var (
		Bool = object.TypeBoolean
		Int  = object.TypeInteger
	)
	tests := []struct {
		input    string
		expected string
	}{
		{
			"5 + true",
			fmt.Sprintf("type mismatch: %s + %s", Int, Bool),
		},
		{
			"5 + true; 5",
			fmt.Sprintf("type mismatch: %s + %s", Int, Bool),
		},
		{
			"-true",
			fmt.Sprintf("unknown operator: -%s", Bool),
		},
		{
			"true + false",
			fmt.Sprintf("invalid operation: %s + %s", Bool, Bool),
		},
		{
			"if (true) { true + false }",
			fmt.Sprintf("invalid operation: %s + %s", Bool, Bool),
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
	}

	for _, subtest := range tests {
		t.Run(subtest.input, func(i *testing.T) {
			evaluated := testParseInput(subtest.input)
			require.IsType(t, &object.Error{}, evaluated)
			err := evaluated.(*object.Error)
			require.Equal(t, subtest.expected, err.Inspect())
		})
	}
}

func TestEval_LetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, subtest := range tests {
		t.Run(subtest.input, func(t *testing.T) {
			evaluated := testParseInput(subtest.input)
			require.IsType(t, &object.Integer{}, evaluated)
			integer := evaluated.(*object.Integer)
			require.Equal(t, subtest.expected, integer.Value)
		})
	}
}

func TestEval_FunctionObject(t *testing.T) {
	tests := []struct {
		input      string
		parameters []string
	}{
		{"fn(x) { x + 2; };", []string{"x"}},
	}

	for _, subtest := range tests {
		t.Run(subtest.input, func(t *testing.T) {
			obj := testParseInput(subtest.input)
			require.IsType(t, &object.Function{}, obj)
			fn := obj.(*object.Function)
			params := make([]string, 0, len(fn.Parameters))
			for _, p := range fn.Parameters {
				params = append(params, p.Value)
			}
			require.Equal(t, subtest.parameters, params)
		})
	}
}

func TestEval_FunctionObjectCall(t *testing.T) {
	fib := `
let fib = fn(x) {
  if (x < 2) {
    return x;
  }
  return fib(x - 1) + fib(x - 2);
}
`
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"fn(x) { return x + 2; }(2);", 4},
		{"let identity = fn(x) { return x; }; identity(2)", 2},
		{"let double = fn(x) { return x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { return x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { return x + y; }; add(5, add(5, 5));", 15},
		{fmt.Sprintf("%s; fib(0);", fib), 0},
		{fmt.Sprintf("%s; fib(1);", fib), 1},
		{fmt.Sprintf("%s; fib(2);", fib), 1},
		{fmt.Sprintf("%s; fib(3);", fib), 2},
		{fmt.Sprintf("%s; fib(4);", fib), 3},
		{fmt.Sprintf("%s; fib(5);", fib), 5},
		{`fn(x) { return x; }("string")`, "string"},
	}

	for _, subtest := range tests {
		t.Run(subtest.input, func(t *testing.T) {
			obj := testParseInput(subtest.input)
			switch obj := obj.(type) {
			case *object.Integer:
				require.Equal(t, int64(subtest.expected.(int)), obj.Value)
			case *object.String:
				require.Equal(t, subtest.expected, obj.Value)
			default:
				t.Fatalf("unknown object type %T", obj)
			}
		})
	}
}

func TestEval_Closers(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`
let newAddr = fn(x) {
  return fn(y) { return x + y; };
}

let addTwo = newAddr(2);
addTwo(2);
`,

			int64(4),
		},
		{
			`
let counter = fn(x) {
  if (x > 100) {
    return true;
  } else {
    let foobar = 9999;
    return counter(x + 1);
  }
}
counter(0);
`,
			true,
		},
	}
	for _, subtest := range tests {
		t.Run(subtest.input, func(t *testing.T) {
			obj := testParseInput(subtest.input)
			switch obj := obj.(type) {
			case *object.Integer:
				require.Equal(t, subtest.expected, obj.Value)
			case *object.Boolean:
				require.Equal(t, subtest.expected, obj.Value)
			default:
				t.Fatalf("unknown type %T", obj)
			}
		})
	}
}


func TestEval_Builtins(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`[len(""), len("123456")]`, []interface{}{0, 6}},
		{`len("string")`, 6},
		{`len("")`, 0},
	}

	for _, subtest := range tests {
		t.Run(subtest.input, func(t *testing.T) {
			obj := testParseInput(subtest.input)
			testResult(t, obj, subtest.expected)
		})
	}
}

func TestEval_IndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`[1][0]`, 1},
		{`"string"[0]`, "s"},
		{`[1, 2, 3, 4, 5][3]`, 4},
		{`[1, 2, 3, 4, 5][3]`, 4},
		{`[1, 2, 3, 4, 5][3*2/3]`, 3},
		{`[1, 2, 3, 4, 5][fn(x){ return x; }(2)]`, 3},
		{`[1, 2, 3, 4, 5][6]`, "index out of range"},
		{`let x = [1, 2, 3*4]; let y = 1; x[y + 1]`, 12},
		{`[1, 2, 3, 4, 5]["6"]`, "expected integer, got str"},
		{`[1, 2, 3, 4, 5][fn(x){ return x; }(1)]`, 2},
	}

	for _, subtest := range tests {
		t.Run(subtest.input, func(t *testing.T) {
			obj := testParseInput(subtest.input)
			testResult(t, obj, subtest.expected)
		})
	}
}

func testResult(t *testing.T, obj object.Object, expected interface{}) {
	switch obj := obj.(type) {
	case *object.Integer:
		require.Equal(t, int64(expected.(int)), obj.Value)
	case *object.Boolean:
		require.Equal(t, expected, obj.Value)
	case *object.String:
		require.Equal(t, expected, obj.Value)
	case *object.Error:
		require.Contains(t, obj.Message, expected)
	case *object.List:
		exp, ok := expected.([]interface{})
		if !ok {
			require.FailNow(t, "expected should be a list")
		}
		for k := range obj.Values {
			testResult(t, obj.Values[k], exp[k])
		}
	default:
		t.Fatalf("unknown type %T", obj)
	}
}
