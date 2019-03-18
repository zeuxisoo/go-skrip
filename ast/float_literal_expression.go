package ast

import (
	"github.com/zeuxisoo/go-skrip/token"
)

type FloatLiteralExpression struct {
	Token token.Token
	Value float64
}

func (f *FloatLiteralExpression) expressionNode() {
}

// Implement methods for Node interface
func (f *FloatLiteralExpression) TokenLiteral() string {
	return f.Token.Literal
}

func (f *FloatLiteralExpression) String() string {
	return f.Token.Literal
}
