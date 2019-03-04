package parser

import (
	"github.com/zeuxisoo/go-skriplang/lexer"
	"github.com/zeuxisoo/go-skriplang/token"
	"github.com/zeuxisoo/go-skriplang/ast"
)

type (
	prefixParseFunction ast.Expression
)

type Parser struct {
	lexer 	*lexer.Lexer
	errors  errorStrings

	prefixParseFunctions map[token.Type]prefixParseFunction
}

func (p *Parser) registerPrefix(tokenType token.Type, callback prefixParseFunction) {
	p.prefixParseFunctions[tokenType] = callback
}

func NewParser(lexer *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer: 	lexer,
		errors: []string{},
	}

	parser.prefixParseFunctions = make(map[token.Type]prefixParseFunction)

	return parser
}
