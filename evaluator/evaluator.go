package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
	"monkey/token"
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
		return &object.Boolean{Value: node.Value}
	case *ast.PrefixExpression:
		right := Eval(node.Right)

		switch node.Operator {
		case token.BANG:
			result := &object.Boolean{}

			switch right := right.(type) {
			case *object.Integer:
				result.Value = !(right.Value > 0)
			case *object.Boolean:
				result.Value = !right.Value
			default:
				return &object.Null{}
			}

			return result
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

		fmt.Printf("Not implemented %s!\n", node.Operator)
		return &object.Null{}
	}
	fmt.Printf("Not implemented %T!\n", node)
	return &object.Null{}
}
