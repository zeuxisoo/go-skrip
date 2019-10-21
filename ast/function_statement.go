package ast

import (
	"bytes"
	"strings"

	"github.com/zeuxisoo/go-skrip/token"
)

type FunctionStatement struct {
	Token    token.Token
	Name     *IdentifierExpression
	Function *FunctionLiteralExpression
}

func (f *FunctionStatement) statementNode() {
}

// Implement methods for Node interface
func (f *FunctionStatement) TokenLiteral() string {
	return f.Token.Literal
}

func (f *FunctionStatement) String() string {
	var out bytes.Buffer

	parameters := []string{}
	for _, parameter := range f.Function.Parameters {
		parameters = append(parameters, parameter.String())
	}

	out.WriteString("func ")                        // func
	out.WriteString(f.Name.String())                // name
	out.WriteString("(")                            // (
	out.WriteString(strings.Join(parameters, ", ")) // 	param1, param2, etc
	out.WriteString(")")                            // )
	out.WriteString(" { ")                          // {
	out.WriteString(f.Function.Block.String())      // 	block
	out.WriteString(" }")                           // }

	return out.String()
}
