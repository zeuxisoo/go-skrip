package ast

import (
    "github.com/zeuxisoo/go-skrip/token"
)

type NilLiteralExpression struct {
    Token token.Token
}

func (n *NilLiteralExpression) expressionNode() {
}

// Implement methods for Node interface
func (n *NilLiteralExpression) TokenLiteral() string {
    return n.Token.Literal
}

func (n *NilLiteralExpression) String() string {
    return n.Token.Literal
}
