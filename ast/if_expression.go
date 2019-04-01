package ast

import (
	"bytes"

	"github.com/zeuxisoo/go-skrip/token"
)

type IfScene struct {
	Condition Expression
	Block     *BlockStatement
}

type IfExpression struct {
	Token 		token.Token
	Scenes		[]*IfScene
	Alternative *BlockStatement
}

func (i *IfExpression) expressionNode() {
}

// Implement methods for Node interface
func (i *IfExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IfExpression) String() string {
	var out bytes.Buffer

	for index, scene := range i.Scenes {
		if index == 0 {
			out.WriteString("if ")
			out.WriteString(scene.Condition.String())
			out.WriteString(" { ")
			out.WriteString(scene.Block.String())
			out.WriteString(" } ")
		}else{
			out.WriteString("else if ")
			out.WriteString(scene.Condition.String())
			out.WriteString(" { ")
			out.WriteString(scene.Block.String())
			out.WriteString(" } ")
		}
	}

	if i.Alternative != nil {
		out.WriteString("else")
		out.WriteString(" { ")
		out.WriteString(i.Alternative.String())
		out.WriteString(" }")
	}

	return out.String()
}
