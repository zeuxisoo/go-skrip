package evaluator

import (
	"fmt"

	"github.com/zeuxisoo/go-skrip/ast"
	"github.com/zeuxisoo/go-skrip/object"
)

var (
	NIL   = &object.Nil{}
	TRUE  = &object.Boolean{ Value: true }
	FALSE = &object.Boolean{ Value: false }
)

var builtIns = map[string]*object.BuiltIn{}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	// Statements
	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	// Expressions
	case *ast.IntegerLiteralExpression:
		return evalIntegerLiteralExpression(node, env)
	case *ast.FloatLiteralExpression:
		return evalFloatLiteralExpression(node, env)
	case *ast.StringLiteralExpression:
		return evalStringLiteralExpression(node, env)
	case *ast.IdentifierExpression:
		return evalIdentifierExpression(node, env)
	case *ast.BooleanExpression:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.HashLiteralExpression:
		return evalHashLiteralExpression(node, env)
	}

	return NIL
}

func RegisterBuiltIn(name string, function object.BuiltInFunction) {
	builtIns[name] = &object.BuiltIn{
		Function: function,
	}
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

	if isError(obj) == true {
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

func evalFloatLiteralExpression(float *ast.FloatLiteralExpression, env *object.Environment) object.Object {
	return &object.Float{
		Value: float.Value,
	}
}

func evalStringLiteralExpression(str *ast.StringLiteralExpression, env *object.Environment) object.Object {
	return &object.String{
		Value: str.Value,
	}
}

func evalIdentifierExpression(identifer *ast.IdentifierExpression, env *object.Environment) object.Object {
	if value, ok := env.Get(identifer.Value); ok {
		return value
	}

	if builtIn, ok := builtIns[identifer.Value]; ok {
		return builtIn
	}

	return newError("Identifier not found: " + identifer.Value)
}

func evalHashLiteralExpression(hash *ast.HashLiteralExpression, env *object.Environment) object.Object {
	hashObject := &object.Hash{
		Order: []object.HashKey{},
		Pairs: make(map[object.HashKey]object.HashPair),
	}

	for keyNode, valueNode := range hash.Pairs {
		// Get key
		key := Eval(keyNode, env)
		if isError(key) == true {
			return key
		}

		// Ensure the key must be hashable
		hashableKey, ok := key.(object.Hashable)
		if ok == false {
			return newError("Cannot use %s as hash key", key.Type())
		}

		// Get value
		value := Eval(valueNode, env)
		if isError(value) == true {
			return value
		}

		// If hashable key object is not exists in order slice, add it for iterable
		hashedKey := hashableKey.HashKey()

		_, exists := hashObject.Pairs[hashedKey]
		if exists == false {
			hashObject.Order = append(hashObject.Order, hashedKey)
		}

		// Add pair to hash object
		hashObject.Pairs[hashedKey] = object.HashPair{
			Key  : key,
			Value: value,
		}
	}

	return hashObject
}

//
func nativeBoolToBooleanObject(value bool) object.Object {
	if value == true {
		return TRUE
	}

	return FALSE
}

// Helper functions
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJECT
	}

	return false
}

func newError(format string, values ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, values...),
	}
}
