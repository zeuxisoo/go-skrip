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

func TestFunctionStatement(t *testing.T) {
	Convey("Function statement testing", t, func() {
		source := `func funcName(x, y) { x + y; };`

		theLexer   := lexer.NewLexer(source)
		theParser  := NewParser(theLexer)
		theProgram := theParser.Parse()

		Convey("Parse program check", func() {
			testParserError(theParser)
			testParserProgramLength(theProgram, 1)
		})

		functionStatement, ok := theProgram.Statements[0].(*ast.FunctionStatement)
		Convey("Can convert to function statement", func() {
			So(ok, ShouldBeTrue)
		})

		Convey("Function name should equals funcName", func() {
			So(functionStatement.Name.String(), ShouldEqual, "funcName")
		})

		Convey("Function parameter length should equals 2", func() {
			So(len(functionStatement.Function.Parameters), ShouldEqual, 2)
		})

		Convey("Function Parameter names should equal x and y", func() {
			testLiteralExpression(functionStatement.Function.Parameters[0], "x")
			testLiteralExpression(functionStatement.Function.Parameters[1], "y")
		})

		functionBlockStatement, ok := functionStatement.Function.Block.Statements[0].(*ast.ExpressionStatement)
		Convey("Can convert function block to expression statement", func() {
			So(ok, ShouldBeTrue)
		})

		Convey("Function block should be equals x+y", func() {
			testInfixExpression(functionBlockStatement.Expression, "x", "+", "y")
		})
	})
}

func TestIntegerLiteralExpression(t *testing.T) {
	Convey("Integer literal expression test", t, func() {
		source := `5;`

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

		integerLiteralExpression, ok := statement.Expression.(*ast.IntegerLiteralExpression)
		Convey("Can convert to integer literal expression", func() {
			So(ok, ShouldBeTrue)
		})

		Convey("Integer float expression value should be equal 5", func() {
			So(integerLiteralExpression.Value, ShouldEqual, 5)
			So(integerLiteralExpression.TokenLiteral(), ShouldEqual, "5")
		})
	})
}

func TestFloatLiteralExpression(t *testing.T) {
	Convey("Float literal expression test", t, func() {
		source := `12.34;`

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

		float, ok := statement.Expression.(*ast.FloatLiteralExpression)
		Convey("Can convert to float literal expression", func() {
			So(ok, ShouldBeTrue)
		})

		Convey("Float literal expression value should be equal 12.34", func() {
			So(float.Value, ShouldEqual, 12.34)
			So(float.TokenLiteral(), ShouldEqual, "12.34")
		})
	})
}

func TestStringLiteralExpression(t *testing.T) {
	Convey("String literal expression test", t, func() {
		source := `"Hello World";`

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

		stringLiteralExpression, ok := statement.Expression.(*ast.StringLiteralExpression)
		Convey("Can convert to string literal expression", func() {
			So(ok, ShouldBeTrue)
		})

		Convey(`String literal expression value should be equal "Hello World"`, func() {
			So(stringLiteralExpression.Value, ShouldEqual, "Hello World")
			So(stringLiteralExpression.TokenLiteral(), ShouldEqual, "Hello World")
		})
	})
}

func TestLetStatementFunctionLiteralExpression(t *testing.T) {
	Convey("Let statement function literal expression test", t, func() {
		source := `let foo = func(x, y) { x + y; }`

		theLexer   := lexer.NewLexer(source)
		theParser  := NewParser(theLexer)
		theProgram := theParser.Parse()

		Convey("Parse program check", func() {
			testParserError(theParser)
			testParserProgramLength(theProgram, 1)
		})

		letStatement, ok := theProgram.Statements[0].(*ast.LetStatement)
		Convey("Can convert to let statement", func() {
			So(ok, ShouldBeTrue)
		})

		Convey("Let function name should be foo", func() {
			So(letStatement.Name.String(), ShouldEqual, "foo")
		})

		functionLiteralExpression, ok := letStatement.Value.(*ast.FunctionLiteralExpression)
		Convey("Can convert to function literal expression", func() {
			So(ok, ShouldBeTrue)
		})

		Convey("Function parameter length should be equal 2", func() {
			So(len(functionLiteralExpression.Parameters), ShouldEqual, 2)
		})

		Convey("Function parameter should be x and y", func() {
			testLiteralExpression(functionLiteralExpression.Parameters[0], "x")
			testLiteralExpression(functionLiteralExpression.Parameters[1], "y")
		})

		Convey("Function body statement length should be equal 1", func() {
			So(len(functionLiteralExpression.Block.Statements), ShouldEqual, 1)
		})

		functionBlockStatement, ok := functionLiteralExpression.Block.Statements[0].(*ast.ExpressionStatement)
		Convey("Can convert function block to expression statement", func() {
			So(ok, ShouldBeTrue)
		})

		Convey("Function block should equals x+y", func() {
			testInfixExpression(functionBlockStatement.Expression, "x", "+", "y")
		})
	})
}

