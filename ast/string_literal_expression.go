package ast

import (
	"github.com/zeuxisoo/go-skrip/token"
)

type StringLiteralExpression struct {
	Token token.Token
	Value string
}

func (s *StringLiteralExpression) expressionNode() {
}

// Implement methods for Node interface
func (s *StringLiteralExpression) TokenLiteral() string {
	return s.Token.Literal
}

func (s *StringLiteralExpression) String() string {
	return s.Token.Literal
}
