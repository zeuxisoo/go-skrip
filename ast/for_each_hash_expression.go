package ast

import (
	"bytes"

	"github.com/zeuxisoo/go-skrip/token"
)

type ForEachHashExpression struct {
	Token 		token.Token
	Key   		string
	Value 		string
	Data	 	Expression
	Block 		*BlockStatement
}

func (f *ForEachHashExpression) expressionNode() {
}

// Implement methods for Node interface
func (f *ForEachHashExpression) TokenLiteral() string {
	return f.Token.Literal
}

func (f *ForEachHashExpression) String() string {
	var out bytes.Buffer

	out.WriteString("for")									// for
	out.WriteString(" " + f.Key + ", " + f.Value + " ")		// 	key, value
	out.WriteString("in")									// in
	out.WriteString(" " + f.Data.String())					// 	{k: v, k: v}
	out.WriteString(" { ")									// {
	out.WriteString(f.Block.String())						// 	...
	out.WriteString(" } ")									// }

	return out.String()
}
