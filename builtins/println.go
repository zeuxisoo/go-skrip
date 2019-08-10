package builtins

import (
	"fmt"

	"github.com/zeuxisoo/go-skrip/object"
)

// Println function: println(arg1, arg2, ...)
func Println(env *object.Environment, arguments ...object.Object) object.Object {
	if len(arguments) == 0 {
		fmt.Println("")

		return NIL
	}

	parameters := make([]interface{}, len(arguments))

	for index, argument := range arguments {
		parameters[index] = argument.Inspect()
	}

	fmt.Println(parameters...)

	return NIL
}
