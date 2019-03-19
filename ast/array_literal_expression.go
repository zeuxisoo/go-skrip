package ast

import (
	"bytes"
	"strings"

	"github.com/zeuxisoo/go-skrip/token"
)

type ArrayLiteralExpression struct {
	Token 		token.Token
	Elements 	[]Expression
}

func (a *ArrayLiteralExpression) expressionNode() {
}

// Implement methods for Node interface
func (a *ArrayLiteralExpression) TokenLiteral() string {
	return a.Token.Literal
}

func (a *ArrayLiteralExpression) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, element := range a.Elements {
		elements = append(elements, element.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
