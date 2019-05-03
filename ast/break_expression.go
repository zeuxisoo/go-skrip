package ast

import (
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
    return b.Token.Literal
}
