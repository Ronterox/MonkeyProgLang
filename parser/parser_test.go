package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func parseSingleInputProgram(t *testing.T, input string) *ast.Program {
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseErrors(t, p)
	return program
}

func parseSingleStatement(t *testing.T, input string) *ast.ExpressionStatement {
	program := parseSingleInputProgram(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected ExpressionStatement, got %T", program.Statements[0])
	}

	return stmt
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`
	stmt := parseSingleStatement(t, input)

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}

		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}
		testFunc(value)
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	stmt := parseSingleStatement(t, "{}")

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
}

func TestParsingHashIntegerLiterals(t *testing.T) {
	stmt := parseSingleStatement(t, `{ 1: 20, 2: 4 }`)

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 2 {
		t.Fatalf("len(hash.Pairs) not 3. got=%d", len(hash.Pairs))
	}

	expected := map[int64]int64{
		1: 20,
		2: 4,
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.IntegerLiteral)
		if !ok {
			t.Errorf("key is not a IntegerLiteral got %T", key)
		}
		testIntegerLiteral(t, value, expected[literal.Value])
	}
}

func TestParsingHashBooleanLiterals(t *testing.T) {
	stmt := parseSingleStatement(t, `{ true: 20, false: 4 }`)

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 2 {
		t.Fatalf("len(hash.Pairs) not 3. got=%d", len(hash.Pairs))
	}

	expected := map[bool]int64{
		true:  20,
		false: 4,
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.BooleanLiteral)
		if !ok {
			t.Errorf("key is not a BooleanLiteral got %T", key)
		}
		testIntegerLiteral(t, value, expected[literal.Value])
	}
}

func TestParsingHashStringLiterals(t *testing.T) {
	stmt := parseSingleStatement(t, `{ "hi": 2, "this": 10, "it": 7 }`)

	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Fatalf("len(hash.Pairs) not 3. got=%d", len(hash.Pairs))
	}

	expected := map[string]int64{
		"hi":   2,
		"this": 10,
		"it":   7,
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not a StringLiteral got %T", key)
		}
		testIntegerLiteral(t, value, expected[literal.Value])
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	stmt := parseSingleStatement(t, "[1, 2 * 2, 3 + 3]")

	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	stmt := parseSingleStatement(t, "myArray[1 + 1]")

	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}
	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestIndentifierExpression(t *testing.T) {
	stmt := parseSingleStatement(t, "foobar;")
	testLiteralExpression(t, stmt.Expression, "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	stmt := parseSingleStatement(t, "5;")
	if !testLiteralExpression(t, stmt.Expression, 5) {
		return
	}
}

func TestStringLiteralExpression(t *testing.T) {
	stmt := parseSingleStatement(t, `"foobar";`)
	if !testLiteralExpression(t, stmt.Expression, `"foobar"`) {
		return
	}

	stmt = parseSingleStatement(t, `"foo" "bar";`)
	if !testLiteralExpression(t, stmt.Expression, `"foobar"`) {
		return
	}

	stmt = parseSingleStatement(t, `"a\n";`)
	if !testLiteralExpression(t, stmt.Expression, `"a
"`) {
		return
	}
}

func TestTemplateExpression(t *testing.T) {
	stmt := parseSingleStatement(t, "`a $literal template\\n`")
	exp, ok := stmt.Expression.(*ast.TemplateString)
	if !ok {
		t.Fatalf("expected TemplateString got %T", stmt.Expression)
	}

	if len(exp.Elements) != 3 {
		t.Fatalf("expected 3 elements got %d", len(exp.Elements))
	}

	testLiteralExpression(t, exp.Elements[0], `"a "`)
	testLiteralExpression(t, exp.Elements[1], "literal")
	testLiteralExpression(t, exp.Elements[2], `" template\n"`)

	stmt = parseSingleStatement(t, "`$literal a`")
	exp, ok = stmt.Expression.(*ast.TemplateString)
	if !ok {
		t.Fatalf("expected TemplateString got %T", stmt.Expression)
	}

	if len(exp.Elements) != 3 {
		t.Fatalf("expected 3 elements got %d", len(exp.Elements))
	}

	testLiteralExpression(t, exp.Elements[0], `""`)
	testLiteralExpression(t, exp.Elements[1], "literal")
	testLiteralExpression(t, exp.Elements[2], `" a"`)
}

func TestBooleanExpression(t *testing.T) {
	stmt := parseSingleStatement(t, "true;")
	if !testLiteralExpression(t, stmt.Expression, true) {
		return
	}
}

func TestIfExpression(t *testing.T) {
	stmt := parseSingleStatement(t, "if (x < y) { x }")

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Errorf("expected IfExpression, got %T", exp)
		return
	}

	if exp.TokenLiteral() != "if" {
		t.Errorf("expected TokenLiteral to be if, got %s", exp.TokenLiteral())
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("expected Consequence.Statements to be 1, got %d", len(exp.Consequence.Statements))
	}

	blockStmt, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("expected ExpressionStatement got %T", blockStmt)
		return
	}

	if !testLiteralExpression(t, blockStmt.Expression, "x") {
		return
	}
}

