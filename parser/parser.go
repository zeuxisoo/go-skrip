package parser

import (
	"strings"

	"github.com/zeuxisoo/go-skriplang/lexer"
	"github.com/zeuxisoo/go-skriplang/token"
	"github.com/zeuxisoo/go-skriplang/ast"
)

type Parser struct {
	lexer 	*lexer.Lexer
	errors  errorStrings

	currentToken token.Token
	peekToken 	 token.Token
}

//
func NewParser(lexer *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer: 	lexer,
		errors: []string{},
	}

	return parser
}

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

//
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
