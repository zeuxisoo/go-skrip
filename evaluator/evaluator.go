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
	case *ast.FunctionStatement:
		return evalFunctionStatement(node, env)
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
	case *ast.ArrayLiteralExpression:
		return evalArrayLiteralExpression(node, env)
	case *ast.HashLiteralExpression:
		return evalHashLiteralExpression(node, env)
	case *ast.FunctionLiteralExpression:
		return evalFunctionLiteralExpression(node, env)
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

func evalFunctionStatement(function *ast.FunctionStatement, env *object.Environment) object.Object {
	obj := evalFunctionLiteralExpression(function.Function, env)

	functionObject := obj.(*object.Function)

	return functionObject
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

func evalArrayLiteralExpression(array *ast.ArrayLiteralExpression, env *object.Environment) object.Object {
	elements := evalExpressions(array.Elements, env)

	if len(elements) == 1 && isError(elements[0]) == true {
		return elements[0]
	}

	return &object.Array{
		Elements: elements,
	}
}

func evalHashLiteralExpression(hash *ast.HashLiteralExpression, env *object.Environment) object.Object {
	hashObject := &object.Hash{
		Order: []object.HashKey{},
		Pairs: make(map[object.HashKey]object.HashPair),
	}

	// Loop by hash order (keys) expressions
	for _, orderKey := range hash.Order {
		// Get key object from evaluated hash order (key) expression
		key := Eval(orderKey, env)
		if isError(key) == true {
			return key
		}

		// Ensure the key object must be hashable
		hashableKey, ok := key.(object.Hashable)
		if ok == false {
			return newError("Cannot use %s as hash key", key.Type())
		}

		// Get hash pair value from hash literal expression pairs by hash order (key) expression
		pairValue, _ := hash.Pairs[orderKey]

		// Get value object from evaluated hash pair value expression
		value := Eval(pairValue, env)
		if isError(value) == true {
			return value
		}

		// If hashable key object is not exists in hash object order slice, add it (for ordered iterable)
		hashedKey := hashableKey.HashKey()

		_, exists := hashObject.Pairs[hashedKey]
		if exists == false {
			hashObject.Order = append(hashObject.Order, hashedKey)
		}

		// Create pair and add to hash object pairs
		hashObject.Pairs[hashedKey] = object.HashPair{
			Key  : key,
			Value: value,
		}
	}

	return hashObject
}

func evalFunctionLiteralExpression(function *ast.FunctionLiteralExpression, env *object.Environment) object.Object {
	return &object.Function{
		Parameters : function.Parameters,
		Block      : function.Block,
		Environment: env,
	}
}

//
func nativeBoolToBooleanObject(value bool) object.Object {
	if value == true {
		return TRUE
	}

	return FALSE
}

func evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var objects []object.Object

	for _, expression := range expressions {
		evaluated := Eval(expression, env)
		if isError(evaluated) == true {
			return []object.Object{ evaluated }
		}

		objects = append(objects, evaluated)
	}

	return objects
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
