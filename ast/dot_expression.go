package ast

import (
	"bytes"

	"github.com/zeuxisoo/go-skrip/token"
)

type DotExpression struct {
	Token token.Token
	Left  Expression
	Item  Expression
}

func (d *DotExpression) expressionNode() {
}

// Implement methods for Node interface
func (d *DotExpression) TokenLiteral() string {
	return d.Token.Literal
}

func (d *DotExpression) String() string {
	var out bytes.Buffer

	out.WriteString(d.Left.String()) // object
	out.WriteString(".")             // .
	out.WriteString(d.Item.String()) // item

	return out.String()
}
