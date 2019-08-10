package builtins

import (
	"github.com/zeuxisoo/go-skrip/object"
)

var (
	NIL = &object.Nil{}
)

// BuiltIns function list
var BuiltIns = map[string]*object.BuiltIn{
	"print":   &object.BuiltIn{Function: Print},
	"println": &object.BuiltIn{Function: Println},
}
