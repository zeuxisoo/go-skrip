package ast

import (
	"bytes"

	"github.com/zeuxisoo/go-skriplang/token"
)

type BlockStatement struct {
	Token 		token.Token
	Statements 	[]Statement
}

func (b *BlockStatement) statementNode() {
}

// Implement methods for Node interface
func (b *BlockStatement) TokenLiteral() string {
	return b.Token.Literal
}

func (b *BlockStatement) String() string {
	var out bytes.Buffer

	for _, statement := range b.Statements {
		out.WriteString(statement.String())
	}

	return out.String()
}
