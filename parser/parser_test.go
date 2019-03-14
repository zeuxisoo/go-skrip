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
			checkParserError(theParser)
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
func checkParserError(parser *Parser) {
	Convey("Parser error test", func() {
		parserErrors       := parser.Errors()
		parserErrorsLength := len(parserErrors)

		So(parserErrorsLength, ShouldEqual, 0)
	})
}

func testParserProgramLength(program *ast.Program) {
	Convey("Paser program length test", func() {
		So(len(program.Statements), ShouldEqual, 1)
	})
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
	Convey("Identifier test", func() {
		identifier, ok := expression.(*ast.IdentifierExpression)

		So(ok, ShouldBeTrue)

		Convey(runMessage("1. Got: %s, Expected: %s", identifier.Value, value), func() {
			So(identifier.Value, ShouldEqual, value)
		})

		Convey(runMessage("2. Got: %s, Expected: %s", identifier.TokenLiteral(), value), func() {
			So(identifier.TokenLiteral(), ShouldEqual, value)
		})
	})
}

func testBooleanExpression(expression ast.Expression, value bool) {
	Convey("Boolean test", func() {
		boolean, ok := expression.(*ast.BooleanExpression)

		So(ok, ShouldBeTrue)

		Convey(runMessage("1. Got: %t, Expected: %t", boolean.Value, value), func() {
			So(boolean.Value, ShouldEqual, value)
		})

		Convey(runMessage("2. Got: %s, Expected: %t", boolean.TokenLiteral(), value), func() {
			So(boolean.TokenLiteral(), ShouldEqual, fmt.Sprintf("%t", value))
		})
	})
}

func testIntegerLiteralExpression(expression ast.Expression, value int64) {
	Convey("Integer literal test", func() {
		integer, ok := expression.(*ast.IntegerLiteralExpression)

		So(ok, ShouldBeTrue)

		Convey(runMessage("1. Got: %d, Expected: %d", integer.Value, value), func() {
			So(integer.Value, ShouldEqual, value)
		})

		Convey(runMessage("2. Got: %s, Expected: %d", integer.TokenLiteral(), value), func() {
			So(integer.TokenLiteral(), ShouldEqual, fmt.Sprintf("%d", value))
		})
	})
}

func testFloatLiteralExpression(expression ast.Expression, value float64) {
	Convey("Float literal test", func() {
		float, ok := expression.(*ast.FloatLiteralExpression)

		So(ok, ShouldBeTrue)

		Convey(runMessage("1. Got: %f, Expected: %f", float.Value, value), func() {
			So(float.Value, ShouldEqual, value)
		})
	})
}

// Helper functions for common
func runMessage(format string, values ...interface{}) string {
	return fmt.Sprintf(format, values...)
}
