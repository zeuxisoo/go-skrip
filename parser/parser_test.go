package parser

import (
	"fmt"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/zeuxisoo/go-skrip/lexer"
	"github.com/zeuxisoo/go-skrip/ast"
)

type expectedLetStatement struct {
	source 		string
	identifier 	string
	value		interface{}
}

type expectedReturnStatement struct {
	source 		string
	returnValue interface{}
}

// Test case
func TestLetStatement(t *testing.T) {
	Convey("Let statement testing", t, func() {
		expectedStatements := []expectedLetStatement{
			{ "let a = 5;",		"a",	5 },
			{ "let b = 5.1",	"b",	5.1 },
			{ "let c = true",	"c",	true },
			{ "let d = c",		"d",	"c" },
		}

		testLetStatement(expectedStatements)
	})
}

func TestBadLetStatement(t *testing.T) {
	Convey("Bad let statement testing", t, func() {
		sources := []string{ "let", "let x;" }

		for _, source := range sources {
			theLexer   := lexer.NewLexer(source)
			theParser  := NewParser(theLexer)
			theProgram := theParser.Parse()

			So(theProgram, ShouldNotBeNil)
			So(len(theParser.errors), ShouldBeGreaterThanOrEqualTo, 1)
			So(len(theParser.Errors()), ShouldEqual, len(theParser.errors))
		}
	})
}

func TestReturnStatement(t *testing.T) {
	Convey("Return statement testing", t, func() {
		expectedStatements := []expectedReturnStatement{
			{ "return 5;", 		5 },
			{ "return 10.1;", 	10.1 },
			{ "return true;", 	true },
			{ "return foo;", 	"foo" },
		}

		testReturnStatement(expectedStatements)
	})
}

func TestIdentifierExpression(t *testing.T) {
	Convey("Identifier expression test", t, func() {
		source := `foobar;`;

		theLexer   := lexer.NewLexer(source)
		theParser  := NewParser(theLexer)
		theProgram := theParser.Parse()

		testParserError(theParser)
		testParserProgramLength(theProgram, 1)

		Convey("can convert to expression statement", func() {
			statement, ok := theProgram.Statements[0].(*ast.ExpressionStatement)

			So(ok, ShouldBeTrue)

			Convey(`Check the value should be equal "foo"`, func() {
				identifier, ok := statement.Expression.(*ast.IdentifierExpression)

				So(ok, ShouldBeTrue)
				So(identifier.Value, ShouldEqual, "foobar")
				So(identifier.TokenLiteral(), ShouldEqual, "foobar")
			})
		})
	})
}

func TestIntegerLiteralExpression(t *testing.T) {
	Convey("Integer literal expression test", t, func() {
		source := `5;`

		theLexer   := lexer.NewLexer(source)
		theParser  := NewParser(theLexer)
		theProgram := theParser.Parse()

		testParserError(theParser)
		testParserProgramLength(theProgram, 1)

		Convey("Can convert to expression statement", func() {
			statement, ok := theProgram.Statements[0].(*ast.ExpressionStatement)

			So(ok, ShouldBeTrue)

			Convey("Check the value should be equal 5", func() {
				integer, ok := statement.Expression.(*ast.IntegerLiteralExpression)

				So(ok, ShouldBeTrue)
				So(integer.Value, ShouldEqual, 5)
				So(integer.TokenLiteral(), ShouldEqual, "5")
			})
		})
	})
}

func TestFloatLiteralExpression(t *testing.T) {
	Convey("Float literal expression test", t, func() {
		source := `12.34;`

		theLexer   := lexer.NewLexer(source)
		theParser  := NewParser(theLexer)
		theProgram := theParser.Parse()

		testParserError(theParser)
		testParserProgramLength(theProgram, 1)

		Convey("Can convert to expression statement", func() {
			statement, ok := theProgram.Statements[0].(*ast.ExpressionStatement)

			So(ok, ShouldBeTrue)

			Convey("Check the value should be equal 12.34", func() {
				float, ok := statement.Expression.(*ast.FloatLiteralExpression)

				So(ok, ShouldBeTrue)
				So(float.Value, ShouldEqual, 12.34)
				So(float.TokenLiteral(), ShouldEqual, "12.34")
			})
		})
	})
}

