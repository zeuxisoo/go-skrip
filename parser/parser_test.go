package parser

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/zeuxisoo/go-skriplang/lexer"
	"github.com/zeuxisoo/go-skriplang/ast"
)

type expectedStatement struct {
	source 		string
	identifier 	string
	value		interface{}
}

// Test case
func TestLetStatement(t *testing.T) {
	Convey("Let statement testing", t, func() {
		expectedStatements := []expectedStatement{
			{ "let a = 5;",		"a",	5 },
			{ "let b = 5.1",	"b",	5.1 },
			{ "let c = true",	"c",	true },
			{ "let d = c",		"d",	"c" },
		}

		testLetStatement(expectedStatements)
	})
}

// Sub method for test case
func testLetStatement(expectedStatements []expectedStatement) {
	for index, currentStatement := range expectedStatements {
		message := runMessage("Running %d, Source: %s", index, currentStatement.source)

		theLexer   := lexer.NewLexer(currentStatement.source)
		theParser  := NewParser(theLexer)
		theProgram := theParser.Parse()

		statement := theProgram.Statements[0]
		letStatement, ok := statement.(*ast.LetStatement)

		Convey(message, func() {
			testParserError(theParser)
			testParserProgramLength(theProgram)

			// Identifier
			So(statement.TokenLiteral(), ShouldEqual, "let")
			So(ok, ShouldNotBeNil)
			So(letStatement.Name.Value, ShouldEqual, currentStatement.identifier)
			So(letStatement.Name.TokenLiteral(), ShouldEqual, currentStatement.identifier)

			// Value
			testLiteralExpression(letStatement.Value, currentStatement.value)
		})
	}
}

// Sub method function for sub method
func testParserError(parser *Parser) {
	parserErrors       := parser.Errors()
	parserErrorsLength := len(parserErrors)

	So(parserErrorsLength, ShouldEqual, 0)
}

func testParserProgramLength(program *ast.Program) {
	So(len(program.Statements), ShouldEqual, 1)
}

func testLiteralExpression(expression ast.Expression, expected interface{}) {
	switch value := expected.(type) {
	case int:
		testIntegerLiteralExpression(expression, int64(value))
	case int64:
		testIntegerLiteralExpression(expression, value)
	case string:
		testIdentifierExpression(expression, value)
	case bool:
		testBooleanExpression(expression, value)
	case float32:
		testFloatLiteralExpression(expression, float64(value))
	case float64:
		testFloatLiteralExpression(expression, value)
	}
}

// Callback function for sub method function "testLiteralExpression"
func testIdentifierExpression(expression ast.Expression, value string) {
	identifier, ok := expression.(*ast.IdentifierExpression)

	identifierValue 	   := identifier.Value
	identifierTokenLiteral := identifier.TokenLiteral()

	Convey(
		runMessage(
			"Identifier test, Value: %s, TokenLiteral: %s, Expected: %s",
			identifierValue, identifierTokenLiteral, value,
		),
		func() {
			So(ok, ShouldBeTrue)
			So(identifier.Value, ShouldEqual, value)
			So(identifier.TokenLiteral(), ShouldEqual, value)
		},
	)
}

func testBooleanExpression(expression ast.Expression, value bool) {
	boolean, ok := expression.(*ast.BooleanExpression)

	booleanValue 		:= boolean.Value
	booleanTokenLiteral := boolean.TokenLiteral()

	Convey(
		runMessage(
			"Boolean test, Value: %t, TokenLiteral: %s, Expected: %s",
			booleanValue, booleanTokenLiteral, fmt.Sprintf("%t", value),
		),
		func() {
			So(ok, ShouldBeTrue)
			So(boolean.Value, ShouldEqual, value)
			So(boolean.TokenLiteral(), ShouldEqual, fmt.Sprintf("%t", value))
		},
	)
}

func testIntegerLiteralExpression(expression ast.Expression, value int64) {
	integer, ok := expression.(*ast.IntegerLiteralExpression)

	integerValue 		:= integer.Value
	integerTokenLiteral := integer.TokenLiteral()

	Convey(
		runMessage(
			"Integer test, Value: %d, TokenLiteral: %s, Expected: %d",
			integerValue, integerTokenLiteral, value,
		),
		func() {
			So(ok, ShouldBeTrue)
			So(integer.Value, ShouldEqual, value)
			So(integer.TokenLiteral(), ShouldEqual, fmt.Sprintf("%d", value))
		},
	)
}

func testFloatLiteralExpression(expression ast.Expression, value float64) {
	float, ok := expression.(*ast.FloatLiteralExpression)

	floatValue 		  := fmt.Sprintf("%.1f", float.Value)
	floatTokenLiteral := float.TokenLiteral()
	expectedValue     := fmt.Sprintf("%.1f", value)

	Convey(
		runMessage(
			"Float literal test, Value: %s, TokenLiteral: %s, Expected: %s",
			floatValue, floatTokenLiteral, expectedValue,
		),
		func() {
			So(ok, ShouldBeTrue)
			So(floatValue, ShouldEqual, expectedValue)
			So(floatTokenLiteral, ShouldEqual, expectedValue)
		},
	)
}

// Helper functions for common
func runMessage(format string, values ...interface{}) string {
	return fmt.Sprintf(format, values...)
}
