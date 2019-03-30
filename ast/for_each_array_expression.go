package ast

import (
	"bytes"

	"github.com/zeuxisoo/go-skrip/token"
)

type ForEachArrayExpression struct {
	Token 		token.Token
	Value 		string
	Data	 	Expression
	Block 		*BlockStatement
}

func (f *ForEachArrayExpression) expressionNode() {
}

// Implement methods for Node interface
func (f *ForEachArrayExpression) TokenLiteral() string {
	return f.Token.Literal
}

func (f *ForEachArrayExpression) String() string {
	var out bytes.Buffer

	out.WriteString("for")					// for
	out.WriteString(" " + f.Value + " ")	//	value
	out.WriteString("in")					// in
	out.WriteString(" " + f.Data.String())	// 	{k: v, k: v}
	out.WriteString(" { ")					// {
	out.WriteString(f.Block.String())		// 	...
	out.WriteString(" } ")					// }

	return out.String()
}