func TestBooleanExpression(t *testing.T) {
	Convey("Boolean expression test", t, func() {
		expectedExpressions := []struct{
			source string
			value  bool
		}{
			{ "true;",  true },
			{ "false;", false },
		}

		for index, expression := range expectedExpressions {
			message := runMessage("Running %d, Source: %s", index, expression.source)

			theLexer   := lexer.NewLexer(expression.source)
			theParser  := NewParser(theLexer)
			theProgram := theParser.Parse()

			Convey(message, func() {
				testParserError(theParser)
				testParserProgramLength(theProgram, 1)

				Convey("Can convert to expression statement", func() {
					statement, ok := theProgram.Statements[0].(*ast.ExpressionStatement)

					So(ok, ShouldBeTrue)

					Convey(
						runMessage("Check the value should be equal %s", strconv.FormatBool(expression.value)),
						func() {
							boolean := statement.Expression.(*ast.BooleanExpression)

							So(boolean.Value, ShouldEqual, expression.value)
						},
					)
				})
			})
		}
	})
}

func TestPrefixExpression(t *testing.T) {
	Convey("Prefix expression test", t, func() {
		expectedExpressions := []struct{
			source 		string
			operator 	string
			value 		interface{}
		}{
			{ "!5", 	 "!", 	5 },
			{ "-10", 	 "-", 	10 },
			{ "!foobar", "!",	"foobar" },
			{ "-foobar", "-", 	"foobar"} ,
			{ "!true",	 "!",	true },
			{ "!false",  "!",   false },
		}

		for index, expression := range expectedExpressions {
			message := runMessage("Running %d, Source: %s", index, expression.source)

			theLexer   := lexer.NewLexer(expression.source)
			theParser  := NewParser(theLexer)
			theProgram := theParser.Parse()

			Convey(message, func() {
				testParserError(theParser)
				testParserProgramLength(theProgram, 1)

				Convey("Can convert to expression statement", func() {
					statement, ok := theProgram.Statements[0].(*ast.ExpressionStatement)

					So(ok, ShouldBeTrue)

					Convey(
						runMessage("Check the operator should be equal %s", expression.operator),
						func() {
							prefix := statement.Expression.(*ast.PrefixExpression)

							So(prefix.Operator, ShouldEqual, expression.operator)

							testLiteralExpression(prefix.Right, expression.value)
						},
					)
				})
			})
		}
	})
}