func TestLetStatementFunctionParameterParsing(t *testing.T) {
	Convey("Let statement function parameter parsing test", t, func() {
		expectedStatements := []struct{
			source 		string
			parameters 	[]string
		}{
			{ "let foo = func() {};", 			[]string{} },
			{ "let foo = func(x) {};", 			[]string{ "x" } },
			{ "let foo = func(x, y, z) {};",  	[]string{ "x", "y", "z" } },
		}

		for index, expected := range expectedStatements {
			message := runMessage("Running %d, Source: %s", index, expected.source)

			theLexer   := lexer.NewLexer(expected.source)
			theParser  := NewParser(theLexer)
			theProgram := theParser.Parse()

			Convey(message, func() {
				Convey("Parse program check", func() {
					testParserError(theParser)
					testParserProgramLength(theProgram, 1)
				})

				letStatement, ok := theProgram.Statements[0].(*ast.LetStatement)
				Convey("Can convert to expression statement", func() {
					So(ok, ShouldBeTrue)
				})

				Convey("Let function name should be foo", func() {
					So(letStatement.Name.String(), ShouldEqual, "foo")
				})

				functionLiteralExpression, ok := letStatement.Value.(*ast.FunctionLiteralExpression)
				Convey("Can convert to function literal expression", func() {
					So(ok, ShouldBeTrue)
				})

				expectedFunctionParameterLength := len(expected.parameters)
				Convey(runMessage(
					"Function parameter length should be equals %d",
					expectedFunctionParameterLength,
				), func() {
					So(len(functionLiteralExpression.Parameters), ShouldEqual, expectedFunctionParameterLength)
				})

				for index2, parameter := range expected.parameters {
					Convey(runMessage(
						"Running: %d, Expected paramter: %s",
						index2, parameter,
					), func() {
						testLiteralExpression(functionLiteralExpression.Parameters[index2], parameter)
					})
				}
			})
		}
	})
}

func TestNoNamedFunctionLiteralExpression(t *testing.T) {
	Convey("No named function literal expression test", t, func() {
		source := `func(x, y) { x + y; }`

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

		functionLiteralExpression, ok := statement.Expression.(*ast.FunctionLiteralExpression)
		Convey("Can convert to function literal expression", func() {
			So(ok, ShouldBeTrue)
		})

		Convey("Function parameter length should be equal 2", func() {
			So(len(functionLiteralExpression.Parameters), ShouldEqual, 2)
		})

		Convey("Function parameter should be x and y", func() {
			testLiteralExpression(functionLiteralExpression.Parameters[0], "x")
			testLiteralExpression(functionLiteralExpression.Parameters[1], "y")
		})

		Convey("Function body statement length should be equal 1", func() {
			So(len(functionLiteralExpression.Block.Statements), ShouldEqual, 1)
		})

		functionBlockStatement, ok := functionLiteralExpression.Block.Statements[0].(*ast.ExpressionStatement)
		Convey("Can convert function block to expression statement", func() {
			So(ok, ShouldBeTrue)
		})

		Convey("Function block should equals x+y", func() {
			testInfixExpression(functionBlockStatement.Expression, "x", "+", "y")
		})
	})
}

func TestNoNamedFunctionParameterParsing(t *testing.T) {
	Convey("No named function parameter parsing test", t, func() {
		expectedStatements := []struct{
			source 		string
			parameters 	[]string
		}{
			{ "func() {};", 		[]string{} },
			{ "func(x) {};", 		[]string{ "x" } },
			{ "func(x, y, z) {};",  []string{ "x", "y", "z" } },
		}

		for index, expected := range expectedStatements {
			message := runMessage("Running %d, Source: %s", index, expected.source)

			theLexer   := lexer.NewLexer(expected.source)
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

				functionLiteralExpression, ok := statement.Expression.(*ast.FunctionLiteralExpression)
				Convey("Can convert to function literal expression", func() {
					So(ok, ShouldBeTrue)
				})

				expectedFunctionParameterLength := len(expected.parameters)
				Convey(runMessage(
					"Function parameter length should be equals %d",
					expectedFunctionParameterLength,
				), func() {
					So(len(functionLiteralExpression.Parameters), ShouldEqual, expectedFunctionParameterLength)
				})

				for index2, parameter := range expected.parameters {
					Convey(runMessage(
						"Running: %d, Expected paramter: %s",
						index2, parameter,
					), func() {
						testLiteralExpression(functionLiteralExpression.Parameters[index2], parameter)
					})
				}
			})
		}
	})
}

