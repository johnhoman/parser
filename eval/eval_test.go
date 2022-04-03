package eval

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	"mitchlang/lexer"
	"mitchlang/object"
	"mitchlang/parser"
)

func parseInput(in string) object.Object {
	l := lexer.New(in)
	p := parser.New(l)
	env := object.NewEnv()
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
		obj := parseInput(subtest.input)
		require.IsType(t, &object.Integer{}, obj)
		integer := obj.(*object.Integer)
		require.Equal(t, subtest.expected, integer.Value)
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
		obj := parseInput(subtest.input)
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
		obj := parseInput(subtest.input)
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
		evaluated := parseInput(subtest.input)
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
		{"if (true) { 10 }", 10},
		{"if (false) { 10 } else { 20 }", 20},
		{"if (10 == 10) { 10 }", 10},
		{"if (10 != 10) { 10 } else { 20 }", 20},
		{"if (1 < 10) { 10 }", 10},
		{"if (1 > 10) { 10 } else { 20 }", 20},
		{"if (1 > 10) { 10 }", nil},
	}

	for _, subtest := range tests {
		t.Run(subtest.input, func(t *testing.T) {
			evaluated := parseInput(subtest.input)
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
			evaluated := parseInput(subtest.input)
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
			evaluated := parseInput(subtest.input)
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
			evaluated := parseInput(subtest.input)
			require.IsType(t, &object.Integer{}, evaluated)
			integer := evaluated.(*object.Integer)
			require.Equal(t, subtest.expected, integer.Value)
		})
	}
}
