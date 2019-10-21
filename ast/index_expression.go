package ast

import (
	"bytes"

	"github.com/zeuxisoo/go-skrip/token"
)

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (i *IndexExpression) expressionNode() {
}

// Implement methods for Node interface
func (i *IndexExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")              // (
	out.WriteString(i.Left.String())  // object/variable
	out.WriteString("[")              // [
	out.WriteString(i.Index.String()) // index
	out.WriteString("]")              // ]
	out.WriteString(")")              // )

	return out.String()
}