func TestIdentifierExpression(t *testing.T) {
	Convey("Identifier expression test", t, func() {
		source := `foobar;`;

		theLexer   := lexer.NewLexer(source)
		theParser  := NewParser(theLexer)
		theProgram := theParser.Parse()

		Convey("Parse program check", func() {
			testParserError(theParser)
			testParserProgramLength(theProgram, 1)
		})

		statement, ok := theProgram.Statements[0].(*ast.ExpressionStatement)
		Convey("can convert to expression statement", func() {
			So(ok, ShouldBeTrue)
		})

		identifierExpression, ok := statement.Expression.(*ast.IdentifierExpression)
		Convey("can convert to identifer expression", func() {
			So(ok, ShouldBeTrue)
		})

		Convey(`Identifer expression value should be equal "foobar"`, func() {
			So(identifierExpression.Value, ShouldEqual, "foobar")
			So(identifierExpression.TokenLiteral(), ShouldEqual, "foobar")
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
				Convey("Parse program check", func() {
					testParserError(theParser)
					testParserProgramLength(theProgram, 1)
				})

				statement, ok := theProgram.Statements[0].(*ast.ExpressionStatement)
				Convey("Can convert to expression statement", func() {
					So(ok, ShouldBeTrue)
				})

				booleanExpression, ok := statement.Expression.(*ast.BooleanExpression)
				Convey("Can convert to boolean expression", func() {
					So(ok, ShouldBeTrue)
				})

				Convey(runMessage(
					"Boolean expression value should be equal %s",
					strconv.FormatBool(expression.value),
				), func() {
						So(booleanExpression.Value, ShouldEqual, expression.value)
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
				Convey("Parse program check", func() {
					testParserError(theParser)
					testParserProgramLength(theProgram, 1)
				})

				statement, ok := theProgram.Statements[0].(*ast.ExpressionStatement)
				Convey("Can convert to expression statement", func() {
					So(ok, ShouldBeTrue)
				})

				prefixExpression := statement.Expression.(*ast.PrefixExpression)
				Convey("Can convert to prefix expression", func() {
					So(ok, ShouldBeTrue)
				})

				Convey(runMessage(
					"Check the operator should be equal %s",
					expression.operator,
				), func() {
					So(prefixExpression.Operator, ShouldEqual, expression.operator)

					testLiteralExpression(prefixExpression.Right, expression.value)
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

			{ "foo && bar", 	"foo",	"&&",	"bar" },
			{ "foo || bar", 	"foo",	"||",	"bar" },
			{ "true && true",  	true,	"&&",	true},
			{ "true && false", 	true,	"&&",	false },

			{ "1..3", 			1,		"..",	3 },
			{ "foo..bar",		"foo",	"..",	"bar" },
			{ "true..false",	true,	"..",	false },
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
	Convey("Array literal expression test", t, func() {
		source := `[1, 2 * 2, 3 + 3];`

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

		Convey("Element 1, integer literal should equals 1", func() {
			testIntegerLiteralExpression(arrayLiteralExpression.Elements[0], 1)
		})

		Convey("Element 2, infix expression should equals 4", func() {
			testInfixExpression(arrayLiteralExpression.Elements[1], 2, "*", 2)
		})

		Convey("Element 3, infix expression should equals 6", func() {
			testInfixExpression(arrayLiteralExpression.Elements[2], 3, "+", 3)
		})
	})
}

func TestEmptyArrayLieteralExpression(t *testing.T) {
	Convey("Empty array literal expression test", t, func() {
		source := `[];`

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

		Convey("Array element should equals 0", func() {
			So(len(arrayLiteralExpression.Elements), ShouldEqual, 0)
		})
	})
}

func TestHashLiteralExpressionStringKeys(t *testing.T) {
	Convey("Hash literal expression string keys test", t, func() {
		source   := `{ "one": 1, "two": 2, "three":3 };`
		expected := map[string]int{
			"one"  : 1,
			"two"  : 2,
			"three": 3,
		}

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

		hashLiteralExpression, ok := statement.Expression.(*ast.HashLiteralExpression)
		Convey("Can convert to hash literal expression", func() {
			So(ok, ShouldBeTrue)
		})

		Convey("Hash pairs length should equals expected pairs length", func() {
			So(len(hashLiteralExpression.Pairs), ShouldEqual, len(expected))
		})

		Convey("Hash values should matched", func() {
			for keyExpression, valueExpression := range hashLiteralExpression.Pairs {
				keyString     := keyExpression.(*ast.StringLiteralExpression)
				expectedValue := expected[keyString.String()]

				testIntegerLiteralExpression(valueExpression, int64(expectedValue))
			}
		})
	})
}

func TestHashLiteralExpressionBooleanKeys(t *testing.T) {
	Convey("Hash literal expression boolean keys test", t, func() {
		source   := `{ true: 1, false: 2 };`
		expected := map[string]int{
			"true" : 1,
			"false": 2,
		}

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

		hashLiteralExpression, ok := statement.Expression.(*ast.HashLiteralExpression)
		Convey("Can convert to hash literal expression", func() {
			So(ok, ShouldBeTrue)
		})

		Convey("Hash pairs length should equals expected pairs length", func() {
			So(len(hashLiteralExpression.Pairs), ShouldEqual, len(expected))
		})

		Convey("Hash values should matched", func() {
			for keyExpression, valueExpression := range hashLiteralExpression.Pairs {
				keyString     := keyExpression.(*ast.BooleanExpression)
				expectedValue := expected[keyString.String()]

				testIntegerLiteralExpression(valueExpression, int64(expectedValue))
			}
		})
	})
}

func TestHashLiteralExpressionIntegerKeys(t *testing.T) {
	Convey("Hash literal expression integer keys test", t, func() {
		source   := `{ 1: 1, 2: 2, 3: 3 };`
		expected := map[string]int{
			"1": 1,
			"2": 2,
			"3": 3,
		}

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

		hashLiteralExpression, ok := statement.Expression.(*ast.HashLiteralExpression)
		Convey("Can convert to hash literal expression", func() {
			So(ok, ShouldBeTrue)
		})

		Convey("Hash pairs length should equals expected pairs length", func() {
			So(len(hashLiteralExpression.Pairs), ShouldEqual, len(expected))
		})

		Convey("Hash values should matched", func() {
			for keyExpression, valueExpression := range hashLiteralExpression.Pairs {
				keyString     := keyExpression.(*ast.IntegerLiteralExpression)
				expectedValue := expected[keyString.String()]

				testIntegerLiteralExpression(valueExpression, int64(expectedValue))
			}
		})
	})
}

func TestHashLiteralExpressionWithExpressionValues(t *testing.T) {
	Convey("Hash literal expression with expression values test", t, func() {
		type expectedValue struct {
			left 		interface{}
			operator  	string
			right 		interface{}
		}

		source      := `{ "one": 1 + 2, "two": 10 - 7, "three": 15 / 3 };`
		expectedMap := map[string]expectedValue{
			"one"  : expectedValue{ 1, "+", 2 },
			"two"  : expectedValue{ 10, "-", 7 },
			"three": expectedValue{ 15, "/", 3 },
		}

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

		hashLiteralExpression, ok := statement.Expression.(*ast.HashLiteralExpression)
		Convey("Can convert to hash literal expression", func() {
			So(ok, ShouldBeTrue)
		})

		Convey("Hash pairs length should equals expected pairs length", func() {
			So(len(hashLiteralExpression.Pairs), ShouldEqual, len(expectedMap))
		})

		Convey("Hash pairs expected test", func() {
			// Convert the hash[Expression]Expression to hash[string]Expression
			hashWithStringKeys := make(map[string]ast.Expression)
			for keyExpression, valueExpression := range hashLiteralExpression.Pairs {
				keyStringExpression, _ := keyExpression.(*ast.StringLiteralExpression)

				hashWithStringKeys[keyStringExpression.String()] = valueExpression
			}

			// And then loop expectedMap not hashLiteralExpression.Pairs
			// to fix the for loop(hashLiteralExpression.Pairs) inside convey will make the test message/result sorting incorrect
			for key, value := range expectedMap {
				Convey(runMessage("Expected test case: %s", key), func() {
					hashValueExpression := hashWithStringKeys[key]

					testInfixExpression(hashValueExpression, value.left, value.operator, value.right)
				})
			}
		})
	})
}

func TestEmptyHashLiteralExpression(t *testing.T) {
	Convey("Empty hash literal expression test", t, func() {
		source := `{};`

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

		hashLiteralExpression, ok := statement.Expression.(*ast.HashLiteralExpression)
		Convey("Can convert to hash literal expression", func() {
			So(ok, ShouldBeTrue)
		})

		Convey("Hash pairs length should equals 0", func() {
			So(len(hashLiteralExpression.Pairs), ShouldEqual, 0)
		})
	})
}

func TestIndexExpression(t *testing.T) {
	Convey("Index expression test", t, func() {
		source := `myArray[1+2];`

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

		indexExpression, ok := statement.Expression.(*ast.IndexExpression)
		Convey("Can convert to index expression", func() {
			So(ok, ShouldBeTrue)
		})

		Convey("Left expression should equals myArray", func() {
			testIdentifierExpression(indexExpression.Left, "myArray")
		})

		Convey("Index expression should equals [1+2]", func() {
			testInfixExpression(indexExpression.Index, 1, "+", 2)
		})
	})
}

func TestGroupedExpression(t *testing.T) {
	Convey("Grouped expression test", t, func() {
		expectedExpressions := []struct{
			source   string
			expected string
		}{
			{ "1 + (2 + 3) + 4", 	"((1 + (2 + 3)) + 4)" },
			{ "(5 + 5) * 2", 		"((5 + 5) * 2)" },
			{ "2 / (5 + 5)",		"(2 / (5 + 5))" },
			{ "-(5 + 5)", 			"(-(5 + 5))" },
			{ "!(true == true)", 	"(!(true == true))" },
		}

		for index, expression := range expectedExpressions {
			Convey(runMessage("Running: %d, Source: %s", index, expression.source), func() {
				theLexer   := lexer.NewLexer(expression.source)
				theParser  := NewParser(theLexer)
				theProgram := theParser.Parse()

				Convey("Parse program check", func() {
					testParserError(theParser)
					testParserProgramLength(theProgram, 1)
				})

				Convey(runMessage("Expected: %s", expression.expected), func() {
					So(theProgram.String(), ShouldEqual, expression.expected)
				})
			})
		}
	})
}

func TestCallExpression(t *testing.T) {
	Convey("Call expression test", t, func() {
		expectedExpressions := []struct{
			source   string
			expected string
		}{
			{ "add(1, 2 * 3, 4 + 5)",				"add(1, (2 * 3), (4 + 5))" },
			{ "a + add(b * c) + d", 				"((a + add((b * c))) + d)"},
			{ "add(a * b[4], b[5], 2 * [6, 7][1])", "add((a * (b[4])), (b[5]), (2 * ([6, 7][1])))"},
		}

		for index, expression := range expectedExpressions {
			Convey(runMessage("Running: %d, Source: %s", index, expression.source), func() {
				theLexer   := lexer.NewLexer(expression.source)
				theParser  := NewParser(theLexer)
				theProgram := theParser.Parse()

				Convey("Parse program check", func() {
					testParserError(theParser)
					testParserProgramLength(theProgram, 1)
				})

				Convey(runMessage("Expected: %s", expression.expected), func() {
					So(theProgram.String(), ShouldEqual, expression.expected)
				})
			})
		}
	})
}

func TestIfExpression(t *testing.T) {
	Convey("If expression test", t, func() {
		source := `if (a < b) { c };`

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

		ifExpression, ok := statement.Expression.(*ast.IfExpression)
		Convey("Can convert to if expression", func() {
			So(ok, ShouldBeTrue)
		})

		Convey("If condition test", func() {
			testInfixExpression(ifExpression.Condition, "a", "<", "b")
		})

		Convey("If condition block length should equals 1", func() {
			So(len(ifExpression.Block.Statements), ShouldEqual, 1)
		})

		block, ok := ifExpression.Block.Statements[0].(*ast.ExpressionStatement)
		Convey("Can convert if condition block to expression statement", func() {
			So(ok, ShouldBeTrue)
		})

		Convey("Identifier should be named c", func() {
			testIdentifierExpression(block.Expression, "c")
		})

		Convey("Else alternative block should be nil", func() {
			So(ifExpression.Alternative, ShouldBeNil)
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
