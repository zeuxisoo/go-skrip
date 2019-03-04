package ast

import (
	"github.com/zeuxisoo/go-skriplang/token"
)

type Identifier struct {
	Token token.Token
	Value string
}

func (i Identifier) expressionNode() {

}
