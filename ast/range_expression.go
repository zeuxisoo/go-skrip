package ast

import (
	"bytes"

	"github.com/zeuxisoo/go-skrip/token"
)

type RangeExpression struct {
	Token token.Token
	Start Expression
	End   Expression
}

func (i *RangeExpression) expressionNode() {
}

// Implement methods for Node interface
func (i *RangeExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *RangeExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")              // (
	out.WriteString(i.Start.String()) // object/variable
	out.WriteString("..")             // ..
	out.WriteString(i.End.String())   // object/variable
	out.WriteString(")")              // )

	return out.String()
}
