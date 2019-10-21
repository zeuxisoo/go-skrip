package ast

import (
	"bytes"
	"strings"

	"github.com/zeuxisoo/go-skrip/token"
)

type FunctionLiteralExpression struct {
	Token      token.Token
	Parameters []*IdentifierExpression
	Block      *BlockStatement
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

	// Only for expression:
	// let foo = func(params) { block }
	out.WriteString(f.TokenLiteral())               // functionName
	out.WriteString("(")                            // (
	out.WriteString(strings.Join(parameters, ", ")) // 	parameter1, parameter2, etc
	out.WriteString(") ")                           // )
	out.WriteString("{ ")                           // {
	out.WriteString(f.Block.String())               // 	block
	out.WriteString(" }")                           // }

	return out.String()
}
