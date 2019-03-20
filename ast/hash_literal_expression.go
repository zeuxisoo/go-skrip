package ast

import (
	"bytes"
	"strings"

	"github.com/zeuxisoo/go-skrip/token"
)

type HashLiteralExpression struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (h *HashLiteralExpression) expressionNode() {
}

// Implement methods for Node interface
func (h *HashLiteralExpression) TokenLiteral() string {
	return h.Token.Literal
}

func (h *HashLiteralExpression) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range h.Pairs {
		pairs = append(pairs, key.String() + ":" + value.String())
	}

	out.WriteString("{")						// {
	out.WriteString(strings.Join(pairs, ", "))	// key:value, key:value ...
	out.WriteString("}")						// }

	return out.String()
}
