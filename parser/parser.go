package parser

import (
	"strings"

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
	peekToken 	 token.Token
}

//
func (p *Parser) Parse() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for !p.CurrentTokenEquals(token.EOF) {
		statement := p.parseStatement()

		// Add statement node into root program root
		if statement != nil && strings.TrimSpace(statement.String()) != "" {
			program.Statements = append(program.Statements, statement)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) CurrentTokenEquals(tokenType token.Type) bool {
	return p.currentToken.Type == tokenType
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.EOF:
		return nil
	default:
		// TODO: default action
		return nil
	}
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken    = p.lexer.NextToken()
}

//
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
