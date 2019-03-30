package parser

import (
	"fmt"
	"strconv"

	"github.com/zeuxisoo/go-skrip/lexer"
	"github.com/zeuxisoo/go-skrip/token"
	"github.com/zeuxisoo/go-skrip/ast"
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

	parser.nextToken()	// set the current token
	parser.nextToken() 	// set the peek token

	parser.prefixParseFunctions = make(map[token.Type]prefixParseFunction)
	parser.registerPrefixParseFunction(token.INT, parser.parseIntegerLiteral)
	parser.registerPrefixParseFunction(token.FLOAT, parser.parseFloatLiteral)
	parser.registerPrefixParseFunction(token.STRING, parser.parseStringLiteral)
	parser.registerPrefixParseFunction(token.FUNCTION, parser.parseFunctionLiteral)
	parser.registerPrefixParseFunction(token.IDENTIFIER, parser.parseIdentifier)
	parser.registerPrefixParseFunction(token.TRUE, parser.parseBoolean)
	parser.registerPrefixParseFunction(token.FALSE, parser.parseBoolean)
	parser.registerPrefixParseFunction(token.BANG, parser.parsePrefixExpression)
	parser.registerPrefixParseFunction(token.MINUS, parser.parsePrefixExpression)
	parser.registerPrefixParseFunction(token.LEFT_BRACKET, parser.parseArrayLiteral)
	parser.registerPrefixParseFunction(token.LEFT_BRACE, parser.parseHashLiteral)
	parser.registerPrefixParseFunction(token.LEFT_PARENTHESIS, parser.parseGroupedExpression)
	parser.registerPrefixParseFunction(token.IF, parser.parseIfExpression)
	parser.registerPrefixParseFunction(token.FOR, parser.parseForExpression)

	parser.infixParseFunctions = make(map[token.Type]infixParseFunction)
	parser.registerInfixParseFunction(token.PLUS, parser.parseInfixExpression)
	parser.registerInfixParseFunction(token.MINUS, parser.parseInfixExpression)
	parser.registerInfixParseFunction(token.ASTERISK, parser.parseInfixExpression)
	parser.registerInfixParseFunction(token.SLASH, parser.parseInfixExpression)
	parser.registerInfixParseFunction(token.LT, parser.parseInfixExpression)
	parser.registerInfixParseFunction(token.GT, parser.parseInfixExpression)
	parser.registerInfixParseFunction(token.LTEQ, parser.parseInfixExpression)
	parser.registerInfixParseFunction(token.GTEQ, parser.parseInfixExpression)
	parser.registerInfixParseFunction(token.EQ, parser.parseInfixExpression)
	parser.registerInfixParseFunction(token.NOT_EQ, parser.parseInfixExpression)
	parser.registerInfixParseFunction(token.AND, parser.parseInfixExpression)
	parser.registerInfixParseFunction(token.OR, parser.parseInfixExpression)
	parser.registerInfixParseFunction(token.RANGE, parser.parseInfixExpression)
	parser.registerInfixParseFunction(token.LEFT_BRACKET, parser.parseIndexExpression)
	parser.registerInfixParseFunction(token.LEFT_PARENTHESIS, parser.parseCallExpression)


	return parser
}

func (p *Parser) Parse() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for !p.currentTokenTypeIs(token.EOF) {
		statement := p.parseStatement()

		// Add statement node into root program root
		if statement != nil {
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
	case token.RETURN:
		return p.parseReturnStatement()
	case token.FUNCTION:
		// If next token is token.identifier, parse by function statement e.g. "func name() {}"
		// otherwise, parse by function literal expression e.g. "func() {}"
		if p.peekTokenTypeIs(token.IDENTIFIER) == true {
			return p.parseFunctionStatement()
		}else{
			return p.parseExpressionStatement()
		}
	default:
		return p.parseExpressionStatement()
	}
}

