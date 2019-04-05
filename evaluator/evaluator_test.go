package evaluator

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/zeuxisoo/go-skrip/lexer"
	"github.com/zeuxisoo/go-skrip/parser"
	"github.com/zeuxisoo/go-skrip/object"
)

func TestIntegerLiteralExpression(t *testing.T) {
	Convey("Integer literal expression eval test", t, func() {
		expecteds := []struct{
			source string
			result int64
		}{
			{ "5",  5 },
			{ "10", 10 },
		}

		for index, expected := range expecteds {
			Convey(runMessage("Running: %d, ", index), func() {
				evaluated := testEval(expected.source)

				Convey(runMessage("Source: %s", expected.source), func() {
					testIntegerObject(evaluated, expected.result)
				})
			})
		}
	})
}

func TestFloatLiteralExpression(t *testing.T) {
	Convey("Float literal expression eval test", t, func() {
		expecteds := []struct{
			source string
			result float64
		}{
			{ "5.0",  5.0 },
			{ "10.3", 10.3 },
		}

		for index, expected := range expecteds {
			Convey(runMessage("Running: %d, ", index), func() {
				evaluated := testEval(expected.source)

				Convey(runMessage("Source: %s", expected.source), func() {
					testFloatObject(evaluated, expected.result)
				})
			})
		}
	})
}

func TestStringLiteralExpression(t *testing.T) {
	Convey("String literal expression eval test", t, func() {
		expecteds := []struct{
			source string
			result string
		}{
			{ `"foo"`,    "foo" },
			{ `"foobar"`, "foobar" },
		}

		for index, expected := range expecteds {
			Convey(runMessage("Running: %d, ", index), func() {
				evaluated := testEval(expected.source)

				Convey(runMessage("Source: %s", expected.source), func() {
					testStringObject(evaluated, expected.result)
				})
			})
		}
	})
}

func TestIdentifierExpression(t *testing.T) {
	Convey("Identifier expression test", t, func() {
		Convey("Error handling test", func() {
			expecteds := []struct{
				source string
				result string
			}{
				{ "foo",    "Identifier not found: foo" },
				{ "foobar", "Identifier not found: foobar" },
			}

			for index, expected := range expecteds {
				Convey(runMessage("Running: %d, ", index), func() {
					evaluated := testEval(expected.source)

					Convey(runMessage("Source: %s", expected.source), func() {
						testErrorObject(evaluated, expected.result)
					})
				})
			}
		})

		Convey("Get identifier from environment test", func() {
			expecteds := []struct{
				source string	// identifier
				value  object.Object
				result interface{}
			}{
				{ "foo",	&object.String{ Value: "fooString" },	"fooString" },
				{ "bar",	&object.Integer{ Value: 5 },			5 },
			}

			environment := object.NewEnvironment()
			for _, expected := range expecteds {
				environment.Set(expected.source, expected.value)
			}

			for index, expected := range expecteds {
				Convey(runMessage("Running: %d", index), func() {
					evaluated := testEvalWithEnv(expected.source, environment)

					Convey(runMessage("Source: %s, expected: %t", expected.source, expected.result), func() {
						testLiteralObject(evaluated, expected.result)
					})
				})
			}
		})

		Convey("Register and Get built-in function test", func() {
			RegisterBuiltIn(
				"fooFunction",
				func(environment *object.Environment, arguments ...object.Object) object.Object {
					return &object.String{
						Value: "foo function",
					}
				},
			)

			testBuiltInObject(testEval("fooFunction"), "foo function")
		})
	})
}

func TestBooleanExpression(t *testing.T) {
	Convey("Boolean expression test", t, func() {
		expecteds := []struct{
			source string
			result bool
		}{
			{ "true",  true },
			{ "false", false },
		}

		for index, expected := range expecteds {
			Convey(runMessage("Running: %d, ", index), func() {
				evaluated := testEval(expected.source)

				Convey(runMessage("Source: %s", expected.source), func() {
					testLiteralObject(evaluated, expected.result)
				})
			})
		}
	})
}

func TestReturnStatement(t *testing.T) {
	Convey("Return statement test", t, func() {
		expecteds := []struct{
			source string
			result interface{}
		}{
			{ "return 10", 10 },
			{ "return 15.5", 15.5 },
		}

		for index, expected := range expecteds {
			Convey(runMessage("Running: %d, ", index), func() {
				evaluated := testEval(expected.source)

				Convey(runMessage("Source: %s", expected.source), func() {
					testDecimalObject(evaluated, expected.result)
				})
			})
		}
	})
}

//
func testEval(source string) object.Object {
	return testEvalWithEnv(source, object.NewEnvironment())
}

func testEvalWithEnv(source string, environment *object.Environment) object.Object {
	theLexer       := lexer.NewLexer(source)
	theParser      := parser.NewParser(theLexer)
	theProgarm 	   := theParser.Parse()

	return Eval(theProgarm, environment)
}

func testLiteralObject(obj object.Object, expected interface{}) {
	switch expected := expected.(type) {
	case int:
		testIntegerObject(obj, int64(expected))
	case int64:
		testIntegerObject(obj, expected)
	case string:
		testStringObject(obj, expected)
	case bool:
		testBooleanObject(obj, expected)
	case float32:
		testFloatObject(obj, float64(expected))
	case float64:
		testFloatObject(obj, expected)
	}
}

func testDecimalObject(obj object.Object, expected interface{}) {
	testLiteralObject(obj, expected)
}

func testIntegerObject(obj object.Object, expected int64) {
	result, ok := obj.(*object.Integer)

	Convey("Can convert to object (integer)", func() {
		So(ok, ShouldBeTrue)
	})

	Convey(runMessage("Object result should be equals %d", expected), func() {
		So(result.Value, ShouldEqual, expected)
	})
}

func testFloatObject(obj object.Object, expected float64) {
	result, ok := obj.(*object.Float)

	Convey("Can convert to object (float)", func() {
		So(ok, ShouldBeTrue)
	})

	Convey(runMessage("Object result should be equals %f", expected), func() {
		So(result.Value, ShouldEqual, expected)
	})
}

func testStringObject(obj object.Object, expected string) {
	result, ok := obj.(*object.String)

	Convey("Can convert to object (string)", func() {
		So(ok, ShouldBeTrue)
	})

	Convey(runMessage("Object result should be equals %s", expected), func() {
		So(result.Value, ShouldEqual, expected)
	})
}

func testBooleanObject(obj object.Object, expected bool) {
	result, ok := obj.(*object.Boolean)

	Convey("Can convert to object (boolean)", func() {
		So(ok, ShouldBeTrue)
	})

	Convey(runMessage("Object result should be equals %t", expected), func() {
		So(result.Value, ShouldEqual, expected)
	})
}

func testBuiltInObject(obj object.Object, expected string) {
	result, ok := obj.(*object.BuiltIn)

	Convey("Can convert to object (built-in)", func() {
		So(ok, ShouldBeTrue)
	})

	testLiteralObject(result.Function(object.NewEnvironment()), expected)
}

func testErrorObject(obj object.Object, expected string) {
	result, ok := obj.(*object.Error)

	Convey("Can convert to object (error)", func() {
		So(ok, ShouldBeTrue)
	})

	Convey(runMessage("Object result should be equals %s", expected), func() {
		So(result.Message, ShouldEqual, expected)
	})
}

// Helper functions for common
func runMessage(format string, values ...interface{}) string {
	return fmt.Sprintf(format, values...)
}
