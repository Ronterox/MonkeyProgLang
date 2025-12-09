package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/token"
	"os"
	"strings"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func buildArray(node *ast.ArrayLiteral, env *object.Environment) object.Object {
	elems := []object.Object{}
	for _, e := range node.Elements {
		val := Eval(e, env)
		if isError(val) {
			return val
		}
		elems = append(elems, val)
	}
	return &object.Array{Elements: elems}
}

func buildHash(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := map[object.HashKey]object.HashPair{}
	for k, v := range node.Pairs {
		key := Eval(k, env)
		if isError(key) {
			return key
		}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("%T=(%v) not yet implemented as hash key!", key, key)
		}

		val := Eval(v, env)
		if isError(val) {
			return val
		}

		pairs[hashKey.HashKey()] = object.HashPair{Key: key, Value: val}
	}
	return &object.Hash{Pairs: pairs}
}

func buildIndex(node *ast.IndexExpression, env *object.Environment) object.Object {
	left := Eval(node.Left, env)
	if isError(left) {
		return left
	}

	if arr, ok := left.(*object.Array); ok {
		right := Eval(node.Index, env)
		if isError(right) {
			return right
		}
		if index, ok := right.(*object.Integer); ok {
			idx := index.Value
			if idx < 0 || idx >= int64(len(arr.Elements)) {
				return NULL
			}
			return arr.Elements[idx]
		}
		return newError("indexing by %s is not yet supported", right.Type())
	} else if hash, ok := left.(*object.Hash); ok {
		right := Eval(node.Index, env)
		if isError(right) {
			return right
		}
		if index, ok := right.(object.Hashable); ok {
			if val, ok := hash.Pairs[index.HashKey()]; ok {
				return val.Value
			}
			return NULL
		}
		return newError("indexing by %s is not yet supported", right.Type())
	}

	return newError("indexing not supported for %s yet", left.Type())

}

