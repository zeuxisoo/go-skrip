package evaluator

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/zeuxisoo/go-skrip/lexer"
	"github.com/zeuxisoo/go-skrip/parser"
	"github.com/zeuxisoo/go-skrip/object"
)

//
type expectedFunctions struct {
	source          string
	parameterLength int
	blockLength     int
}

//
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

func TestArrayLiteralExpression(t *testing.T) {
	Convey("Array literal expression test", t, func() {
		expecteds := []struct{
			source   string
			length   int
			elements []string
		}{
			{ `[1, 2, 3]`,         3, []string{ "1", "2", "3" } },
			{ `[5.1, 6.2, 7.3]`,   3, []string{ "5.1", "6.2", "7.3" } },
			{ `["a", "b", "c"]`,   3, []string{ "a", "b", "c" } },
			{ `[5.1, "a", 2, 1]`,  4, []string{ "5.1", "a", "2", "1" } },
		}

		for index, expected := range expecteds {
			Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
				evaluated := testEval(expected.source)

				array, ok := evaluated.(*object.Array)
				Convey("Can convert to object (array)", func() {
					So(ok, ShouldBeTrue)
				})

				Convey(runMessage("Elements length should equals %d", expected.length), func() {
					So(len(array.Elements), ShouldEqual, expected.length)
				})

				//
				compareElements := []string{}
				for _, element := range array.Elements {
					compareElements = append(compareElements, element.Inspect())
				}

				Convey(runMessage(`Elements should equals %s`, expected.elements), func() {
					So(compareElements, ShouldResemble, expected.elements)
				})
			})
		}
	})
}

func TestHashLiteralExpression(t *testing.T) {
	Convey("Hash literal expression test", t, func() {
		expecteds := []struct{
			source string
			length int
			order  []string
		}{
			{ `{ "foo": 1, "bar": 2 }`,        2, []string{ "foo:1", "bar:2" } },
			{ `{ 1: "foo", 2: "bar" }`,        2, []string{ "1:foo", "2:bar" } },
			{ `{ 5.5: "foo", 6.6: "bar" }`,    2, []string{ "5.5:foo", "6.6:bar" } },
			{ `{ true: "foo", false: "bar" }`, 2, []string{ "true:foo", "false:bar" } },

			{ `{ "z": 10, "d": 20, "a": 1 }`,           3, []string{ "z:10", "d:20", "a:1" } },
			{ `{ 20: "c", 10: "h", 30: "e", 12: "d" }`, 4, []string{ "20:c", "10:h", "30:e", "12:d" } },
			{ `{ "k": 1, 2.2: "g", 1: "5", "e": "9" }`, 4, []string{ "k:1", "2.2:g", "1:5", "e:9" } },
		}

		for index, expected := range expecteds {
			Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
				evaluated := testEval(expected.source)

				hash, ok := evaluated.(*object.Hash)
				Convey("Can convert to object (hash)", func() {
					So(ok, ShouldBeTrue)
				})

				Convey(runMessage("Order (keys) length should equals %d", expected.length), func() {
					So(len(hash.Order), ShouldEqual, expected.length)
				})

				Convey(runMessage("Pairs length should equals %d", expected.length), func() {
					So(len(hash.Pairs), ShouldEqual, expected.length)
				})

				//
				compareOrders := make([]string, 0)
				for _, key := range hash.Order {
					pair := hash.Pairs[key]

					pairValue := fmt.Sprintf("%s:%s", pair.Key.Inspect(), pair.Value.Inspect())
					Convey(runMessage(`Pair "%s" should be in %s`, pairValue, expected.order), func() {
						So(pairValue, ShouldBeIn, expected.order)
					})

					compareOrders = append(compareOrders, pairValue)
				}

				Convey(runMessage(`Order should equals %s`, expected.order), func() {
					So(compareOrders, ShouldResemble, expected.order)
				})
			})
		}
	})
}

func TestFunctionLiteralExpression(t *testing.T) {
	Convey("Function literal expression test", t, func() {
		expecteds := []struct{
			source          string
			parameterLength int
			blockLength     int
		}{
			{ "func(a, b, c) { d }", 3, 1 },
			{ "func(a, b) { c; d }", 2, 2 },
		}

		for index, expected := range expecteds {
			Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
				evaluated := testEval(expected.source)

				testFunctionObject(evaluated, expected)
			})
		}
	})
}

