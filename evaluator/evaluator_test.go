package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestEvalInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-10", -10},
		{"-5", -5},
		{"--5", 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testInteger(t, evaluated, tt.expected)
	}
}

func TestEvalBoolean(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"!true", false},
		{"!false", true},
		{"!0", true},
		{"!5", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBoolean(t, evaluated, tt.expected)
	}
}

func testBoolean(t *testing.T, evaluated object.Object, expected bool) bool {
	obj, ok := evaluated.(*object.Boolean)
	if !ok {
		t.Errorf("expected Boolean got %T", evaluated)
		return false
	}

	if obj.Value != expected {
		t.Errorf("expected %t got %t", expected, obj.Value)
		return false
	}
	return true
}

func testInteger(t *testing.T, evaluated object.Object, expected int64) bool {
	obj, ok := evaluated.(*object.Integer)
	if !ok {
		t.Errorf("expected Integer got %T", evaluated)
		return false
	}

	if obj.Value != expected {
		t.Errorf("expected %d got %d", expected, obj.Value)
		return false
	}
	return true
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return Eval(program)
}
