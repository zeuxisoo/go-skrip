package ast

import (
	"bytes"

	"github.com/zeuxisoo/go-skrip/token"
)

type IfExpression struct {
	Token 		token.Token
	Condition	Expression
	Block		*BlockStatement
	Alternative *BlockStatement
}

func (i *IfExpression) expressionNode() {
}

// Implement methods for Node interface
func (i *IfExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if ")
	out.WriteString(i.Condition.String())
	out.WriteString(" { ")
	out.WriteString(i.Block.String())
	out.WriteString(" } ")

	if i.Alternative != nil {
		out.WriteString("else")
		out.WriteString(" { ")
		out.WriteString(i.Alternative.String())
		out.WriteString(" } ")
	}

	return out.String()
}
