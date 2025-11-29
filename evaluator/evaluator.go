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
			if ret, ok := result.(*object.Return); ok {
				return ret.Value
			}
		}
		return result
	case *ast.BlockStatement:
		var result object.Object
		for _, stmt := range node.Statements {
			result = Eval(stmt)
			if _, ok := result.(*object.Return); ok {
				return result
			}
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
				return NULL
			}
		case token.MINUS:
			switch right := right.(type) {
			case *object.Integer:
				return &object.Integer{Value: -right.Value}
			}
			return NULL
		}

		fmt.Printf("Not implemented operator %s!\n", node.Operator)
		return NULL
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)

		if left.Type() == object.INTEGER && right.Type() == object.INTEGER {
			left := left.(*object.Integer)
			right := right.(*object.Integer)

			switch node.Operator {
			case token.PLUS:
				return &object.Integer{Value: left.Value + right.Value}
			case token.ASTERISK:
				return &object.Integer{Value: left.Value * right.Value}
			case token.MINUS:
				return &object.Integer{Value: left.Value - right.Value}
			case token.SLASH:
				return &object.Integer{Value: left.Value / right.Value}
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
		} else if left.Type() == object.BOOLEAN && right.Type() == object.BOOLEAN {
			left := left.(*object.Boolean)
			right := right.(*object.Boolean)

			switch node.Operator {
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
		}

		fmt.Printf("Operation %s between %s and %s not implemented!\n", node.Operator, left.Type(), right.Type())
		return NULL
	case *ast.IfExpression:
		cond := Eval(node.Condition)
		switch cond := cond.(type) {
		case *object.Boolean:
			if cond == TRUE {
				return Eval(node.Consequence)
			} else if node.Alternative != nil {
				return Eval(node.Alternative)
			}
		case *object.Integer:
			if cond.Value > 0 {
				return Eval(node.Consequence)
			}
			return Eval(node.Alternative)
		}
		return NULL
	case *ast.ReturnStatement:
		return &object.Return{Value: Eval(node.RetValue)}
	}
	fmt.Printf("Not implemented eval for %T!\n", node)
	return NULL
}
