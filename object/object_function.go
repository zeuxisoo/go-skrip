package object

import (
	"bytes"
	"strings"

	"github.com/zeuxisoo/go-skrip/ast"
)

type Function struct {
	Parameters  []*ast.IdentifierExpression
	Block       *ast.BlockStatement
	Environment *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJECT
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	parameters := []string{}
	for _, parameter := range f.Parameters {
		parameters = append(parameters, parameter.String())
	}

	out.WriteString("func")
	out.WriteString("(")
	out.WriteString(strings.Join(parameters, ", "))
	out.WriteString(")")
	out.WriteString(" { ")
	out.WriteString(f.Block.String())
	out.WriteString(" } ")

	return out.String()
}
