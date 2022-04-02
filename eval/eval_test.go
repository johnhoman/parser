package eval

import (
	"testing"

	"github.com/stretchr/testify/require"

	"mitchlang/lexer"
	"mitchlang/object"
	"mitchlang/parser"
)

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
		l := lexer.New(subtest.input)
		p := parser.New(l)
		program := p.ParseProgram()
		obj := Eval(program)
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
	}

	for _, subtest := range tests {
		l := lexer.New(subtest.input)
		p := parser.New(l)
		program := p.ParseProgram()
		obj := Eval(program)
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
		l := lexer.New(subtest.input)
		p := parser.New(l)
		program := p.ParseProgram()

		evaluated := Eval(program)
		require.IsType(t, &object.Boolean{}, evaluated)
		boolean := evaluated.(*object.Boolean)
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
		l := lexer.New(subtest.input)
		p := parser.New(l)
		program := p.ParseProgram()

		evaluated := Eval(program)
		require.IsType(t, &object.Integer{}, evaluated)
		integer := evaluated.(*object.Integer)
		require.Equal(t, subtest.expected, integer.Value)
	}
}
