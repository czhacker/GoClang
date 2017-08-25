package evaluator

import (
	"GoClang/ast"
	"GoClang/object"
	"fmt"
)

var (
	TRUE = &object.Boolean{Value:true}
	FALSE = &object.Boolean{Value:false}
	NULL = &object.Null{}
)


func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.IntegerLiteral:
		return &object.Integer{Value:node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)

	case *ast.BlockStatement:
		return evalBlockStatements(node)

	case *ast.IfExpression:
		return evalIfExpression(node)

	case *ast.ReturnStatement:
		value := Eval(node.ReturnValue)
		return &object.ReturnValue{Value:value}

	}
	return nil
}
func evalProgram(node *ast.Program) object.Object {
	var result object.Object

	for _, statement := range node.Statements {
		result = Eval(statement)

		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}
	return result
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)

		if returnValue, ok := result.(*object.ReturnValue); ok{
			return returnValue.Value
		}
	}
	return result
}

func evalBlockStatements(node *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range node.Statements {
		result = Eval(statement)

		if result != NULL && result.Type() == object.RETURN_VALUE_OBJ {
			return result
		}
	}
	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalPrefixExpression(operator string, obj object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(obj)
	case "-":
		return evalMinusOperatorExpression(obj)
	default:
		return NULL
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBooleanObject(left == right)
	case operator == "!=":
		return nativeBoolToBooleanObject(left != right)
	default:
		return NULL
	}
}

func evalBangOperatorExpression(obj object.Object) object.Object {
	switch obj {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusOperatorExpression(obj object.Object) object.Object  {
	result, ok := obj.(*object.Integer)
	if !ok {
		return NULL
	}
	return &object.Integer{Value:-result.Value}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value:leftValue + rightValue}
	case "-":
		return &object.Integer{Value:leftValue - rightValue}
	case "*":
		return &object.Integer{Value:leftValue * rightValue}
	case "/":
		if rightValue == 0 {
			fmt.Printf("DIVIDEND ZERO IS ILLEGAL!\n")
			return NULL
		}
		return &object.Integer{Value:leftValue / rightValue}
	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return NULL
	}
}


func evalIfExpression(node *ast.IfExpression) object.Object {
	condition := Eval(node.Condition)
	if isTruthy(condition){
		return Eval(node.Consequence)
	}else if node.Alternative != nil{
		return Eval(node.Alternative)
	}else{
		return NULL
	}
}

func isTruthy(condition object.Object) bool {
	switch condition {
	case FALSE:
		return false
	case NULL:
		return false
	case TRUE:
		return true
	default:
		return true
	}
}

