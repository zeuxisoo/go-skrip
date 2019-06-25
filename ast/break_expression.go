package ast

import (
	"bytes"

	"github.com/zeuxisoo/go-skrip/token"
)

type BreakExpression struct {
	Token token.Token
}

func (b *BreakExpression) expressionNode() {
}

// Implement methods for Node interface
func (b *BreakExpression) TokenLiteral() string {
	return b.Token.Literal
}

func (b *BreakExpression) String() string {
	var out bytes.Buffer

	out.WriteString(b.Token.Literal)
	out.WriteString("; ")

	return out.String()
}
