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

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}
	return false
}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		var result object.Object
		for _, stmt := range node.Statements {
			result = Eval(stmt, env)
			switch result := result.(type) {
			case *object.Return:
				return result.Value
			case *object.Error:
				return result
			}
		}
		return result
	case *ast.BlockStatement:
		var result object.Object
		for _, stmt := range node.Statements {
			result = Eval(stmt, env)
			switch result.(type) {
			case *object.Return, *object.Error:
				return result
			}
		}
		return result
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BooleanLiteral:
		if node.Value {
			return TRUE
		}
		return FALSE
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

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
				return newError("Not implemented %s for %s", token.BANG, right.Type())
			}
		case token.MINUS:
			switch right := right.(type) {
			case *object.Integer:
				return &object.Integer{Value: -right.Value}
			}
			return newError("Not implemented %s for %s", token.MINUS, right.Type())
		}

		fmt.Printf("Not implemented operator %s!\n", node.Operator)
		return NULL
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

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
		return newError("Operation %s between %s and %s not implemented!", node.Operator, left.Type(), right.Type())
	case *ast.IfExpression:
		cond := Eval(node.Condition, env)
		if isError(cond) {
			return cond
		}

		switch cond := cond.(type) {
		case *object.Boolean:
			if cond == TRUE {
				return Eval(node.Consequence, env)
			} else if node.Alternative != nil {
				return Eval(node.Alternative, env)
			}
		case *object.Integer:
			if cond.Value > 0 {
				return Eval(node.Consequence, env)
			}
			return Eval(node.Alternative, env)
		}
		return NULL
	case *ast.ReturnStatement:
		ret := Eval(node.RetValue, env)
		if isError(ret) {
			return ret
		}
		return &object.Return{Value: ret}
	case *ast.LetStatement:
		value := Eval(node.Value, env)
		if isError(value) {
			return value
		}
		return env.Set(node.Name.Value, value)
	case *ast.Identifier:
		if value, ok := env.Get(node.Value); ok {
			return value
		}
		return newError("identifier not found: %s", node.Value)
	case *ast.FunctionLiteral:
		return &object.Function{Parameters: node.Parameters, Body: node.Body, Env: env}
	case *ast.CallExpression:
		caller := Eval(node.Function, env)
		if isError(caller) {
			return caller
		}

		if fn, ok := caller.(*object.Function); ok {
			if len(node.Arguments) < len(fn.Parameters) {
				return newError("function missing %d parameters", len(fn.Parameters)-len(node.Arguments))
			}

			// TODO: Understand this fucking recursion
			fnEnv := fn.Env.SmartCopy()
			for i, p := range fn.Parameters {
				arg := Eval(node.Arguments[i], fnEnv)
				if isError(arg) {
					return arg
				}
				fnEnv.Set(p.Value, arg)
			}

			ret := Eval(fn.Body, fnEnv)
			if ret, ok := ret.(*object.Return); ok {
				return ret.Value
			}
			return ret
		}

		return newError("expected FUNCTION call got %s call!", caller.Type())
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	}
	return newError("Not implemented eval for %T!", node)
}
