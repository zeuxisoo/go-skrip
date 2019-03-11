package parser

import (
	"fmt"
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
	case token.LET:
		return p.parseLetStatement()
	case token.EOF:
		return nil
	default:
		// TODO: default action
		return nil
	}
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	// TODO: parse expression
	return nil
}

//
func (p *Parser) parseLetStatement() *ast.LetStatement {
	// Set the LetStatement Token value is "let token strcut"
	statement := &ast.LetStatement{
		Token: p.currentToken,
	}

	// If next token is identifier
	//		call nextToken() to set the current token point to this
	// otherwise, the token is not identifier
	//		return nil
	if p.expectPeekTokenType(token.IDENTIFIER) == false {
		return nil
	}

	// Set the LetStatement Name point to IdentifierExpression struct
	// and set the variable Token struct and variable name
	statement.Name = &ast.IdentifierExpression{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	// Ensure that next token is assign symbol, and set the current token point to this
	if p.expectPeekTokenType(token.ASSIGN) == false {
		return nil
	}

	// Move the current token again let it point to variable value
	p.nextToken()

	// Set variable value by parsed expression
	statement.Value = p.parseExpression(LOWEST)

	//
	if p.peekTokenTypeIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}



func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken    = p.lexer.NextToken()
}

func (p *Parser) expectPeekTokenType(tokenType token.Type) bool {
	if p.peekTokenTypeIs(tokenType) {
		p.nextToken()
		return true
	}

	p.peekTokenTypeError(tokenType)

	return false
}

func (p *Parser) peekTokenTypeIs(tokenType token.Type) bool {
	return p.peekToken.Type == tokenType
}

func (p *Parser) peekTokenTypeError(tokenType token.Type) {
	message := fmt.Sprintf("Expected peek token type should be %s, but got %s", tokenType, p.peekToken.Type)
	p.errors = append(p.errors, message)
}
