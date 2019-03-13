package parser

import (
	"fmt"
	"strings"
	"strconv"

	"github.com/zeuxisoo/go-skriplang/lexer"
	"github.com/zeuxisoo/go-skriplang/token"
	"github.com/zeuxisoo/go-skriplang/ast"
)

type (
	prefixParseFunction func() ast.Expression
	infixParseFunction	func(ast.Expression) ast.Expression
)

type Parser struct {
	lexer 	*lexer.Lexer
	errors  errorStrings

	currentToken token.Token
	peekToken 	 token.Token

	prefixParseFunctions map[token.Type]prefixParseFunction
	infixParseFunctions	 map[token.Type]infixParseFunction
}

// Public functions
func NewParser(lexer *lexer.Lexer) *Parser {
	parser := &Parser{
		lexer: 	lexer,
		errors: []string{},
	}

	parser.prefixParseFunctions = make(map[token.Type]prefixParseFunction)
	parser.registerPrefixParseFunction(token.INT, parser.parseIntegerLiteral)
	parser.registerPrefixParseFunction(token.FLOAT, parser.parseFloatLiteral)
	parser.registerPrefixParseFunction(token.STRING, parser.parseStringLiteral)
	parser.registerPrefixParseFunction(token.FUNCTION, parser.parseFunctionLiteral)

	return parser
}

func (p *Parser) Parse() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for !p.currentTokenTypeIs(token.EOF) {
		statement := p.parseStatement()

		// Add statement node into root program root
		if statement != nil && strings.TrimSpace(statement.String()) != "" {
			program.Statements = append(program.Statements, statement)
		}

		p.nextToken()
	}

	return program
}

func (p *Parser) Errors() []string {
	return p.errors
}

// Parse functions
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
	// Get the prefix parse callback function name like (keywords: IF, FUNCTION, etc)
	prefixParseFunction := p.prefixParseFunctions[p.currentToken.Type]

	if prefixParseFunction == nil {
		p.noPrefixParseFunctionError(p.currentToken.Type)

		return nil
	}

	// Fire the prefix parse callback function
	leftExpression := prefixParseFunction()

	// Loop each token
	// 		unitil found semicolon token
	// 		when current token precedence is greater than LOWEST precedence
	for p.peekTokenTypeIs(token.SEMICOLON) == false && precedence < p.peekPrecedence() {
		// Get the  infix parse callback function name like (operator: + plus, - minus, etc)
		infixParseFunction := p.infixParseFunctions[p.peekToken.Type]

		if infixParseFunction == nil {
			return leftExpression
		}else{
			p.nextToken()

			leftExpression = infixParseFunction(leftExpression)
		}
	}

	return leftExpression
}

// Parse statement functions
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

// Parse prefix functions
func (p *Parser) parseIntegerLiteral() ast.Expression {
	integerLiteralExpression := &ast.IntegerLiteralExpression{
		Token: p.currentToken,
	}

	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(
			p.errors,
			fmt.Sprintf("Line: %d, Can not parse %q as integer", p.currentToken.LineNumber, p.currentToken.Literal),
		)

		return nil
	}

	integerLiteralExpression.Value = value

	return integerLiteralExpression
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	floatLiteralExpression := &ast.FloatLiteralExpression{
		Token: p.currentToken,
	}

	value, err := strconv.ParseFloat(p.currentToken.Literal, 64)
	if err != nil {
		p.errors = append(
			p.errors,
			fmt.Sprintf("Line: %d, Can not parse %q as float", p.currentToken.LineNumber, p.currentToken.Literal),
		)

		return nil
	}

	floatLiteralExpression.Value = value

	return floatLiteralExpression
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteralExpression{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	functionLiteralExpression := &ast.FunctionLiteralExpression{
		Token: p.currentToken,
	}

	// Expect next token is "("
	if p.expectPeekTokenType(token.LEFT_PARENTHESIS) == false {
		return nil
	}

	functionLiteralExpression.Parameters = p.parsePrefixFunctionParameters()

	// Expect next token is "{"
	if p.expectPeekTokenType(token.LEFT_BRACE) == false {
		return nil
	}

	functionLiteralExpression.Block = p.parseBlockStatement()

	return functionLiteralExpression
}

// Helper functions
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

func (p *Parser) currentTokenTypeIs(tokenType token.Type) bool {
	return p.currentToken.Type == tokenType
}

func (p *Parser) peekTokenTypeIs(tokenType token.Type) bool {
	return p.peekToken.Type == tokenType
}

func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedences[p.currentToken.Type]; ok {
		return precedence
	}

	return LOWEST
}

func (p *Parser) registerPrefixParseFunction(tokenType token.Type, callback prefixParseFunction) {
	p.prefixParseFunctions[tokenType] = callback
}

// Helper function for parse function, block and etc
func (p *Parser) parsePrefixFunctionParameters() []*ast.IdentifierExpression {
	identifierExpressions := []*ast.IdentifierExpression{}

	// If the next token is ")", it means no arguments
	// so, move to next token and return empty arguments
	if p.peekTokenTypeIs(token.RIGHT_PARENTHESIS) {
		p.nextToken()

		return identifierExpressions
	}

	// Current in "(", so move it to next token
	p.nextToken()

	// Loop until found ")"
	for p.currentTokenTypeIs(token.RIGHT_PARENTHESIS) == false {
		// Append the parameter identifier to parameter identifiers
		identifierExpression := &ast.IdentifierExpression{
			Token: p.currentToken,
			Value: p.currentToken.Literal,
		}
		identifierExpressions = append(identifierExpressions, identifierExpression)

		// Move to next token
		p.nextToken()

		// If current token is ",", skip it then move to next token
		if p.currentTokenTypeIs(token.COMMA) == true {
			p.nextToken()
		}
	}

	return identifierExpressions
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	blockStatement := &ast.BlockStatement{
		Token: p.currentToken,
	}

	blockStatement.Statements = []ast.Statement{}

	// Move to next token from "{"
	p.nextToken()

	// Loop until found "}"
	for p.currentTokenTypeIs(token.RIGHT_BRACE) == false {
		statement := p.parseStatement()

		if statement != nil {
			blockStatement.Statements = append(blockStatement.Statements, statement)
		}

		p.nextToken()
	}

	return blockStatement
}

// Error handle functions
func (p *Parser) peekTokenTypeError(tokenType token.Type) {
	message := fmt.Sprintf("Line: %d, Expected peek token type should be %s, but got %s", p.currentToken.LineNumber, tokenType, p.peekToken.Type)
	p.errors = append(p.errors, message)
}

func (p *Parser) noPrefixParseFunctionError(tokenType token.Type) {
	message := fmt.Sprintf("Line: %d, Can not found related prefix parse function for %s", p.currentToken.LineNumber, tokenType)
	p.errors = append(p.errors, message)
}