// Parse statement functions
func (p *Parser) parseLetStatement() *ast.LetStatement {
	// Set the LetStatement Token value is "let token struct"
	statement := &ast.LetStatement{
		Token: p.currentToken,
	}

	// If next token is identifier
	//		call nextToken() to set the current token point to this
	// otherwise, the token is not identifier
	//		return nil
	if p.expectPeekTokenTypeIs(token.IDENTIFIER) == false {
		return nil
	}

	// Set the LetStatement Name point to IdentifierExpression struct
	// and set the variable Token struct and variable name
	statement.Name = &ast.IdentifierExpression{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	// Ensure that next token is assign symbol, and set the current token point to this
	if p.expectPeekTokenTypeIs(token.ASSIGN) == false {
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

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{
		Token: p.currentToken,
	}

	// Move the current token to return value
	p.nextToken()

	// Parse the return value expression
	statement.ReturnValue = p.parseExpression(LOWEST)

	//
	if p.peekTokenTypeIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseFunctionStatement() *ast.FunctionStatement {
	// Set up function statement struct
	statement := &ast.FunctionStatement{
		Token: p.currentToken,
	}

	// Move to next token function name
	p.nextToken()

	// Parse and set function name
	statement.Name = &ast.IdentifierExpression{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	// Parse function literal expression
	statement.Function = p.parseFunctionLiteral().(*ast.FunctionLiteralExpression)

	//
	if p.peekTokenTypeIs(token.SEMICOLON) == true {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{
		Token: p.currentToken,
	}
	statement.Expression = p.parseExpression(LOWEST)

	for p.peekTokenTypeIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

// Parse statement functions
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// Workflow
	// E.g. 5 + 5;
	//  current token is 5,
	// 	1. so it is token.INT
	//  2. so it registered in []prefixParseFunctions list
	//	  2.1. so it will not call noPrefixParseFunctionError() add error message
	//  4. fire the prefix parse function from prefixParseFunctions[token.INT]
	//  5. continue lookup each token when found token.SEMICOLON and ensure precedence less than next token precedence

	// Get the prefix parse callback function from registered list like (if, function, !, -, 5, 5.1, "text")
	prefixParseFunction := p.prefixParseFunctions[p.currentToken.Type]

	// if the current token type is not registered in prefix parse function list, add the error
	if prefixParseFunction == nil {
		p.noPrefixParseFunctionError(p.currentToken.Type)

		return nil
	}

	// If the current token type is registered in prefix parse function list, fire the prefix parse function
	leftExpression := prefixParseFunction()

	// Continue lookup the following tokens
	// Loop each token
	// 		unitil found semicolon token
	// 		when current token precedence is greater than LOWEST precedence
	for p.peekTokenTypeIs(token.SEMICOLON) == false && precedence < p.peekPrecedence() {
		// Get the  infix parse callback function name from registered list like (operator: + plus, - minus, etc)
		infixParseFunction := p.infixParseFunctions[p.peekToken.Type]

		// if the next token type is not registered in infix parse function list, only return parsed current token
		// otherwise, set next token to current token, and pass the parsed previous token to infix parse function
		// e.g. 5 + 6;
		// - current token: 5
		// 	 - left expression: 5
		// - peek token: +
		// - infix parse function found
		// - set current token: +
		// - pass left express (5) to found infix parse function
		if infixParseFunction == nil {
			return leftExpression
		}else{
			p.nextToken()

			leftExpression = infixParseFunction(leftExpression)
		}
	}

	return leftExpression
}

func (p *Parser) parseExpressionList(endTokenType token.Type) []ast.Expression {
	expressions := []ast.Expression{}

	// If next token is equals end token type like "]" and "}" etc
	// update the current and next token and then return the expression list
	if p.peekTokenTypeIs(endTokenType) {
		p.nextToken()

		return expressions
	}

	// Otherwise, set the current token to first element, and set next token like "," or ")" or "}" etc
	p.nextToken()

	// Parse the first element and add to expression list
	expressions = append(expressions, p.parseExpression(LOWEST))

	// Loop when found comma again and again
	// and add each found element into expression list
	for p.peekTokenTypeIs(token.COMMA) == true {
		p.nextToken()	// set current token to ","
		p.nextToken()	// set current token to next element

		expressions = append(expressions, p.parseExpression(LOWEST))
	}

	// If next token is equals end token type like "]" and "}" etc
	// update the current token to this end token and next token
	if p.expectPeekTokenTypeIs(endTokenType) == false {
		return nil
	}

	return expressions
}

// Parse prefix/infix functions
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
	if p.expectPeekTokenTypeIs(token.LEFT_PARENTHESIS) == false {
		return nil
	}

	functionLiteralExpression.Parameters = p.parseFunctionParameters()

	// Expect next token is "{"
	if p.expectPeekTokenTypeIs(token.LEFT_BRACE) == false {
		return nil
	}

	functionLiteralExpression.Block = p.parseBlockStatement()

	return functionLiteralExpression
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	arrayLiteralExpression := &ast.ArrayLiteralExpression{
		Token: p.currentToken,
	}
	arrayLiteralExpression.Elements = p.parseExpressionList(token.RIGHT_BRACKET)

	return arrayLiteralExpression
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hashLiteralExpression := &ast.HashLiteralExpression{
		Token: p.currentToken,
	}

	hashLiteralExpression.Pairs = make(map[ast.Expression]ast.Expression)

	// Loop until found "}"
	for p.peekTokenTypeIs(token.RIGHT_BRACE) == false {
		// Set current token to key token
		p.nextToken()

		// Parse current/key token expression and assign to key variable
		key := p.parseExpression(LOWEST)

		// If next token is not ":", return nil. Otherwise update current token to this ":"
		if p.expectPeekTokenTypeIs(token.COLON) == false {
			return nil
		}

		// Set current token to value token
		p.nextToken()

		// Parse current/value token expression and assign to value variable
		value := p.parseExpression(LOWEST)

		// Update the pairs map data
		hashLiteralExpression.Pairs[key] = value

		// If next token is not "}" and it will expect next token it is "," and update the current token to this
		// otherwise return nil to break the loop
		if p.peekTokenTypeIs(token.RIGHT_BRACE) == false && p.expectPeekTokenTypeIs(token.COMMA) == false {
			return nil
		}
	}

	// End of loop, if the next token is not "}", return nil
	if p.expectPeekTokenTypeIs(token.RIGHT_BRACE) == false {
		return nil
	}

	return hashLiteralExpression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	// Move current token "(" to next token
	p.nextToken()

	expression := p.parseExpression(LOWEST)

	// If the next token is ")", set current token
	// otherwise return nil
	if p.expectPeekTokenTypeIs(token.RIGHT_PARENTHESIS) == false {
		return nil
	}

	return expression
}

func (p *Parser) parseIfExpression() ast.Expression {
	ifExpression := &ast.IfExpression{
		Token: p.currentToken,
	}

	// If next token is "(", set current token to this,
	// otherwise, return nil
	if p.expectPeekTokenTypeIs(token.LEFT_PARENTHESIS) == false {
		return nil
	}

	// Move current token to if condition
	p.nextToken()

	// Parse the condition between "(" and ")"
	ifExpression.Condition = p.parseExpression(LOWEST)

	// If next token is ")", set current token to this
	// otherwise, return nil
	if p.expectPeekTokenTypeIs(token.RIGHT_PARENTHESIS) == false {
		return nil
	}

	// If next token is "{", set current token to this
	// otherwise, return nil
	if p.expectPeekTokenTypeIs(token.LEFT_BRACE) == false {
		return nil
	}

	// Parse the block after if condition
	ifExpression.Block = p.parseBlockStatement()

	// When found else after if condition block, parse the else block
	if p.peekTokenTypeIs(token.ELSE) == true {
		p.nextToken()

		if p.expectPeekTokenTypeIs(token.LEFT_BRACE) == false {
			return nil
		}

		ifExpression.Alternative = p.parseBlockStatement()
	}

	return ifExpression
}

func (p *Parser) parseForExpression() ast.Expression {
	// When found "{", mean "for { ... }"
	if p.peekTokenTypeIs(token.LEFT_BRACE) == true {
		return p.parseForEverExpression(p.currentToken)
	}

	// Save current token (token.FOR) for forEachHash and forEachArray
	tokenFor := p.currentToken

	// If next token is not identifier, stop it and return nil
	// otherwise set current token to this
	if p.expectPeekTokenTypeIs(token.IDENTIFIER) == false {
		return nil
	}

	// When next token is ",", mean "for key, value in hash { ... }"
	// otherwise mean "for value in array { ... }"
	if p.peekTokenTypeIs(token.COMMA) == true {
		return p.parseForEachHashExpression(tokenFor, p.currentToken)
	}

	return p.parseForEachArrayOrRangeExpression(tokenFor, p.currentToken)
}

func (p *Parser) parseIdentifier() ast.Expression {
	identifier := &ast.IdentifierExpression{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	return identifier
}

func (p *Parser) parseBoolean() ast.Expression {
	boolean := &ast.BooleanExpression{
		Token: p.currentToken,
		Value: p.currentTokenTypeIs(token.TRUE),
	}

	return boolean
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	prefix := &ast.PrefixExpression{
		Token   : p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.nextToken()

	prefix.Right = p.parseExpression(PREFIX)

	return prefix
}

func (p *Parser) parseInfixExpression(leftExpression ast.Expression) ast.Expression {
	infix := &ast.InfixExpression{
		Token   : p.currentToken,
		Left    : leftExpression,
		Operator: p.currentToken.Literal,
	}

	precedence := p.currentPrecedence()

	p.nextToken()

	infix.Right = p.parseExpression(precedence)

	return infix
}

func (p *Parser) parseIndexExpression(leftExpression ast.Expression) ast.Expression {
	index := &ast.IndexExpression{
		Token: p.currentToken,
		Left : leftExpression,
	}

	p.nextToken()

	index.Index = p.parseExpression(LOWEST)

	// If next token is }, update the current token to this
	// otherwise, return nil
	if p.expectPeekTokenTypeIs(token.RIGHT_BRACKET) == false {
		return nil
	}

	return index
}

func (p *Parser) parseCallExpression(leftExpression ast.Expression) ast.Expression {
	call := &ast.CallExpression{
		Token   : p.currentToken,
		Function: leftExpression,
	}

	call.Arguments = p.parseExpressionList(token.RIGHT_PARENTHESIS)

	return call
}

// Helper functions
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken    = p.lexer.NextToken()
}

func (p *Parser) expectPeekTokenTypeIs(tokenType token.Type) bool {
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
	if precedence, ok := precedences[p.peekToken.Type]; ok {
		return precedence
	}

	return LOWEST
}

func (p *Parser) currentPrecedence() int {
	if precedence, ok := precedences[p.currentToken.Type]; ok {
		return precedence
	}

	return LOWEST
}

func (p *Parser) registerPrefixParseFunction(tokenType token.Type, callback prefixParseFunction) {
	p.prefixParseFunctions[tokenType] = callback
}

func (p *Parser) registerInfixParseFunction(tokenType token.Type, callback infixParseFunction) {
	p.infixParseFunctions[tokenType] = callback
}

// Helper function for parse prefix function like function arguments, function block and so on
func (p *Parser) parseFunctionParameters() []*ast.IdentifierExpression {
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

// Helper function for parse for loop case
func (p *Parser) parseForEverExpression(currentToken token.Token) ast.Expression {
	forEverExpression := &ast.ForEverExpression{
		Token: currentToken,
	}

	if p.expectPeekTokenTypeIs(token.LEFT_BRACE) == false {
		return nil
	}

	forEverExpression.Block = p.parseBlockStatement()

	return forEverExpression
}

func (p *Parser) parseForEachHashExpression(tokenFor token.Token, currentToken token.Token) ast.Expression {
	forEachHashExpression := &ast.ForEachHashExpression{
		Token: tokenFor,
		Key  : currentToken.Literal,
	}

	// If next token is ",", set current token to it
	// otherwise, return nil and stop
	if p.expectPeekTokenTypeIs(token.COMMA) == false {
		return nil
	}

	// If next token is identifier, set current token to it
	// otherwise, return nil and stop
	if p.expectPeekTokenTypeIs(token.IDENTIFIER) == false {
		return nil
	}

	// Set for loop value name
	forEachHashExpression.Value = p.currentToken.Literal

	// If next token is "in", set current token to it
	// otherwise, return nil and stop
	if p.expectPeekTokenTypeIs(token.IN) == false {
		return nil
	}

	// Enter to loop data
	p.nextToken()

	// Parse the loop data
	forEachHashExpression.Iterable = p.parseExpression(LOWEST)

	// If next token is "{", set current token to it
	// otherwise, return nil and stop
	if p.expectPeekTokenTypeIs(token.LEFT_BRACE) == false {
		return nil
	}

	forEachHashExpression.Block = p.parseBlockStatement()

	return forEachHashExpression
}

func (p *Parser) parseForEachArrayOrRangeExpression(tokenFor token.Token, currentToken token.Token) ast.Expression {
	forEachArrayOrRangeExpression := &ast.ForEachArrayOrRangeExpression{
		Token: tokenFor,
		Value: currentToken.Literal,
	}

	// If next token is "in", set current token to it
	// otherwise, return nil and stop
	if p.expectPeekTokenTypeIs(token.IN) == false {
		return nil
	}

	// Enter to loop data
	p.nextToken()

	// Parse the loop data
	forEachArrayOrRangeExpression.Iterable = p.parseExpression(LOWEST)

	// If next token is "{", set current token to it
	// otherwise, return nil and stop
	if p.expectPeekTokenTypeIs(token.LEFT_BRACE) == false {
		return nil
	}

	forEachArrayOrRangeExpression.Block = p.parseBlockStatement()

	return forEachArrayOrRangeExpression
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