func TestCallExpression(t *testing.T) {
	Convey("Call expression test", t, func() {
		expecteds := []struct{
			source string
			result interface{}
		}{
			{ `func a() { return 123; }; a();`,   123 },
			{ `func a() { return 12.3; }; a();`,  12.3 },
			{ `func a() { return "123"; }; a();`, "123" },

			{ `func a(b) { return b; }; a("foo");`, "foo" },
			{ `func a(b, c, d) { return d; }; a("foo", 123, 4.5);`, 4.5 },

			{ `func a() { let b = "foo"; return b; }; a();`, "foo" },
		}

		for index, expected := range expecteds {
			Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
				evaluated := testEval(expected.source)

				testLiteralObject(evaluated, expected.result)
			})
		}
	})
}

func TestIndexExpression(t *testing.T) {
	Convey("Index expression test", t, func() {
		Convey("For array object", func() {
			expecteds := []struct{
				source string
				result interface{}
			}{
				{ "[1, 2, 3][2]",       3 },
				{ `[1.1, 2.2, 3.3][0]`, 1.1 },
				{ `["a", "b", "c"][1]`, "b" },
			}

			for index, expected := range expecteds {
				Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
					evaluated := testEval(expected.source)

					testLiteralObject(evaluated, expected.result)
				})
			}
		})

		Convey("For hash object", func() {
			expecteds := []struct{
				source string
				result interface{}
			}{
				{ `{ 1: "a", "2": 7.2, 3.1: 50 }[1]`,   "a"},
				{ `{ 1: "a", "2": 7.2, 3.1: 50 }["2"]`, 7.2},
				{ `{ 1: "a", "2": 7.2, 3.1: 50 }[3.1]`, 50},
			}

			for index, expected := range expecteds {
				Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
					evaluated := testEval(expected.source)

					testLiteralObject(evaluated, expected.result)
				})
			}
		})
	})
}

// Statements
func TestLetStatement(t *testing.T) {
	Convey("Let statement test", t, func() {
		expecteds := []struct{
			source string
			value  interface{}
		}{
			{ `let a = 5;`,     5 },
			{ `let b = 5.5;`,   5.5 },
			{ `let c = "foo";`, "foo" },
		}

		for index, expected := range expecteds {
			Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
				evaluated := testEval(expected.source)

				testLiteralObject(evaluated, expected.value)
			})
		}
	})
}

func TestLetStatementWithFunctionLiteralExpression(t *testing.T) {
	Convey("Let statement with function literal expression", t, func() {
		source := "let a = func(a, b) { c };"

		// Should be return function object
		evaluated := testEval(source)

		testFunctionObject(evaluated, expectedFunctions{
			parameterLength: 2,
			blockLength    : 1,
		})
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

func TestFunctionStatement(t *testing.T) {
	Convey("Function statement test", t, func() {
		expecteds := []expectedFunctions{
			{ "func myFunc1(a, b, c) { d }", 3, 1 },
			{ "func myFunc2(a, b) { c; d }", 2, 2 },
		}

		for index, expected := range expecteds {
			Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
				evaluated := testEval(expected.source)

				testFunctionObject(evaluated, expected)
			})
		}
	})
}

func TestBlockStatement(t *testing.T) {
	Convey("Block statement test", t, func() {
		expecteds := []struct{
			source      string
			returnValue interface{}
		}{
			{ "let a = 1;",            1 },
			{ "let a = 1; let b = 2;", 2 },

			{ "let a = 1.1;",              1.1 },
			{ "let a = 1.1; let b = 2.2;", 2.2 },

			{ `let a = "foo";`,               "foo" },
			{ `let a = "foo"; let b = "bar"`, "bar" },

			{ `let a = "";`,                  ""},
			{ `let a = "foobar"; return a;`,  "foobar" },
		}

		for index, expected := range expecteds {
			Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
				evaluated := testEval(expected.source)

				testLiteralObject(evaluated, expected.returnValue)
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

func testFunctionObject(obj object.Object, expected expectedFunctions) {
	function, ok := obj.(*object.Function)
	Convey("Can convert to object (function)", func() {
		So(ok, ShouldBeTrue)
	})

	Convey(runMessage("Function parameters length should be equals %d", expected.parameterLength), func() {
		So(len(function.Parameters), ShouldEqual, expected.parameterLength)
	})

	Convey(runMessage("Function block should be equals %d", expected.blockLength), func() {
		So(len(function.Block.Statements), ShouldEqual, expected.blockLength)
	})
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
