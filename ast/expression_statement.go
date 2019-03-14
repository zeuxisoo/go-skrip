package ast

import (
	"github.com/zeuxisoo/go-skriplang/token"
)

type ExpressionStatement struct {
	Token		token.Token
	Expression 	Expression
}

func (e *ExpressionStatement) statementNode() {
}

// Implement methods for Node interface
func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}

func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}

	return ""
}
