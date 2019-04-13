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
	case *ast.LetStatement:
		return evalLetStatement(node, env)
	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)
	case *ast.FunctionStatement:
		return evalFunctionStatement(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
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
	case *ast.CallExpression:
		return evalCallExpression(node, env)
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

func evalLetStatement(let *ast.LetStatement, env *object.Environment) object.Object {
	obj := Eval(let.Value, env)
	if isError(obj) == true {
		return obj
	}

	env.Set(let.Name.Value, obj)

	return obj
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

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var obj object.Object

	for _, statement := range block.Statements {
		obj := Eval(statement, env)
		if obj != nil {
			objectType := obj.Type()

			if objectType == object.RETURN_VALUE_OBJECT || objectType == object.ERROR_OBJECT {
				return obj
			}
		}
	}

	return obj
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

func evalCallExpression(call *ast.CallExpression, env *object.Environment) object.Object {
	// E.g. myFunction(argument1, argument2, ...)

	// Evaluate call myFunction
	function := Eval(call.Function, env)
	if isError(function) == true {
		return function
	}

	// Evaluate call argument1, argument2
	arguments := evalExpressions(call.Arguments, env)
	if len(arguments) == 1 && isError(arguments[0]) == true {
		return arguments[0]
	}

	// Apply to call arguments to function
	result := applyFunction(env, function, arguments)
	if isError(result) == true {
		return newError("Error calling %s: %s", call.Function, result.Inspect())
	}

	return result
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

func applyFunction(env *object.Environment, function object.Object, arguments []object.Object) object.Object {
	switch fn := function.(type) {
	// custom function
	case *object.Function:
		extendEnvironment, err := extendFunctionEnvironment(fn, arguments)
		if err != nil {
			return err
		}

		evaluated := Eval(fn.Block, extendEnvironment)

		return unwrapReturnValue(evaluated)
	// built-in function
	case *object.BuiltIn:
		return fn.Function(env, arguments...)
	default:
		return newError("%s is not a function", fn.Type())
	}
}

func extendFunctionEnvironment(function *object.Function, arguments []object.Object) (*object.Environment, *object.Error) {
	// Create scoped environment for current function
	environment := object.NewEnclosedEnvironment(function.Environment)

	if len(arguments) != len(function.Parameters) {
		return nil, newError(
			"not enough arguments for %s function, Got: %s, Expected: %s",
			function.Inspect(), arguments, function.Parameters,
		)
	}

	// Setup variable by parameter is the name, arguments is the value
	for index, parameter := range function.Parameters {
		environment.Set(parameter.Value, arguments[index])
	}

	return environment, nil
}

func unwrapReturnValue(obj object.Object) object.Object {
	// Return value only if current object is return value object
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue
	}

	return obj
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
