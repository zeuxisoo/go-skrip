package parser

import (
	"github.com/zeuxisoo/go-skriplang/lexer"
)

type Parser struct {
	lexer 	*lexer.Lexer
	errors  errorStrings
}

func NewParser(lexer *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer: 	lexer,
		errors: []string{},
	}

	return parser
}
