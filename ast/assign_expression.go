package ast

import (
	"bytes"

	"github.com/zeuxisoo/go-skrip/token"
)

type AssignExpression struct {
	Token token.Token
	Left  Expression
	Value Expression
}

func (a *AssignExpression) expressionNode() {
}

func (a *AssignExpression) TokenLiteral() string {
	return a.Token.Literal
}

func (a *AssignExpression) String() string {
	var out bytes.Buffer

	out.WriteString(a.Left.String())
	out.WriteString(" = ")
	out.WriteString(a.Value.String())
	out.WriteString(";")

	return out.String()
}
