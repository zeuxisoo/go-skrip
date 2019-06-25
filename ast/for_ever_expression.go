package ast

import (
	"bytes"

	"github.com/zeuxisoo/go-skrip/token"
)

type ForEverExpression struct {
	Token token.Token
	Block *BlockStatement
}

func (f *ForEverExpression) expressionNode() {
}

// Implement methods for Node interface
func (f *ForEverExpression) TokenLiteral() string {
	return f.Token.Literal
}

func (f *ForEverExpression) String() string {
	var out bytes.Buffer

	out.WriteString("for")
	out.WriteString(" { ")
	out.WriteString(f.Block.String())
	out.WriteString("} ")

	return out.String()
}
