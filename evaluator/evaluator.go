package evaluator

import (
	"fmt"
	"strconv"

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
	case *ast.RangeExpression:
		return evalRangeExpression(node, env)
	case *ast.CallExpression:
		return evalCallExpression(node, env)
	case *ast.IndexExpression:
		return evalIndexExpression(node, env)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, env)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) == true {
			return left
		}

		right := Eval(node.Right, env)
		if isError(right) == true {
			return right
		}

		return evalInfixExpression(left, node.Operator, right, env)
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

	// Set function name to environment like let a = func() {}, a will be variable
	env.Set(function.Name.Value, obj)

	// Cast object.Object to object.Function
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

func evalRangeExpression(rng *ast.RangeExpression, env *object.Environment) object.Object {
	start := Eval(rng.Start, env)
	if isError(start) == true {
		return start
	}

	end := Eval(rng.End, env)
	if isError(end) == true {
		return end
	}

	switch {
	// int..int
	case start.Type() == object.INTEGER_OBJECT && end.Type() == object.INTEGER_OBJECT:
		return evalRangeIntegerExpression(start, end)
	// float..float
	case start.Type() == object.FLOAT_OBJECT && end.Type() == object.FLOAT_OBJECT:
		return evalRangeFloatExpression(start, end)
	// string..string
	case start.Type() == object.STRING_OBJECT && end.Type() == object.STRING_OBJECT:
		return evalRangeStringExpression(start, end)
	default:
		return newError(
			"Range operator not support for %s (%s) to %s (%s)",
			start.Inspect(), start.Type(), end.Inspect(), end.Type(),
		)
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

func evalIndexExpression(index *ast.IndexExpression, env *object.Environment) object.Object {
	left := Eval(index.Left, env)
	if isError(left) == true {
		return left
	}

	idx := Eval(index.Index, env)
	if isError(idx) == true {
		return idx
	}

	switch {
	// array[integer]
	case left.Type() == object.ARRAY_OBJECT && idx.Type() == object.INTEGER_OBJECT:
		return evalArrayIndexExpression(left, idx)
	// hash[hashable]
	case left.Type() == object.HASH_OBJECT:
		return evalHashIndexExpression(left, idx)
	// string[integer]
	case left.Type() == object.STRING_OBJECT && idx.Type() == object.INTEGER_OBJECT:
		return evalStringIndexExpression(left, idx)
	default:
		return newError("Index operator not support for %s on %s", idx.Inspect(), left.Type())
	}
}

func evalPrefixExpression(prefix *ast.PrefixExpression, env *object.Environment) object.Object {
	right := Eval(prefix.Right, env)
	if isError(right) == true {
		return right
	}

	switch prefix.Operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	case "+":
		return evalPlusPrefixOperatorExpression(right)
	default:
		return newError("Unknown operator %s with %s", prefix.Operator, right.Type())
	}
}

func evalInfixExpression(left object.Object, operator string, right object.Object, env *object.Environment)  object.Object {
	switch {
	// and
	case operator == "&&":
		return nativeBoolToBooleanObject(objectToNativeBoolean(left) && objectToNativeBoolean(right))
	// or
	case operator == "||":
		return nativeBoolToBooleanObject(objectToNativeBoolean(left) || objectToNativeBoolean(right))
	// int operator int
	case left.Type() == object.INTEGER_OBJECT && right.Type() == object.INTEGER_OBJECT:
		return evalIntegerIntegerInfixExpression(left, operator, right)
	// int operator float
	case left.Type() == object.INTEGER_OBJECT && right.Type() == object.FLOAT_OBJECT:
		return evalIntegerFloatInfixExpression(left, operator, right)
	// float operator float
	case left.Type() == object.FLOAT_OBJECT && right.Type() == object.FLOAT_OBJECT:
		return evalFloatFloatInfixExpression(left, operator, right)
	// float operator int
	case left.Type() == object.FLOAT_OBJECT && right.Type() == object.INTEGER_OBJECT:
		return evalFloatIntegerInfixExpression(left, operator, right)
	// string operator string
	case left.Type() == object.STRING_OBJECT && right.Type() == object.STRING_OBJECT:
		return evalStringStringInfixExpression(left, operator, right)
	// array operator array
	case left.Type() == object.ARRAY_OBJECT && right.Type() == object.ARRAY_OBJECT:
		return evalArrayArrayInfixExpression(left, operator, right, env)
	// hash operator hash
	case left.Type() == object.HASH_OBJECT && right.Type() == object.HASH_OBJECT:
		return evalHashHashInfixExpression(left, operator, right)
	// TODO: equals when left data type and right data type are different
	case operator == "==":
		return nil
	// TODO: not equals
	case operator == "!=":
		return nil
	case left.Type() != right.Type():
		return newError("Type mismatch %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("Unknown operator %s %s %s", left.Type(), operator, right.Type())
	}
}

// For boolean expression
func nativeBoolToBooleanObject(value bool) object.Object {
	if value == true {
		return TRUE
	}

	return FALSE
}

// For call expression
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

// For range expression
func evalRangeIntegerExpression(start object.Object, end object.Object) object.Object {
	startObject := start.(*object.Integer)
	endObject   := end.(*object.Integer)

	elements := make([]object.Object, 0)
	for i := startObject.Value; i < endObject.Value; i++ {
		elements = append(elements, &object.Integer{
			Value: i,
		})
	}

	return &object.Array{
		Elements: elements,
	}
}

func evalRangeFloatExpression(start object.Object, end object.Object) object.Object {
	startObject := start.(*object.Float)
	endObject   := end.(*object.Float)

	elements := make([]object.Object, 0)
	for i := startObject.Value; i < endObject.Value; i += 0.1 {
		elements = append(elements, &object.Float{
			Value: i,
		})
	}

	return &object.Array{
		Elements: elements,
	}
}

func evalRangeStringExpression(start object.Object, end object.Object) object.Object {
	startObject := start.(*object.String)
	endObject   := end.(*object.String)

	if len(startObject.Value) > 1 {
		return newError("Range start value must be char only")
	}

	if len(endObject.Value) > 1 {
		return newError("Range end value must be char only")
	}

	elements := make([]object.Object, 0)

	startByte := int32([]rune(startObject.Value)[0])
	endByte   := int32([]rune(endObject.Value)[0])

	if startByte >= endByte {
		// E.g. z -> a
		for i := startByte; i > endByte; i-- {
			elements = append(elements, &object.String{
				Value: string(string(i)),
			})
		}
	}else{
		// E.g. a -> z
		for i := startByte; i < endByte; i++ {
			elements = append(elements, &object.String{
				Value: string(string(i)),
			})
		}
	}

	return &object.Array{
		Elements: elements,
	}
}

// For index expression
func evalArrayIndexExpression(left object.Object, index object.Object) object.Object {
	// for array[integer]
	arrayObject := left.(*object.Array)
	indexObject := index.(*object.Integer)

	indexValue := indexObject.Value
	maxLength  := int64(len(arrayObject.Elements) - 1)

	if indexValue < 0 || indexValue > maxLength {
		return NIL
	}

	return arrayObject.Elements[indexValue]
}

func evalHashIndexExpression(left object.Object, index object.Object) object.Object {
	// for hash[hashable]
	hashObject := left.(*object.Hash)

	// Make sure the string object is hashable and call HashKey method to get hash key
	key, ok := index.(object.Hashable)
	if ok == false {
		return newError("Cannot use %s as hash key", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if ok == false {
		return NIL
	}

	return pair.Value
}

func evalStringIndexExpression(left object.Object, index object.Object) object.Object {
	// for string[integer]
	stringObject := left.(*object.String)
	indexObject  := index.(*object.Integer)

	stringValue := stringObject.Value
	indexValue  := indexObject.Value

	maxLength   := int64(len(stringValue) - 1)

	if indexValue < 0 || indexValue > maxLength {
		return NIL
	}

	return &object.String{
		Value: string(stringObject.Value[indexValue]),
	}
}

// For prefix expression
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NIL:
		return TRUE
	default:
		switch obj := right.(type) {
		case *object.Integer:
			if obj.Value == 0 {
				return TRUE
			}
			return FALSE
		case *object.Float:
			if obj.Value == 0.0 {
				return TRUE
			}
			return FALSE
		case *object.String:
			if len(obj.Value) == 0 {
				return TRUE
			}
			return FALSE
		case *object.Array:
			if len(obj.Elements) == 0 {
				return TRUE
			}
			return FALSE
		case *object.Hash:
			if len(obj.Pairs) == 0 {
				return TRUE
			}
			return FALSE
		default:
			return FALSE
		}
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	switch obj := right.(type) {
	case *object.Integer:
		return &object.Integer{
			Value: -obj.Value,
		}
	case *object.Float:
		return &object.Float{
			Value: -obj.Value,
		}
	default:
		return newError("Unnown operator - with %s", right.Type())
	}
}

func evalPlusPrefixOperatorExpression(right object.Object) object.Object {
	return right
}

// For infix expression
func evalIntegerIntegerInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	leftInteger  := left.(*object.Integer)
	rightInteger := right.(*object.Integer)

	leftValue := leftInteger.Value
	rightValue := rightInteger.Value

	switch operator {
	case "+":
		return &object.Integer{ Value: leftValue + rightValue }
	case "-":
		return &object.Integer{ Value: leftValue - rightValue }
	case "*":
		return &object.Integer{ Value: leftValue * rightValue }
	case "/":
		return &object.Integer{ Value: leftValue / rightValue }
	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case "<=":
		return nativeBoolToBooleanObject(leftValue <= rightValue)
	case ">=":
		return nativeBoolToBooleanObject(leftValue >= rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return newError("Unknown operator %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerFloatInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	leftInteger := left.(*object.Integer)
	rightFloat  := right.(*object.Float)

	leftValue  := float64(leftInteger.Value)
	rightValue := rightFloat.Value

	switch operator {
	case "+":
		return &object.Float{ Value: humanFloat(leftValue + rightValue) }
	case "-":
		return &object.Float{ Value: humanFloat(leftValue - rightValue) }
	case "*":
		return &object.Float{ Value: humanFloat(leftValue * rightValue) }
	case "/":
		return &object.Float{ Value: humanFloat(leftValue / rightValue) }
	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case "<=":
		return nativeBoolToBooleanObject(leftValue <= rightValue)
	case ">=":
		return nativeBoolToBooleanObject(leftValue >= rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return newError("Unknown operator %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatFloatInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	leftFloat := left.(*object.Float)
	rightFloat  := right.(*object.Float)

	leftValue  := leftFloat.Value
	rightValue := rightFloat.Value

	switch operator {
	case "+":
		return &object.Float{ Value: humanFloat(leftValue + rightValue) }
	case "-":
		return &object.Float{ Value: humanFloat(leftValue - rightValue) }
	case "*":
		return &object.Float{ Value: humanFloat(leftValue * rightValue) }
	case "/":
		return &object.Float{ Value: humanFloat(leftValue / rightValue) }
	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case "<=":
		return nativeBoolToBooleanObject(leftValue <= rightValue)
	case ">=":
		return nativeBoolToBooleanObject(leftValue >= rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return newError("Unknown operator %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalFloatIntegerInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	leftFloat    := left.(*object.Float)
	rightInteger := right.(*object.Integer)

	leftValue  := leftFloat.Value
	rightValue := float64(rightInteger.Value)

	switch operator {
	case "+":
		return &object.Float{ Value: humanFloat(leftValue + rightValue) }
	case "-":
		return &object.Float{ Value: humanFloat(leftValue - rightValue) }
	case "*":
		return &object.Float{ Value: humanFloat(leftValue * rightValue) }
	case "/":
		return &object.Float{ Value: humanFloat(leftValue / rightValue) }
	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case "<=":
		return nativeBoolToBooleanObject(leftValue <= rightValue)
	case ">=":
		return nativeBoolToBooleanObject(leftValue >= rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return newError("Unknown operator %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalStringStringInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	leftString  := left.(*object.String)
	rightString := right.(*object.String)

	leftValue  := leftString.Value
	rightValue := rightString.Value

	switch operator {
	case "+":
		return &object.String{ Value: leftValue + rightValue }
	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case "<=":
		return nativeBoolToBooleanObject(leftValue <= rightValue)
	case ">=":
		return nativeBoolToBooleanObject(leftValue >= rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return newError("Unknown operator %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalArrayArrayInfixExpression(left object.Object, operator string, right object.Object, env *object.Environment) object.Object {
	leftArray  := left.(*object.Array)
	rightArray := right.(*object.Array)

	leftElements  := leftArray.Elements
	rightElements := rightArray.Elements

	switch operator {
	case "+":
		return &object.Array{ Elements: append(leftElements, rightElements...) }
	case "==":
		if len(leftElements) != len(rightElements) {
			return FALSE
		}

		for i := range leftElements {
			compareResult := evalInfixExpression(leftElements[i], "==", rightElements[i], env)

			if objectToNativeBoolean(compareResult) != true {
				return FALSE
			}
		}

		return TRUE
	case "!=":
		if evalArrayArrayInfixExpression(left, "==", right, env) == TRUE {
			return FALSE
		}

		return TRUE
	default:
		return newError("Unknown operator %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalHashHashInfixExpression(left object.Object, operator string, right object.Object) object.Object {
	leftHash  := left.(*object.Hash)
	rightHash := right.(*object.Hash)

	leftPairs  := leftHash.Pairs
	rightPairs := rightHash.Pairs

	switch operator {
	case "+":
		for _, hashKey := range rightHash.Order {
			pair, _ := rightPairs[hashKey]

			if hashable, ok := pair.Key.(object.Hashable); ok {
				hashKey := hashable.HashKey()

				// Add the hash key into hash order if the hash key is not exists in hash order
				_, exists := leftPairs[hashKey]
				if exists == false {
					leftHash.Order = append(leftHash.Order, hashKey)
				}

				leftPairs[hashKey] = object.HashPair{
					Key  : pair.Key,
					Value: pair.Value,
				}
			}else{
				return newError("Cannot use %s as hash key", pair.Key.Type())
			}
		}

		return &object.Hash{
			Order: leftHash.Order,
			Pairs: leftPairs,
		}
	case "==":
		if len(leftPairs) != len(rightPairs) {
			return FALSE
		}

		matchCount := 0
		for leftPairKey, leftPair := range leftPairs {
			for rightPairKey, rightPair := range rightPairs {
				if leftPairKey.Value == rightPairKey.Value &&
					leftPair.Key.Inspect() == rightPair.Key.Inspect() &&
					leftPair.Value.Inspect() == rightPair.Value.Inspect() {
						matchCount = matchCount + 1
					}
			}
		}

		if matchCount == len(leftPairs) {
			return TRUE
		}

		return FALSE
	case "!=":
		if evalHashHashInfixExpression(left, "==", right) == TRUE {
			return FALSE
		}

		return TRUE
	default:
		return newError("Unknown operator %s %s %s", left.Type(), operator, right.Type())
	}
}

// Helper functions
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

func objectToNativeBoolean(obj object.Object) bool {
	if ret, ok := obj.(*object.ReturnValue); ok {
		obj = ret.Value
	}

	switch o := obj.(type) {
	case *object.Boolean:
		return o.Value
	case *object.Nil:
		return false
	case *object.String:
		return len(o.Value) != 0
	case *object.Integer:
		if o.Value == 0 {
			return false
		}
		return true
	case *object.Float:
		if o.Value == 0.0 {
			return false
		}
		return true
	case *object.Array:
		if len(o.Elements) == 0 {
			return false
		}
		return true
	case *object.Hash:
		if len(o.Pairs) == 0 {
			return false
		}
		return true
	default:
		return true
	}
}

// Try to make the floating point more humanized
// E.g.
// In standard.
// - 3 * 2.3 will be 6.8999999999999995
// - 1 - 2.3 will be -1.2999999999999998
// In human.
// - 3 * 2.3 will be 6.9
// - 1 - 2.3 will be -1.3
// More information
// - https://stackoverflow.com/questions/588004/is-floating-point-math-broken
func humanFloat(value float64) float64 {
	humanFloat := fmt.Sprintf("%f", value)
	realFloat, _ := strconv.ParseFloat(humanFloat, 64)

	return realFloat
}

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
