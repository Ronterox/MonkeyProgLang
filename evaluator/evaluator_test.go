package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}
	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d", len(result.Elements))
	}
	testInteger(t, result.Elements[0], 1)
	testInteger(t, result.Elements[1], 4)
	testInteger(t, result.Elements[2], 6)
}

func TestArrayIndexing(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"[1, 2 * 2, 3 + 3][0]", 1},
		{"[1, 2 * 2, 3 + 3][1]", 4},
		{"[1, 2 * 2, 3 + 3][2]", 6},
		{"let arr = [1, 2 * 2, 3 + 3]; arr[2];", 6},
		{"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];", 6},
		{"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]", 2},
		{"[1, 2, 3][3]", nil},
		{"[1, 2, 3][-1]", nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if expected, ok := tt.expected.(int); ok {
			testInteger(t, evaluated, int64(expected))
		} else {
			testNull(t, evaluated)
		}
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
		{`first([1, 2])`, 1},
		{`last([1, 2])`, 2},
		{`first([])`, nil},
		{`last([])`, nil},
		{`first(rest([1, 2]))`, 2},
		{`first(rest([2]))`, nil},
		{`
			let arr = [1, 2, 3];
			rest(rest(arr))
			first(arr)
			`,
			1,
		},
		{`
			let arr = [1, 2, 3];
			push(arr, 1)
			last(push(arr, 2))
			`,
			2,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testInteger(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		default:
			testNull(t, evaluated)
		}
	}
}

func TestFunction(t *testing.T) {
	input := "fn(x) { x + 2; };"
	evaluated := testEval(input)

	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("expected Function got %T=(%v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("expected 1 parameter got %d", len(fn.Parameters))
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("expected x as parameter got %s", fn.Parameters[0].String())
	}

	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("expected body to be %s got %s", expectedBody, fn.Body.String())
	}
}

func TestFunctionCall(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testInteger(t, testEval(tt.input), tt.expected)
	}
}

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
		{
			"foobar",
			"identifier not found: foobar",
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

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
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
		{"if (null) { 3 }", nil},
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

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}
	for _, tt := range tests {
		testInteger(t, testEval(tt.input), tt.expected)
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
		{
			`
			let f = fn(x) {
			  return x;
			  x + 10;
			};
			f(10);`,
			10,
		},
		{
			`
			let f = fn(x) {
			   let result = x + 10;
			   return result;
			   return 10;
			};
			f(10);`,
			20,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testInteger(t, evaluated, tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
		let newAdder = fn(x) {
			fn(y) { x + y };
		};
		let addTwo = newAdder(2);
		addTwo(2);
	`
	testInteger(t, testEval(input), 4)
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

func TestEvalString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"foobar"`, "foobar"},
		{`"foo and bar"`, "foo and bar"},
		{`"foo\nand\nbar"`, "foo\\nand\\nbar"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testString(t, evaluated, tt.expected)
	}
}

func testString(t *testing.T, evaluated object.Object, expected string) bool {
	obj, ok := evaluated.(*object.String)
	if !ok {
		t.Errorf("expected String got %T=(%v)", evaluated, evaluated)
		return false
	}

	if obj.Value != expected {
		t.Errorf("expected %s got %s", expected, obj.Value)
		return false
	}
	return true
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
		t.Errorf("expected Integer got %T=(%v)", evaluated, evaluated)
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
	return Eval(program, object.NewEnvironment())
}