func buildCall(node *ast.CallExpression, env *object.Environment) object.Object {
	caller := Eval(node.Function, env)
	if isError(caller) {
		return caller
	}

	if fn, ok := caller.(*object.Builtin); ok {
		args := []object.Object{}
		for _, a := range node.Arguments {
			arg := Eval(a, env)
			if isError(arg) {
				return arg
			}
			args = append(args, arg)
		}
		return fn.Fn(args...)
	}

	if fn, ok := caller.(*object.Function); ok {
		if len(node.Arguments) < len(fn.Parameters) {
			return newError("function missing %d parameters", len(fn.Parameters)-len(node.Arguments))
		}

		// TODO: Understand this fucking recursion
		fnEnv := fn.Env.SmartCopy()
		for i, p := range fn.Parameters {
			arg := Eval(node.Arguments[i], env)
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
}

func buildIf(node *ast.IfExpression, env *object.Environment) object.Object {
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
}

func buildInfix(node *ast.InfixExpression, env *object.Environment) object.Object {
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
		case token.LE:
			if left.Value <= right.Value {
				return TRUE
			}
			return FALSE
		case token.GE:
			if left.Value >= right.Value {
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
	} else if left.Type() == object.STRING && right.Type() == object.STRING {
		left := left.(*object.String)
		right := right.(*object.String)

		switch node.Operator {
		case token.PLUS:
			return &object.String{Value: left.Value + right.Value}
		case token.MINUS:
			return &object.String{Value: strings.ReplaceAll(left.Value, right.Value, "")}
		}
	} else if left.Type() == object.INTEGER && right.Type() == object.STRING {
		left := left.(*object.Integer)
		right := right.(*object.String)

		switch node.Operator {
		case token.PLUS:
			return &object.String{Value: fmt.Sprintf("%d%s", left.Value, right.Value)}
		case token.ASTERISK:
			return &object.String{Value: strings.Repeat(right.Value, int(left.Value))}
		}
	} else if left.Type() == object.STRING && right.Type() == object.INTEGER {
		left := left.(*object.String)
		right := right.(*object.Integer)

		switch node.Operator {
		case token.PLUS:
			return &object.String{Value: fmt.Sprintf("%s%d", left.Value, right.Value)}
		case token.ASTERISK:
			return &object.String{Value: strings.Repeat(left.Value, int(right.Value))}
		}
	}

	return newError("Operation %s between %s and %s not implemented!", node.Operator, left.Type(), right.Type())
}

func buildPrefix(node *ast.PrefixExpression, env *object.Environment) object.Object {
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
}

func buildBuiltin(node *ast.Identifier, env *object.Environment) object.Object {
	switch node.Value {
	case "null":
		return NULL
	case "len":
		return &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				switch obj := args[0].(type) {
				case *object.String:
					return &object.Integer{Value: int64(len(obj.Value))}
				case *object.Array:
					return &object.Integer{Value: int64(len(obj.Elements))}
				}
				return newError("argument to `len` not supported, got %s", args[0].Type())
			},
		}
	case "head":
		return &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				switch e := args[0].(type) {
				case *object.Array:
					if len(e.Elements) == 0 {
						return NULL
					}
					return e.Elements[0]
				case *object.String:
					if len(e.Value) == 0 {
						return NULL
					}
					return &object.String{Value: string(e.Value[0])}
				}
				return newError("head is not implemented for %s", args[0].Type())
			},
		}
	case "last":
		return &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				switch e := args[0].(type) {
				case *object.Array:
					if length := len(e.Elements); length > 0 {
						return e.Elements[length-1]
					}
					return NULL
				case *object.String:
					if len(e.Value) == 0 {
						return NULL
					}
					return &object.String{Value: string(e.Value[len(e.Value)-1])}
				}
				return newError("last is not implemented for %s", args[0].Type())
			},
		}
	case "tail":
		return &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}
				switch e := args[0].(type) {
				case *object.Array:
					newArr := &object.Array{}
					if length := len(e.Elements); length != 0 {
						newArr.Elements = e.Elements[1:length]
					}
					return newArr
				case *object.String:
					if len(e.Value) == 0 {
						return NULL
					}
					return &object.String{Value: string(e.Value[1:])}
				}
				return newError("tail is not implemented for %s", args[0].Type())
			},
		}
	case "push":
		return &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return newError("wrong number of arguments. got=%d, want=2", len(args))
				}
				switch arr := args[0].(type) {
				case *object.Array:
					return &object.Array{Elements: append(arr.Elements, args[1])}
				}
				return newError("push is not implemented for %s", args[0].Type())
			},
		}
	case "string":
		return &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				all := []string{}
				for _, arg := range args {
					all = append(all, arg.Inspect())
				}
				return &object.String{Value: strings.Join(all, "")}
			},
		}
	case "echo":
		return &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				all := []string{}
				for _, arg := range args {
					all = append(all, arg.Inspect())
				}
				fmt.Println(strings.Join(all, ""))
				return NULL
			},
		}
	case "read":
		return &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				switch file := args[0].(type) {
				case *object.String:
					content, err := os.ReadFile(file.Value)
					if err != nil {
						return newError("could not read file %s", file.Value)
					}
					return &object.String{Value: string(content)}
				}
				return newError("argument to `read` not supported yet, got %s", args[0].Type())
			},
		}
	case "eval":
		return &object.Builtin{
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return newError("wrong number of arguments. got=%d, want=1", len(args))
				}

				if text, ok := args[0].(*object.String); ok {
					lexer := lexer.New(text.Value)
					parser := parser.New(lexer)
					return Eval(parser.ParseProgram(), env)
				}
				return newError("argument to `eval` not supported yet, got %s", args[0].Type())
			},
		}
	}
	return nil
}

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
		return buildPrefix(node, env)
	case *ast.InfixExpression:
		return buildInfix(node, env)
	case *ast.IfExpression:
		return buildIf(node, env)
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
		if obj := buildBuiltin(node, env); obj != nil {
			return obj
		}
		return newError("identifier not found: %s", node.Value)
	case *ast.FunctionLiteral:
		return &object.Function{Parameters: node.Parameters, Body: node.Body, Env: env}
	case *ast.CallExpression:
		return buildCall(node, env)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.ArrayLiteral:
		return buildArray(node, env)
	case *ast.IndexExpression:
		return buildIndex(node, env)
	case *ast.HashLiteral:
		return buildHash(node, env)
	}
	return newError("Not implemented eval for %T!", node)
}
