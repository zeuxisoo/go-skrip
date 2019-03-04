package ast

import (
	"bytes"
)

// Program is the root node of each AST in parser produces
type Program struct {
	Statements []Statement
}

// TokenLiteral is implements from Node interface
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, statement := range p.Statements {
		out.WriteString(statement.String())
	}

	return out.String()
}
