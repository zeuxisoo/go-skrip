package evaluator

import (
	"github.com/zeuxisoo/go-skrip/ast"
	"github.com/zeuxisoo/go-skrip/object"
)

var (
	NIL = &object.Nil{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	// Statements
	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)
	// Expressions
	case *ast.IntegerLiteralExpression:
		return evalIntegerLiteralExpression(node, env)
	}

	return NIL
}

// Eval function
func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		}
	}

	return result
}

func evalReturnStatement(ret *ast.ReturnStatement, env *object.Environment) object.Object {
	obj := Eval(ret.ReturnValue, env)

	if isError(obj) {
		return obj
	}

	return &object.ReturnValue{
		Value: obj,
	}
}

func evalIntegerLiteralExpression(integer *ast.IntegerLiteralExpression, env *object.Environment) object.Object {
	return &object.Integer{
		Value: integer.Value,
	}
}

// Helper functions
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJECT
	}

	return false
}
