package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
	"monkey/token"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		var result object.Object
		for _, stmt := range node.Statements {
			result = Eval(stmt)
		}
		return result
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		if node.Value {
			return TRUE
		}
		return FALSE
	case *ast.PrefixExpression:
		right := Eval(node.Right)

		switch node.Operator {
		case token.BANG:
			switch right := right.(type) {
			case *object.Integer:
				if right.Value > 0 {
					return FALSE
				}
				return TRUE
			case *object.Boolean:
				if right == TRUE {
					return FALSE
				}
				return TRUE
			default:
				return &object.Null{}
			}
		case token.MINUS:
			result := &object.Integer{}

			switch right := right.(type) {
			case *object.Integer:
				result.Value = -right.Value
			default:
				return &object.Null{}
			}

			return result
		}

		fmt.Printf("Not implemented operator %s!\n", node.Operator)
		return &object.Null{}
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)

		if left.Type() == object.INTEGER && right.Type() == object.INTEGER {
			left := left.(*object.Integer)
			right := right.(*object.Integer)
			result := &object.Integer{}

			switch node.Operator {
			case token.PLUS:
				result.Value = left.Value + right.Value
			case token.ASTERISK:
				result.Value = left.Value * right.Value
			case token.MINUS:
				result.Value = left.Value - right.Value
			case token.SLASH:
				result.Value = left.Value / right.Value
			case token.GT:
				if left.Value > right.Value {
					return TRUE
				}
				return FALSE
			case token.LT:
				if left.Value < right.Value {
					return TRUE
				}
				return FALSE
			case token.EQ:
				if left.Value == right.Value {
					return TRUE
				}
				return FALSE
			case token.NE:
				if left.Value != right.Value {
					return TRUE
				}
				return FALSE
			}
			return result
		}

		fmt.Printf("Sum between %s and %s not implemented!\n", left.Type(), right.Type())
		return &object.Null{}
	}
	fmt.Printf("Not implemented eval for %T!\n", node)
	return &object.Null{}
}
