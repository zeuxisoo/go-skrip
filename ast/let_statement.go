package ast

import (
	"bytes"

	"github.com/zeuxisoo/go-skrip/token"
)

type LetStatement struct {
	Token token.Token
	Name  *IdentifierExpression
	Value Expression
}

func (l *LetStatement) statementNode() {
}

// Implement methods for Node interface
func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}

func (l *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(l.TokenLiteral() + " ")	// let
	out.WriteString(l.Name.String())		// variable
	out.WriteString(" = ")					// =

	if l.Value != nil {
		out.WriteString(l.Value.String())	// Value
	}

	out.WriteString(";")					// ;

	return out.String()
}
