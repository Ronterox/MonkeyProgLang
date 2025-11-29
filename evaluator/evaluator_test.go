package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"Operation + between INTEGER and BOOLEAN not implemented!",
		},
		{
			"5 + true; 5;",
			"Operation + between INTEGER and BOOLEAN not implemented!",
		},
		{
			"-true",
			"Not implemented - for BOOLEAN",
		},
		{
			"true + false;",
			"Operation + between BOOLEAN and BOOLEAN not implemented!",
		},
		{
			"5; true + false; 5",
			"Operation + between BOOLEAN and BOOLEAN not implemented!",
		},
		{
			"if (10 > 1) { true + false; }",
			"Operation + between BOOLEAN and BOOLEAN not implemented!",
		},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}
				return 1;
			}
			`,
			"Operation + between BOOLEAN and BOOLEAN not implemented!",
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestEvalIfExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"if (true) { 5 }", 5},
		{"if (false) { 2 }", nil},
		{"if (1 < 2) { 10 }", 10},
		{"if (1) { 10 * 2 }", 20},
		// {"if (null) { 3 }", nil},
		{"if (0) { 2 } else { 5 }", 5},
		{"if (0 + 1) { 3 } else { 5 }", 3},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testInteger(t, evaluated, int64(integer))
		} else {
			testNull(t, evaluated)
		}
	}
}

func testNull(t *testing.T, evaluated object.Object) bool {
	if evaluated != NULL {
		t.Errorf("expected Null got %T=(%v)", evaluated, evaluated)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return 10;
				}
				return 1;
			}
			`, 10,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testInteger(t, evaluated, tt.expected)
	}
}

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
		{"5 + 5", 10},
		{"10 / 2", 5},
		{"2 * 3", 6},
		{"2 - 1", 1},
		{"(2 + 5) * 2", 14},
		{"(2 * 2) * 2", 8},
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
		{"5 > 5", false},
		{"5 == 5", true},
		{"5 < 5", false},
		{"2 < 3", true},
		{"3 > 2", true},
		{"1 != 0", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) != false", true},
		{"true != false", true},
		{"false == true", false},
		{"true == true", true},
		{"false == false", true},
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
