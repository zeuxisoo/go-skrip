package ast

import (
    "github.com/zeuxisoo/go-skrip/token"
)

type ContinueExpression struct {
    Token token.Token
}

func (c *ContinueExpression) expressionNode() {
}

// Implement methods for Node interface
func (c *ContinueExpression) TokenLiteral() string {
    return c.Token.Literal
}

func (c *ContinueExpression) String() string {
    return c.Token.Literal
}
