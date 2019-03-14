package ast

import (
	"github.com/zeuxisoo/go-skriplang/token"
)

type BooleanExpression struct {
	Token token.Token
	Value bool
}

func (b *BooleanExpression) expressionNode() {
}

// Implement methods for Node interface
func (b *BooleanExpression) TokenLiteral() string {
	return b.Token.Literal
}

func (b *BooleanExpression) String() string {
	return b.Token.Literal
}
