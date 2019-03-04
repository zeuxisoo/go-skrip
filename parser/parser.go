package parser

import (
	"github.com/zeuxisoo/go-skriplang/lexer"
	"github.com/zeuxisoo/go-skriplang/token"
	"github.com/zeuxisoo/go-skriplang/ast"
)

type (
	prefixParseFunction func() ast.Expression
)

type Parser struct {
	lexer 	*lexer.Lexer
	errors  errorStrings

	prefixParseFunctions map[token.Type]prefixParseFunction

	currentToken token.Token
}

func (p *Parser) registerPrefix(tokenType token.Type, callback prefixParseFunction) {
	p.prefixParseFunctions[tokenType] = callback
}

func (p *Parser) parserIdentifier() ast.Expression {
	identifier := &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	return identifier
}

func NewParser(lexer *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer: 	lexer,
		errors: []string{},
	}

	parser.prefixParseFunctions = make(map[token.Type]prefixParseFunction)
	parser.registerPrefix(token.IDENTIFIER, parser.parserIdentifier)

	return parser
}
