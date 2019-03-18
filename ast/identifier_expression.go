package ast

import (
	"github.com/zeuxisoo/go-skrip/token"
)

type IdentifierExpression struct {
	Token token.Token
	Value string
}

func (i *IdentifierExpression) expressionNode() {
}

// Implement methods for Node interface
func (i *IdentifierExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IdentifierExpression) String() string {
	return i.Value
}
