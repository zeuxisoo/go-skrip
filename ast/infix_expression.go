package ast

import (
	"bytes"

	"github.com/zeuxisoo/go-skrip/token"
)

type InfixExpression struct {
	Token		token.Token
	Left	 	Expression
	Operator 	string
	Right 		Expression
}

func (i *InfixExpression) expressionNode() {
}

// Implement methods for Node interface
func (i *InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString(" " + i.Operator + " ")
	out.WriteString(i.Right.String())
	out.WriteString(")")

	return out.String()
}
