package ast

import (
	"bytes"
	"strings"

	"github.com/zeuxisoo/go-skriplang/token"
)

type FunctionLiteralExpression struct {
	Token 		token.Token
	Parameters 	[]*IdentifierExpression
	Block		*BlockStatement
}

func (f *FunctionLiteralExpression) expressionNode() {
}

// Implement methods for Node interface
func (f *FunctionLiteralExpression) TokenLiteral() string {
	return f.Token.Literal
}

func (f *FunctionLiteralExpression) String() string {
	var out bytes.Buffer

	parameters := []string{}

	for _, parameter := range f.Parameters {
		parameters = append(parameters, parameter.String())
	}

	out.WriteString(f.TokenLiteral())				// functionName
	out.WriteString("(")							// (
	out.WriteString(strings.Join(parameters, ", "))	// parameter1, parameter2, etc
	out.WriteString(") ")							// )
	out.WriteString(f.Block.String()) 				// { ...... } without "{" and "}"

	return out.String()
}