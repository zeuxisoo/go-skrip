package ast

import (
	"bytes"
	"strings"

	"github.com/zeuxisoo/go-skrip/token"
)

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (c *CallExpression) expressionNode() {
}

// Implement methods for Node interface
func (c *CallExpression) TokenLiteral() string {
	return c.Token.Literal
}

func (c *CallExpression) String() string {
	var out bytes.Buffer

	arguments := []string{}
	for _, argument := range c.Arguments {
		arguments = append(arguments, argument.String())
	}

	out.WriteString(c.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(arguments, ", "))
	out.WriteString(")")

	return out.String()
}