func TestIfElseExpression(t *testing.T) {
	stmt := parseSingleStatement(t, "if (x < y) { x } else { y }")

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Errorf("expected IfExpression, got %T", exp)
		return
	}

	if exp.TokenLiteral() != "if" {
		t.Errorf("expected TokenLiteral to be if, got %s", exp.TokenLiteral())
		return
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("expected Consequence.Statements to be 1, got %d", len(exp.Consequence.Statements))
	}

	blockIfStmt, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("expected ExpressionStatement got %T", blockIfStmt)
		return
	}

	if !testLiteralExpression(t, blockIfStmt.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("expected Alternative.Statements to be 1, got %d", len(exp.Alternative.Statements))
	}

	blockIfElseStmt, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("expected ExpressionStatement got %T", blockIfElseStmt)
		return
	}

	if !testLiteralExpression(t, blockIfElseStmt.Expression, "y") {
		return
	}
}

func TestFunctionLiteral(t *testing.T) {
	stmt := parseSingleStatement(t, "fn(x, y) { x + y }")

	exp, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Errorf("expected FunctionLiteral, got %T", exp)
		return
	}

	if exp.TokenLiteral() != "fn" {
		t.Errorf("expected TokenLiteral to be fn, got %s", exp.TokenLiteral())
		return
	}

	if len(exp.Parameters) != 2 {
		t.Errorf("expected length of parameters to be 2, got %d", len(exp.Parameters))
		return
	}

	if !testLiteralExpression(t, exp.Parameters[0], "x") {
		return
	}
	if !testLiteralExpression(t, exp.Parameters[1], "y") {
		return
	}

	if len(exp.Body.Statements) != 1 {
		t.Errorf("expected Body.Statements to have length of 1, got %d", len(exp.Body.Statements))
		return
	}

	bodyStmt, ok := exp.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("expected ExpressionStatement got %T", bodyStmt)
		return
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestParsingMacro(t *testing.T) {
	stmt := parseSingleStatement(t, "macro(x: int, y: int) { `x + y` }")

	exp, ok := stmt.Expression.(*ast.MacroLiteral)
	if !ok {
		t.Errorf("expected Macro, got %T", exp)
		return
	}

	if exp.TokenLiteral() != "macro" {
		t.Errorf("expected TokenLiteral to be macro, got %s", exp.TokenLiteral())
		return
	}

	if len(exp.Parameters) != 2 {
		t.Errorf("expected length of parameters to be 2, got %d", len(exp.Parameters))
		return
	}

	if !testLiteralExpression(t, exp.Parameters[0], "x") {
		return
	}

	if !testLiteralExpression(t, exp.Pattern[0], "int") {
		return
	}

	if !testLiteralExpression(t, exp.Parameters[1], "y") {
		return
	}

	if !testLiteralExpression(t, exp.Pattern[1], "int") {
		return
	}

	if len(exp.Body.Elements) != 1 {
		t.Fatalf("expected 1 elements got %d", len(exp.Body.Elements))
	}

	testLiteralExpression(t, exp.Body.Elements[0], `"x + y"`)

	// TODO: Write the rest of the tests
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
		stmt := parseSingleStatement(t, tt.input)
		exp := stmt.Expression.(*ast.FunctionLiteral)

		if len(exp.Parameters) != len(tt.expectedParams) {
			t.Errorf("expected length of parameters to be %d, got %d", len(tt.expectedParams), len(exp.Parameters))
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
	stmt := parseSingleStatement(t, "add(1, 2 * 3, 4 + 5)")

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Errorf("expected CallExpression, got %T", exp)
		return
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Errorf("expected length of parameters to be 3, got %d", len(exp.Arguments))
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
		stmt := parseSingleStatement(t, tt.input)

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Errorf("expected PrefixExpression, got %T", stmt.Expression)
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("expected operator to be '%s', got '%s'", exp.Operator, tt.operator)
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
		{"5 >= 5;", 5, ">=", 5},
		{"5 <= 5;", 5, "<=", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		stmt := parseSingleStatement(t, tt.input)
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
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"a | add(b & c) | d",
			"((a | add((b & c))) | d)",
		},
		{
			"5 > 4 == 3 < 4 & c | b",
			"((((5 > 4) == (3 < 4)) & c) | b)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
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
			"a + b % c",
			"(a + (b % c))",
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
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}
	for _, tt := range tests {
		program := parseSingleInputProgram(t, tt.input)
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
		if v[0] == '"' || v[0] == '`' {
			return testString(t, exp, v)
		}
		return testIdentifier(t, exp, v)
	case bool:
		return testBoolean(t, exp, v)
	}
	t.Errorf("type of exp not handled got %T", exp)
	return false
}

func testString(t *testing.T, exp ast.Expression, v string) bool {
	lit, ok := exp.(*ast.StringLiteral)
	if !ok {
		t.Errorf("expected String got %T", exp)
		return false
	}

	if fmt.Sprintf(`"%s"`, lit.Value) != v {
		t.Errorf("expected String '%s' got '%s'", v, fmt.Sprintf(`"%s"`, lit.Value))
		return false
	}

	return true
}

func testBoolean(t *testing.T, exp ast.Expression, v bool) bool {
	opExp, ok := exp.(*ast.BooleanLiteral)
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
	tests := []struct {
		input         string
		expectedIdent string
		expectedValue any
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		program := parseSingleInputProgram(t, tt.input)

		if len(program.Statements) != 1 {
			t.Errorf("Expected 1 statements, got %d", len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdent) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue any
	}{
		{"return 5;", 5},
		{"return 10;", 10},
		{"return add;", "add"},
	}

	for _, tt := range tests {
		program := parseSingleInputProgram(t, tt.input)

		if len(program.Statements) != 1 {
			t.Errorf("Expected 1 statements, got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got %T", stmt)
			continue
		}

		if stmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got %q", stmt.TokenLiteral())
		}

		if !testLiteralExpression(t, stmt.RetValue, tt.expectedValue) {
			return
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
