package evaluator

import (
	"GoClang/lexer"
	"GoClang/object"
	"GoClang/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 + 10 + 10", 40},
		{"2 * 2 * 2", 8},
		{"5 + 2 * 2", 9},
		{"20 + 2 * -10 + 5", 5},
		{"3 * (3 * 3) + 10", 37},
		{"10 / 2 + 3 * 10", 35},
		{"50 / 2 * 2", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{" 1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"true != false", true},
		{"false == false", true},
		{"false == true", false},
		{"true != true", false},
		{"(1 > 2) == true", false},
		{"(1 < 2) == true", true},
		{"(1 > 2) == false", true},
		{"(1 < 2) == false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) {10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"1; return 10; 9", 10},
		{"return 2 * 5; 1", 10},
		{"9; return 2 * 5; 3", 10},
		{"if (10 > 1) { if ( 10 > 1) { return 10; } return 1; }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + false;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false; }", "unknown operator: BOOLEAN + BOOLEAN"},
		{`if(10 > 1){ if(10 > 1){ return true + false; } return 1;}`, "unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar", "identifier not found: foobar"},
		{`"hello" - "world"`, "unknown operator: STRING - STRING"},
		{`{"name": "Monkey"}[fn(x){ x }];`, "unusable as hash key: FUNCTION"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T (%+v)", errObj, errObj)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
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
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	tests := "fn(x) { x + 2; };"
	evaluated := testEval(tests)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].Value != "x" {
		t.Fatalf("parameters is not 'x'. got=%q", fn.Parameters[0].Value)
	}

	expectedBody := "(x + 2)"
	if expectedBody != fn.Body.String() {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let fns = fn(x) { x;}; fns(5);", 5},
		{"let fns = fn(x) { return x;}; fns(5);", 5},
		{"let double = fn(x) { x * 2;}; double(2);", 4},
		{"let add = fn(x,y) { x + y;}; add(2,3);", 5},
		{"let add = fn(x,y) { x + y;}; add(5 + 5, add(5 , 5));", 20},
		{"fn(x){x;}(5);", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
	let newAdder = fn(x) {
		fn(y){ x + y };
	};

	let addTwo = newAdder(2);
	addTwo(3);`
	testIntegerObject(t, testEval(input), 5)
}

func TestString(t *testing.T) {
	input := `"Hello World";`

	evaluated := testEval(input)
	strObj, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T", evaluated)
	}

	if strObj.Value != "Hello World" {
		t.Fatalf("String.Value is not %q. got=%q", "Hello World", strObj.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	strObj, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T", evaluated)
	}

	if strObj.Value != "Hello World!" {
		t.Fatalf("String.Value is not %q. got=%q", "Hello World", strObj.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not support, got=INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))

		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}

		}
	}
}

func TestArrayLiteral(t *testing.T) {
	input := "[1,2*2, 3 + 3]"
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array, got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array have wrong number of elements. got=%d", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], int64(1))
	testIntegerObject(t, result.Elements[1], int64(4))
	testIntegerObject(t, result.Elements[2], int64(6))
}

func TestArrayIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1,2,3][0]", 1},
		{"[1,2,3][1]", 2},
		{"[1,2,3][2]", 3},
		{"let i = 0; [1][i];", 1},
		{"[1,2,3][1 + 1]", 3},
		{"let myArray = [1, 2, 3]; myArray[0]", 1},
		{"[1,2,3][3]", nil},
		{"[1,2,3][-1]", nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}

}

func TestHashLiterals(t *testing.T)  {
	input := `let two = "two";
	{
		"one" : 10 - 9,
		two : 1 + 1,
		"thr" + "ee" : 6 / 2,
		4 : 4,
		true : 5,
		false : 6
	}`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash, got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value:"one"}).HashKey() : 1,
		(&object.String{Value:"two"}).HashKey() : 2,
		(&object.String{Value:"three"}).HashKey() : 3,
		(&object.Integer{Value:4}).HashKey() : 4,
		(&object.Boolean{Value:true}).HashKey() : 5,
		(&object.Boolean{Value:false}).HashKey() : 6,
	}

	if len(expected) != len(result.Pair) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pair))
	}

	for expectedKey, expectedValue := range expected{
		pair, ok := result.Pair[expectedKey]
		if !ok {
			t.Error("no pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T)  {
	tests := []struct{
		input string
		expected interface{}
	}{
		{`{"foo": 5}["foo"]`,5},
		{`{"foo": 5}["bar"]`, nil},
		{`let key = "foo"; {"foo": 5}[key]`, 5},
		{`{}["foo"]`, nil},
		{`{5:5}[5]`, 5},
		{`{true:5}[true]`, 5},
		{`{false: 5}[false]`, 5},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		}else {
			testNullObject(t, evaluated)
		}
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParserProgram()
	env := object.NewEnviroment()
	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T(+%v)", obj, obj)
		return false
	}
	return true
}