func TestInfixExpression(t *testing.T) {
	Convey("Infix expression test", t, func() {
		expectedExpressions := []struct{
			source 		string
			leftValue 	interface{}
			operator 	string
			rightValue 	interface{}
		}{
			{ "10 + 10;", 	10, 	"+", 	10 },
			{ "11 - 11;", 	11, 	"-", 	11 },
			{ "12 * 12;", 	12, 	"*", 	12 },
			{ "13 / 13;", 	13, 	"/", 	13 },
			{ "14 > 14;", 	14, 	">", 	14 },
			{ "15 < 15;", 	15, 	"<", 	15 },
			{ "16 == 16;", 	16, 	"==", 	16 },
			{ "17 != 17;", 	17, 	"!=", 	17 },

			{ "foobar1 + barfoo1;", 	"foobar1", 	"+", 	"barfoo1" },
			{ "foobar2 - barfoo2;", 	"foobar2", 	"-", 	"barfoo2" },
			{ "foobar3 * barfoo3;", 	"foobar3", 	"*", 	"barfoo3" },
			{ "foobar4 / barfoo4;", 	"foobar4", 	"/", 	"barfoo4" },
			{ "foobar5 > barfoo5;", 	"foobar5", 	">", 	"barfoo5" },
			{ "foobar6 < barfoo6;", 	"foobar6", 	"<", 	"barfoo6" },
			{ "foobar7 == barfoo7;", 	"foobar7", 	"==", 	"barfoo7" },
			{ "foobar8 != barfoo8;", 	"foobar8", 	"!=", 	"barfoo8" },

			{ "true == true", 	true, 	"==", 	true },
			{ "true != false",	true, 	"!=", 	false },
			{ "false == false", false, 	"==", 	false },
		}

		for index, expression := range expectedExpressions {
			message := runMessage("Running %d, Source: %s", index, expression.source)

			theLexer   := lexer.NewLexer(expression.source)
			theParser  := NewParser(theLexer)
			theProgram := theParser.Parse()

			Convey(message, func() {
				Convey("Parse program check", func() {
					testParserError(theParser)
					testParserProgramLength(theProgram, 1)
				})

				statement, ok := theProgram.Statements[0].(*ast.ExpressionStatement)
				Convey("Can convert to expression statement", func() {
					So(ok, ShouldBeTrue)
				})

				testInfixExpression(statement.Expression, expression.leftValue, expression.operator, expression.rightValue)

				// infixExpression, ok := statement.Expression.(*ast.InfixExpression)
				// Convey("Can convert to infix expression", func() {
				// 	So(ok, ShouldBeTrue)
				// })

				// Convey("Left expression", func() {
				// 	testLiteralExpression(infixExpression.Left, expression.leftValue)
				// })

				// Convey("Operator check", func() {
				// 	So(infixExpression.Operator, ShouldEqual, expression.operator)
				// })

				// Convey("Right expression", func() {
				// 	testLiteralExpression(infixExpression.Right, expression.rightValue)
				// })
			})
		}
	})
}

func TestOperatorPrecedence(t *testing.T) {
	Convey("Operator precedence test", t, func() {
		expectedExpressions := []struct{
			source 		string
			expected 	string
			length 		int
		}{
			{ "-a * b", 					"((-a) * b)",								1 },
			{ "!-a", 						"(!(-a))",									1 },
			{ "a + b + c", 					"((a + b) + c)", 							1 },
			{ "a + b - c", 					"((a + b) - c)",							1 },
			{ "a * b * c", 					"((a * b) * c)", 							1 },
			{ "a * b / c", 					"((a * b) / c)",							1 },
			{ "a + b / c", 					"(a + (b / c))", 							1 },
			{ "a + b * c + d / e - f", 		"(((a + (b * c)) + (d / e)) - f)",			1 },
			{ "3 + 4; -5 * 5", 				"(3 + 4)((-5) * 5)",						2 },
			{ "5 > 4 == 3 < 4", 			"((5 > 4) == (3 < 4))",						1 },
			{ "5 < 4 != 3 > 4", 			"((5 < 4) != (3 > 4))",						1 },
			{ "3 + 4 * 5 == 3 * 1 + 4 * 5",	"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",	1 },

			{ "true", 	"true",		1 },
			{ "false", 	"false",	1 },

			{ "3 > 5 == false",	"((3 > 5) == false)",	1 },
			{ "3 < 5 == true", 	"((3 < 5) == true)",	1 },
		}

		for index, expression := range expectedExpressions {
			message := runMessage("Running %d, Source: %s", index, expression.source)

			theLexer   := lexer.NewLexer(expression.source)
			theParser  := NewParser(theLexer)
			theProgram := theParser.Parse()

			Convey(message, func() {
				Convey("Parse program check", func() {
					testParserError(theParser)
				})

				Convey(runMessage("Expected length: %d", expression.length), func() {
					testParserProgramLength(theProgram, expression.length)
				})

				Convey(runMessage("Expected: %s", expression.expected), func() {
					So(theProgram.String(), ShouldEqual, expression.expected)
				})
			})
		}
	})
}

