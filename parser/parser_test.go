package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestIndentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected ExpressionStatement, got %T", program.Statements[0])
		return
	}

	testLiteralExpression(t, stmt.Expression, "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected ExpressionStatement, got %T", program.Statements[0])
		return
	}

	if !testLiteralExpression(t, stmt.Expression, 5) {
		return
	}
}

func TestBooleanExpression(t *testing.T) {
	input := "true;"
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected ExpressionStatement, got %T", program.Statements[0])
		return
	}

	if !testLiteralExpression(t, stmt.Expression, true) {
		return
	}
}

func TestIfExpression(t *testing.T) {
	input := "if (x < y) { x }"
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected ExpressionStatement, got %T", program.Statements[0])
		return
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Errorf("Expected IfExpression, got %T", exp)
		return
	}

	if exp.TokenLiteral() != "if" {
		t.Errorf("Expected TokenLiteral to be if, got %s", exp.TokenLiteral())
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("Expected Consequence.Statements to be 1, got %d", len(exp.Consequence.Statements))
	}

	blockStmt, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected ExpressionStatement got %T", blockStmt)
		return
	}

	if !testLiteralExpression(t, blockStmt.Expression, "x") {
		return
	}
}

func TestIfElseExpression(t *testing.T) {
	input := "if (x < y) { x } else { y }"
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected ExpressionStatement, got %T", program.Statements[0])
		return
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Errorf("Expected IfExpression, got %T", exp)
		return
	}

	if exp.TokenLiteral() != "if" {
		t.Errorf("Expected TokenLiteral to be if, got %s", exp.TokenLiteral())
		return
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("Expected Consequence.Statements to be 1, got %d", len(exp.Consequence.Statements))
	}

	blockIfStmt, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected ExpressionStatement got %T", blockIfStmt)
		return
	}

	if !testLiteralExpression(t, blockIfStmt.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("Expected Alternative.Statements to be 1, got %d", len(exp.Alternative.Statements))
	}

	blockIfElseStmt, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected ExpressionStatement got %T", blockIfElseStmt)
		return
	}

	if !testLiteralExpression(t, blockIfElseStmt.Expression, "y") {
		return
	}
}

func TestFunctionLiteral(t *testing.T) {
	input := "fn(x, y) { x + y }"
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected ExpressionStatement, got %T", program.Statements[0])
		return
	}

	exp, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Errorf("Expected FunctionLiteral, got %T", exp)
		return
	}

	if exp.TokenLiteral() != "fn" {
		t.Errorf("Expected TokenLiteral to be fn, got %s", exp.TokenLiteral())
		return
	}

	if len(exp.Parameters) != 2 {
		t.Errorf("Expected length of parameters to be 2, got %d", len(exp.Parameters))
		return
	}

	if !testLiteralExpression(t, exp.Parameters[0], "x") {
		return
	}
	if !testLiteralExpression(t, exp.Parameters[1], "y") {
		return
	}

	if len(exp.Body.Statements) != 1 {
		t.Errorf("Expected Body.Statements to have length of 1, got %d", len(exp.Body.Statements))
		return
	}

	bodyStmt, ok := exp.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected ExpressionStatement got %T", bodyStmt)
		return
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameters(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParseErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		exp := stmt.Expression.(*ast.FunctionLiteral)

		if len(exp.Parameters) != len(tt.expectedParams) {
			t.Errorf("Expected length of parameters to be %d, got %d", len(tt.expectedParams), len(exp.Parameters))
			return
		}

		for i, p := range tt.expectedParams {
			if !testLiteralExpression(t, exp.Parameters[i], p) {
				return
			}
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5)"
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("Expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected ExpressionStatement, got %T", program.Statements[0])
		return
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Errorf("Expected CallExpression, got %T", exp)
		return
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Errorf("Expected length of parameters to be 3, got %d", len(exp.Arguments))
		return
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    any
	}{
		{"!true", "!", true},
		{"!false", "!", false},
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Errorf("Expected 1 statement, got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("Expected ExpressionStatement, got %T", program.Statements[0])
			return
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Errorf("Expected PrefixExpression, got %T", stmt.Expression)
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("Expected operator to be '%s', got '%s'", exp.Operator, tt.operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input    string
		left     any
		operator string
		right    any
	}{
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParseErrors(t, p)

		if len(program.Statements) != 1 {
			t.Errorf("Expected 1 statement, got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("Expected ExpressionStatement, got %T", program.Statements[0])
			return
		}

		if !testInfixExpression(t, stmt.Expression, tt.left, tt.operator, tt.right) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParseErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, i int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("Expected ast.IntegerLiteral, got %T", il)
		return false
	}

	if integ.Value != i {
		t.Errorf("Expected int to be %d got %d", integ.Value, i)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", i) {
		t.Errorf("Expected token literal to be %d got %s", i, integ.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBoolean(t, exp, v)
	}
	t.Errorf("type of exp not handled got %T", exp)
	return false
}

func testBoolean(t *testing.T, exp ast.Expression, v bool) bool {
	opExp, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("Expected BooleanExpression, got %T", exp)
		return false
	}

	if opExp.Value != v {
		t.Fatalf("Expected value to be %t, got %t", opExp.Value, v)
	}

	if opExp.TokenLiteral() != fmt.Sprintf("%t", v) {
		t.Fatalf("Expected TokenLiteral to be %s got %t", opExp.TokenLiteral(), v)
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left any, operator string, right any) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("Expected PrefixExpression, got %T", exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Fatalf("Expected operator to be '%s', got '%s'", opExp.Operator, operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 5;
	let foobar = 838383;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if program == nil {
		t.Error("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Errorf("Expected 3 statements, got %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 838383;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParseErrors(t, p)

	if program == nil {
		t.Error("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Errorf("Expected 3 statements, got %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got %T", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.TokenLiteral())
		}
	}
}
func checkParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("Expected let, got %s", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("Expected LetStatement, got %T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("Expected %s, got %s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("Expected %s, got %s", name, letStmt.Name)
		return false
	}

	return true
}
