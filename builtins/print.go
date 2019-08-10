package builtins

import (
	"fmt"

	"github.com/zeuxisoo/go-skrip/object"
)

// Print function: print(arg1, arg2, ...)
func Print(env *object.Environment, arguments ...object.Object) object.Object {
	if len(arguments) == 0 {
		fmt.Print("")

		return NIL
	}

	parameters := make([]interface{}, len(arguments))

	for index, argument := range arguments {
		parameters[index] = argument.Inspect()
	}

	fmt.Print(parameters...)

	return NIL
}
