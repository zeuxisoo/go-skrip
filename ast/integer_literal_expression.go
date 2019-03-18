package ast

import (
	"github.com/zeuxisoo/go-skrip/token"
)

type IntegerLiteralExpression struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteralExpression) expressionNode() {
}

// Implement methods for Node interface
func (i *IntegerLiteralExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IntegerLiteralExpression) String() string {
	return i.Token.Literal
}
