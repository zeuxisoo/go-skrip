package ast

import (
	"bytes"

	"github.com/zeuxisoo/go-skriplang/token"
)

type ReturnStatement struct {
	Token 		token.Token
	ReturnValue Expression
}

func (r *ReturnStatement) statementNode() {
}

// Implement methods for Node interface
func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

func (r *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(r.TokenLiteral() + " ")		// return

	if r.ReturnValue != nil {
		out.WriteString(r.ReturnValue.String())	// value
	}

	out.WriteString(";")						// ;

	return out.String()
}