func TestArrayLiteralExpression(t *testing.T) {
	Convey("Array literal expression", t, func() {
		source := `[1, 2 * 2, 3 + 3]`

		theLexer   := lexer.NewLexer(source)
		theParser  := NewParser(theLexer)
		theProgram := theParser.Parse()

		Convey("Parse program check", func() {
			testParserError(theParser)
			testParserProgramLength(theProgram, 1)
		})

		statement, ok := theProgram.Statements[0].(*ast.ExpressionStatement)
		Convey("Can convert to expression statement", func() {
			So(ok, ShouldBeTrue)
		})

		arrayLiteralExpression, ok := statement.Expression.(*ast.ArrayLiteralExpression)
		Convey("Can convert to array literal expression", func() {
			So(ok, ShouldBeTrue)
		})

		Convey("Array element should equals 3", func() {
			So(len(arrayLiteralExpression.Elements), ShouldEqual, 3)
		})

		Convey("Element 1 should equals 1", func() {
			testIntegerLiteralExpression(arrayLiteralExpression.Elements[0], 1)
		})

		Convey("Element 2 should equals 4", func() {
			testInfixExpression(arrayLiteralExpression.Elements[1], 2, "*", 2)
		})

		Convey("Element 3 should equals 6", func() {
			testInfixExpression(arrayLiteralExpression.Elements[2], 3, "+", 3)
		})
	})
}

// Sub method for test case
func testLetStatement(expectedStatements []expectedLetStatement) {
	for index, currentStatement := range expectedStatements {
		message := runMessage("Running %d, Source: %s", index, currentStatement.source)

		theLexer   := lexer.NewLexer(currentStatement.source)
		theParser  := NewParser(theLexer)
		theProgram := theParser.Parse()

		statement := theProgram.Statements[0]
		letStatement, ok := statement.(*ast.LetStatement)

		Convey(message, func() {
			testParserError(theParser)
			testParserProgramLength(theProgram, 1)

			// Identifier
			So(ok, ShouldBeTrue)
			So(statement.TokenLiteral(), ShouldEqual, "let")
			So(letStatement.Name.Value, ShouldEqual, currentStatement.identifier)
			So(letStatement.Name.TokenLiteral(), ShouldEqual, currentStatement.identifier)

			// Value
			testLiteralExpression(letStatement.Value, currentStatement.value)
		})
	}
}

func testReturnStatement(expectedStatements []expectedReturnStatement) {
	for index, currentStatement := range expectedStatements {
		message := runMessage("Running %d, Source: %s", index,currentStatement.source)

		theLexer   := lexer.NewLexer(currentStatement.source)
		theParser  := NewParser(theLexer)
		theProgram := theParser.Parse()

		returnStatement, ok := theProgram.Statements[0].(*ast.ReturnStatement)

		Convey(message, func() {
			testParserError(theParser)
			testParserProgramLength(theProgram, 1)

			// Return keywords
			So(ok, ShouldBeTrue)
			So(returnStatement.TokenLiteral(), ShouldEqual, "return")

			// Value
			testLiteralExpression(returnStatement.ReturnValue, currentStatement.returnValue)
		})
	}
}

// Sub method function for sub method
func testParserError(parser *Parser) {
	parserErrors       := parser.Errors()
	parserErrorsLength := len(parserErrors)

	So(parserErrorsLength, ShouldEqual, 0)
}

func testParserProgramLength(program *ast.Program, length int) {
	So(len(program.Statements), ShouldEqual, length)
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

func testInfixExpression(expression ast.Expression, leftValue interface{}, operator string, rightValue interface{}) {
	infixExpression, ok := expression.(*ast.InfixExpression)

	Convey("Can convert to infix expression", func() {
		So(ok, ShouldBeTrue)
	})

	Convey("Left expression", func() {
		testLiteralExpression(infixExpression.Left, leftValue)
	})

	Convey("Operator check", func() {
		So(infixExpression.Operator, ShouldEqual, operator)
	})

	Convey("Right expression", func() {
		testLiteralExpression(infixExpression.Right, rightValue)
	})
}

// Helper functions for common
func runMessage(format string, values ...interface{}) string {
	return fmt.Sprintf(format, values...)
}
