package ast

import (
	"github.com/zeuxisoo/go-skriplang/token"
)

type Identifier struct {
	Token token.Token
	Value string
}

func (i Identifier) expressionNode() {
}

// Implement methods for Node interface
func (i Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i Identifier) String() string {
	return i.Value
}
