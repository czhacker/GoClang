package evaluator

import (
	"GoClang/lexer"
	"GoClang/object"
	"GoClang/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T)  {
	tests := []struct{
		input string
		expected int64
	}{
		{"5",5},
		{"10", 10},
		{"-5",-5},
		{"-10",-10},
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

func TestEvalBooleanExpression(t *testing.T)  {
	tests := []struct{
		input string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{" 1 == 1", true},
		{"1 != 1",false},
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

func TestBangOperator(t *testing.T)  {
	tests := []struct{
		input string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5",true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T)  {
	tests := []struct{
		input string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }",10},
		{"if (1 > 2) {10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _,tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		}else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T)  {
	tests := []struct{
		input string
		expected int64
	}{
		{"return 10;",10},
		{"1; return 10; 9", 10},
		{"return 2 * 5; 1", 10},
		{"9; return 2 * 5; 3", 10},
		{"if (10 > 1) { if ( 10 > 1) { return 10; } return 1; }", 10},
	}

	for _, tt := range tests{
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T)  {
	tests := []struct{
		input string
		expectedMessage string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true","unknown operator: -BOOLEAN"},
		{"true + false;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false; }", "unknown operator: BOOLEAN + BOOLEAN"},
		{`if(10 > 1){ if(10 > 1){ return true + false; } return 1;}`,"unknown operator: BOOLEAN + BOOLEAN"},
		{"foobar","identifier not found: foobar"},
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

func TestLetStatements(t *testing.T)  {
	tests := []struct{
		input string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;",15},
	}

	for _, tt := range tests{
		testIntegerObject(t, testEval(tt.input),tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParserProgram()
	env := object.NewEnviroment()
	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool{
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

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool  {
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
	if obj != NULL{
		t.Errorf("object is not NULL. got=%T(+%v)", obj,obj)
		return false
	}
	return true
}


