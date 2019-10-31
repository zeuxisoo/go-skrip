package evaluator

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/zeuxisoo/go-skrip/lexer"
	"github.com/zeuxisoo/go-skrip/object"
	"github.com/zeuxisoo/go-skrip/parser"
)

//
type expectedFunction struct {
	source          string
	parameterLength int
	blockLength     int
}

type expectedArray struct {
	source   string
	length   int
	elements []string
}

type expectedHash struct {
	source string
	length int
	order  []string
}

//
func TestIntegerLiteralExpression(t *testing.T) {
	Convey("Integer literal expression test", t, func() {
		expecteds := []struct {
			source string
			result int64
		}{
			{"5", 5},
			{"10", 10},
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
	Convey("Float literal expression test", t, func() {
		expecteds := []struct {
			source string
			result float64
		}{
			{"5.0", 5.0},
			{"10.3", 10.3},
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
	Convey("String literal expression test", t, func() {
		expecteds := []struct {
			source string
			result string
		}{
			{`"foo"`, "foo"},
			{`"foobar"`, "foobar"},
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

func TestNilLiteralExpression(t *testing.T) {
	Convey("Nil literal expression test", t, func() {
		expected := struct {
			source string
			result string
		}{
			`nil`, "nil",
		}

		evaluated := testEval(expected.source)

		testNilObject(evaluated, expected.result)
	})
}

func TestIdentifierExpression(t *testing.T) {
	Convey("Identifier expression test", t, func() {
		Convey("Error handling test", func() {
			expecteds := []struct {
				source string
				result string
			}{
				{"foo", "Identifier not found: foo"},
				{"foobar", "Identifier not found: foobar"},
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
			expecteds := []struct {
				source string // identifier
				value  object.Object
				result interface{}
			}{
				{"foo", &object.String{Value: "fooString"}, "fooString"},
				{"bar", &object.Integer{Value: 5}, 5},
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
		expecteds := []struct {
			source string
			result bool
		}{
			{"true", true},
			{"false", false},
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

func TestAssignExpression(t *testing.T) {
	Convey("Assign expression test", t, func() {
		expecteds := []struct {
			source string
			result interface{}
		}{
			{`let a = "foo"; a = "bar"; a`, "bar"},
			{`let a = 12345; a = 56789; a`, 56789},
			{`let a = "foo"; a = 56789; a`, 56789},

			{`let a = {}; a[1] = "foo"; a[1]`, "foo"},
			{`let a = {}; a["foo"] = "bar"; a["foo"]`, "bar"},
			{`let a = {1: "foo", 2: "bar"}; a[2] = "baz"; a[2]`, "baz"},

			{`let a = [1,2,3]; a[1] = "foo"; a[1]`, "foo"},

			{`let a = {}; a.foo = 12; a.bar = "baz"; a.foo;`, 12},
			{`let a = {}; a.foo = 12; a.bar = "baz"; a.bar;`, "baz"},
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
		expecteds := []expectedArray{
			{`[1, 2, 3]`, 3, []string{"1", "2", "3"}},
			{`[5.1, 6.2, 7.3]`, 3, []string{"5.1", "6.2", "7.3"}},
			{`["a", "b", "c"]`, 3, []string{"a", "b", "c"}},
			{`[5.1, "a", 2, 1]`, 4, []string{"5.1", "a", "2", "1"}},
		}

		for index, expected := range expecteds {
			Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
				evaluated := testEval(expected.source)

				testArrayObject(evaluated, expected)
			})
		}
	})
}

func TestHashLiteralExpression(t *testing.T) {
	Convey("Hash literal expression test", t, func() {
		expecteds := []expectedHash{
			{`{ "foo": 1, "bar": 2 }`, 2, []string{"foo:1", "bar:2"}},
			{`{ 1: "foo", 2: "bar" }`, 2, []string{"1:foo", "2:bar"}},
			{`{ 5.5: "foo", 6.6: "bar" }`, 2, []string{"5.5:foo", "6.6:bar"}},
			{`{ true: "foo", false: "bar" }`, 2, []string{"true:foo", "false:bar"}},

			{`{ "z": 10, "d": 20, "a": 1 }`, 3, []string{"z:10", "d:20", "a:1"}},
			{`{ 20: "c", 10: "h", 30: "e", 12: "d" }`, 4, []string{"20:c", "10:h", "30:e", "12:d"}},
			{`{ "k": 1, 2.2: "g", 1: "5", "e": "9" }`, 4, []string{"k:1", "2.2:g", "1:5", "e:9"}},
		}

		for index, expected := range expecteds {
			Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
				evaluated := testEval(expected.source)

				testHashObject(evaluated, expected)
			})
		}
	})
}

func TestFunctionLiteralExpression(t *testing.T) {
	Convey("Function literal expression test", t, func() {
		expecteds := []struct {
			source          string
			parameterLength int
			blockLength     int
		}{
			{"func(a, b, c) { d }", 3, 1},
			{"func(a, b) { c; d }", 2, 2},
		}

		for index, expected := range expecteds {
			Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
				evaluated := testEval(expected.source)

				testFunctionObject(evaluated, expected)
			})
		}
	})
}

func TestRangeExpression(t *testing.T) {
	Convey("Range expression test", t, func() {
		expecteds := []struct {
			source   string
			length   int
			elements []string
		}{
			{`1..5`, 4, []string{"1", "2", "3", "4"}},
			{`3.1..3.6`, 5, []string{"3.1", "3.2", "3.3", "3.4", "3.5"}},
			{`"a".."c"`, 2, []string{"a", "b"}},
			{`"f".."a"`, 5, []string{"f", "e", "d", "c", "b"}},
			{`"z".."v"`, 4, []string{"z", "y", "x", "w"}},
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
					// Format float value to 1 decimal places
					if strings.Contains(element.Inspect(), ".") {
						floatValue, _ := strconv.ParseFloat(element.Inspect(), 64)

						compareElements = append(compareElements, fmt.Sprintf("%0.1f", floatValue))
					} else {
						compareElements = append(compareElements, element.Inspect())
					}
				}

				Convey(runMessage(`Elements should equals %s`, expected.elements), func() {
					So(compareElements, ShouldResemble, expected.elements)
				})
			})
		}
	})
}

func TestCallExpression(t *testing.T) {
	Convey("Call expression test", t, func() {
		expecteds := []struct {
			source string
			result interface{}
		}{
			{`func a() { return 123; }; a();`, 123},
			{`func a() { return 12.3; }; a();`, 12.3},
			{`func a() { return "123"; }; a();`, "123"},

			{`func a(b) { return b; }; a("foo");`, "foo"},
			{`func a(b, c, d) { return d; }; a("foo", 123, 4.5);`, 4.5},

			{`func a() { let b = "foo"; return b; }; a();`, "foo"},
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
			expecteds := []struct {
				source string
				result interface{}
			}{
				{"[1, 2, 3][2]", 3},
				{`[1.1, 2.2, 3.3][0]`, 1.1},
				{`["a", "b", "c"][1]`, "b"},
			}

			for index, expected := range expecteds {
				Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
					evaluated := testEval(expected.source)

					testLiteralObject(evaluated, expected.result)
				})
			}
		})

		Convey("For hash object", func() {
			expecteds := []struct {
				source string
				result interface{}
			}{
				{`{ 1: "a", "2": 7.2, 3.1: 50 }[1]`, "a"},
				{`{ 1: "a", "2": 7.2, 3.1: 50 }["2"]`, 7.2},
				{`{ 1: "a", "2": 7.2, 3.1: 50 }[3.1]`, 50},
			}

			for index, expected := range expecteds {
				Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
					evaluated := testEval(expected.source)

					testLiteralObject(evaluated, expected.result)
				})
			}
		})

		Convey("For string object", func() {
			expecteds := []struct {
				source string
				result interface{}
			}{
				{`"foobar"[0]`, "f"},
				{`"foobar"[3]`, "b"},
				{`"foobar"[5]`, "r"},
			}

			for index, expected := range expecteds {
				Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
					evaluated := testEval(expected.source)

					testLiteralObject(evaluated, expected.result)
				})
			}
		})

		Convey("For error when index type incorrect", func() {
			expecteds := []struct {
				source string
				result string
			}{
				{`[1, 2]["10"]`, "Index operator not support for 10 on ARRAY_OBJECT"},
				{`"foobar"["10"]`, "Index operator not support for 10 on STRING_OBJECT"},
			}

			for index, expected := range expecteds {
				Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
					evaluated := testEval(expected.source)

					testErrorObject(evaluated, expected.result)
				})
			}
		})
	})
}

func TestPrefixExpression(t *testing.T) {
	Convey("Prefix expression test", t, func() {
		Convey("Bang operator", func() {
			expecteds := []struct {
				source string
				result bool
			}{
				{`!1`, false},
				{`!0`, true},

				{`!1.1`, false},
				{`!0.0`, true},

				{`!"foo"`, false},
				{`!""`, true},

				{`![1,2]`, false},
				{`![]`, true},

				{`!{1:"a", 2:"b"}`, false},
				{`!{}`, true},

				{`!!""`, false},
				{`!!!""`, true},
				{`!!0`, false},
				{`!!!1`, false},
				{`!!0.0`, false},
				{`!!!1.1`, false},
			}

			for index, expected := range expecteds {
				Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
					evaluated := testEval(expected.source)

					testLiteralObject(evaluated, expected.result)
				})
			}
		})

		Convey("Minus operator", func() {
			expecteds := []struct {
				source string
				result interface{}
			}{
				{"-5", -5},
				{"-10", -10},

				{"-5.5", -5.5},
				{"-10.10", -10.10},
			}

			for index, expected := range expecteds {
				Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
					evaluated := testEval(expected.source)

					testLiteralObject(evaluated, expected.result)
				})
			}
		})

		Convey("Plus operator", func() {
			expecteds := []struct {
				source string
				result interface{}
			}{
				{"+5", 5},
				{"+10", 10},

				{"+5.5", 5.5},
				{"+10.10", 10.10},
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

func TestInfixExpression(t *testing.T) {
	Convey("Infix expression test", t, func() {
		Convey("&& operator test", func() {
			expecteds := []struct {
				source string
				result interface{}
			}{
				{`1 && 2`, true},
				{`3.3 && 4.4`, true},

				{`"foo" && "bar"`, true},

				{`[] && []`, false},
				{`[1, 2] && []`, false},
				{`[1, 2] && [3, 4]`, true},

				{`{} && {}`, false},
				{`{1: "a"} && {}`, false},
				{`{1: "a"} && {"b": 2}`, true},
			}

			for index, expected := range expecteds {
				Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
					evaluated := testEval(expected.source)

					testLiteralObject(evaluated, expected.result)
				})
			}
		})

		Convey("|| operator test", func() {
			expecteds := []struct {
				source string
				result interface{}
			}{
				{`1 || 2`, true},
				{`3.3 || 4.4`, true},

				{`"foo" || "bar"`, true},

				{`[] || []`, false},
				{`[1, 2] || []`, true},
				{`[1, 2] || [3, 4]`, true},

				{`{} || {}`, false},
				{`{1: "a"} || {}`, true},
				{`{1: "a"} || {"b": 2}`, true},
			}

			for index, expected := range expecteds {
				Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
					evaluated := testEval(expected.source)

					testLiteralObject(evaluated, expected.result)
				})
			}
		})

		Convey("Integer with integer operator test", func() {
			expecteds := []struct {
				source string
				result interface{}
			}{
				{`1 + 2`, 3},
				{`1 - 2`, -1},
				{`3 * 2`, 6},
				{`6 / 2`, 3},

				{`1 < 2`, true},
				{`1 > 2`, false},
				{`1 >= 2`, false},
				{`1 <= 2`, true},
				{`1 == 1`, true},
				{`1 != 2`, true},
			}

			for index, expected := range expecteds {
				Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
					evaluated := testEval(expected.source)

					testLiteralObject(evaluated, expected.result)
				})
			}
		})

		Convey("Integer with float operator test", func() {
			expecteds := []struct {
				source string
				result interface{}
			}{
				{`1 + 2.2`, 3.2},
				{`1 - 2.3`, -1.3},
				{`3 * 2.3`, 6.9},
				{`6 / 2.5`, 2.4},

				{`1 < 2.2`, true},
				{`1 > 2.3`, false},
				{`1 >= 0.4`, true},
				{`1 <= 1.5`, true},
				{`1 == 1.0`, true},
				{`1 != 2.7`, true},
			}

			for index, expected := range expecteds {
				Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
					evaluated := testEval(expected.source)

					testLiteralObject(evaluated, expected.result)
				})
			}
		})

		Convey("Float with float operator test", func() {
			expecteds := []struct {
				source string
				result interface{}
			}{
				{`1.1 + 2.2`, 3.3},
				{`1.3 - 2.3`, -1.0},
				{`3.3 * 2.3`, 7.59},
				{`6.8 / 2.5`, 2.72},

				{`1.3 < 2.2`, true},
				{`1.5 > 2.3`, false},
				{`1.7 >= 0.4`, true},
				{`2.5 <= 1.5`, false},
				{`3.3 == 3.3`, true},
				{`10.5 != 2.7`, true},
			}

			for index, expected := range expecteds {
				Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
					evaluated := testEval(expected.source)

					testLiteralObject(evaluated, expected.result)
				})
			}
		})

		Convey("Float with integer operator test", func() {
			expecteds := []struct {
				source string
				result interface{}
			}{
				{`2.2 + 1`, 3.2},
				{`2.3 - 1`, 1.3},
				{`2.3 * 3`, 6.9},
				{`8.4 / 2`, 4.2},

				{`2.2 < 1`, false},
				{`2.3 > 1`, true},
				{`0.4 >= 1`, false},
				{`1.5 <= 1`, false},
				{`1.0 == 1`, true},
				{`2.7 != 1`, true},
			}

			for index, expected := range expecteds {
				Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
					evaluated := testEval(expected.source)

					testLiteralObject(evaluated, expected.result)
				})
			}
		})

		Convey("String with string operator test", func() {
			expecteds := []struct {
				source string
				result interface{}
			}{
				{`"foo" + "bar"`, "foobar"},

				{`"a" < "b"`, true},
				{`"a" > "b"`, false},
				{`"a" <= "b"`, true},
				{`"a" >= "b"`, false},
				{`"a" == "b"`, false},
			}

			for index, expected := range expecteds {
				Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
					evaluated := testEval(expected.source)

					testLiteralObject(evaluated, expected.result)
				})
			}
		})

		Convey("Array with array operator test", func() {
			Convey("+ operator test", func() {
				expected := expectedArray{
					source:   `[1, 2.2] + ["foo", "bar"]`,
					length:   4,
					elements: []string{"1", "2.2", "foo", "bar"},
				}

				evaluated := testEval(expected.source)

				testArrayObject(evaluated, expected)
			})

			Convey("compare operator test", func() {
				expecteds := []struct {
					source string
					result interface{}
				}{
					//
					{`[1, 2] == [1, 2, 3]`, false},
					{`[1, 2] == [1, 2]`, true},
					{`[1, 2] == [3, 2]`, false},
					{`[1, 2] == [1, 3]`, false},

					{`[1, 2.2, "foo"] == [1, 2.2, "foo"]`, true},

					{`[0.1] == [0.1]`, true},
					{`[0.1] == [0.2]`, false},

					{`["foo"] == ["foo"]`, true},
					{`["foo"] == ["bar"]`, false},

					//
					{`[1, 2] != [1, 2, 3]`, true},
					{`[1, 2] != [1, 2]`, false},
					{`[1, 2] != [3, 2]`, true},
					{`[1, 2] != [1, 3]`, true},

					{`[1, 2.2, "foo"] != [1, 2.2, "foo"]`, false},

					{`[0.1] != [0.1]`, false},
					{`[0.1] != [0.2]`, true},

					{`["foo"] != ["foo"]`, false},
					{`["foo"] != ["bar"]`, true},
				}

				for index, expected := range expecteds {
					Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
						evaluated := testEval(expected.source)

						testLiteralObject(evaluated, expected.result)
					})
				}
			})
		})

		Convey("Hash with hash operator test", func() {
			Convey("+ operator test", func() {
				expected := expectedHash{
					source: `{ 1: 2, 3: 4 } + { 5.5: 6.6, "foo": "bar" }`,
					length: 4,
					order:  []string{"1:2", "3:4", "5.5:6.6", "foo:bar"},
				}

				evaluated := testEval(expected.source)

				testHashObject(evaluated, expected)
			})

			Convey("compare operator test", func() {
				expecteds := []struct {
					source string
					result interface{}
				}{
					//
					{`{ 1: 2, 3.3: 4.4, "foo": "bar" } == { 1: 2, 3.3: 4.4 }`, false},

					{`{ 1: 2, 3.3: 4.4, "foo": "bar" } == { 1: 2, 3.3: 4.4, "foo": "bar" }`, true},
					{`{ 1: 2, 3.3: 4.4, "foo": "bar" } == { 99: 2, 3.3: 4.4, "foo": "bar" }`, false},
					{`{ 1: 2, 3.3: 4.4, "foo": "bar" } == { 1: 99, 3.3: 4.4, "foo": "bar" }`, false},

					//
					{`{ 1: 2, 3.3: 4.4, "foo": "bar" } != { 1: 2, 3.3: 4.4 }`, true},

					{`{ 1: 2, 3.3: 4.4, "foo": "bar" } != { 1: 2, 3.3: 4.4, "foo": "bar" }`, false},
					{`{ 1: 2, 3.3: 4.4, "foo": "bar" } != { 99: 2, 3.3: 4.4, "foo": "bar" }`, true},
					{`{ 1: 2, 3.3: 4.4, "foo": "bar" } != { 1: 99, 3.3: 4.4, "foo": "bar" }`, true},
				}

				for index, expected := range expecteds {
					Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
						evaluated := testEval(expected.source)

						testLiteralObject(evaluated, expected.result)
					})
				}
			})

			Convey("== operator test", func() {
				expecteds := []struct {
					source string
					result interface{}
				}{
					{`true == true`, true},
					{`false == false`, true},

					{`true == false`, false},

					{`(5 > 10) == false`, true},
					{`(5 < 10) == true`, true},

					{`1 == "1"`, false},
					{`[] == {}`, false},
				}

				for index, expected := range expecteds {
					Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
						evaluated := testEval(expected.source)

						testLiteralObject(evaluated, expected.result)
					})
				}
			})

			Convey("!= operator test", func() {
				expecteds := []struct {
					source string
					result interface{}
				}{
					{`true != true`, false},
					{`false != false`, false},

					{`true != false`, true},

					{`(5 > 10) != false`, false},
					{`(5 < 10) != true`, false},

					{`1 != "1"`, true},
					{`[] != {}`, true},
				}

				for index, expected := range expecteds {
					Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
						evaluated := testEval(expected.source)

						testLiteralObject(evaluated, expected.result)
					})
				}
			})
		})
	})
}

func TestBreakExpression(t *testing.T) {
	Convey("Break expression test", t, func() {
		expected := struct {
			source string
			result string
		}{
			`break`, "break",
		}

		evaluated := testEval(expected.source)

		testBreakObject(evaluated, expected.result)
	})
}

func TestContinueExpression(t *testing.T) {
	Convey("Continue expression test", t, func() {
		expected := struct {
			source string
			result string
		}{
			`continue`, "continue",
		}

		evaluated := testEval(expected.source)

		testContinueObject(evaluated, expected.result)
	})
}

func TestIfExpression(t *testing.T) {
	Convey("If expression test", t, func() {
		expecteds := []struct {
			source string
			result interface{}
		}{
			{`if (5 < 10) { return true }`, true},
			{`if (5 > 10) { return true } else { return false }`, false},
			{`if (5 > 10) { return "Y" } else if (5 == 5) { return "=="; }else{ return "N" }`, "=="},
		}

		for index, expected := range expecteds {
			Convey(runMessage("Running: %d, Source: %s", index, expected.source), func() {
				evaluated := testEval(expected.source)

				testLiteralObject(evaluated, expected.result)
			})
		}
	})
}

func TestForEverExpression(t *testing.T) {
	Convey("For ever expresion test", t, func() {
		expecteds := []struct {
			source string
			result interface{}
		}{
			{
				`let a = 0;
				for {
					let a = a + 1;
					if (a < 10) { continue; }else{ break; }
				}
				a;`,
				10,
			},
			{
				`let a = 1;
				for {
					let a = a + 1;

					if (a / 2 == 1) {
						return a;
					}
				}
				a;`,
				2,
			},
		}

		for index, expected := range expecteds {
			Convey(runMessage("Running: %d, Source: %s", index, codeToSingleLine(expected.source)), func() {
				evaluated := testEval(expected.source)

				testLiteralObject(evaluated, expected.result)
			})
		}
	})
}

func TestForEachArrayOrRangeExpression(t *testing.T) {
	Convey("For each array or range test", t, func() {
		expecteds := []struct {
			source string
			result interface{}
		}{
			{
				`let a = 0;
				for b in 1..4 {
					let a = a + b;
				}
				a;`,
				6,
			},
			{
				`let a = 0;
				for b in 0.1..0.4 {
					let a = a + b;
				}
				a;`,
				0.6,
			},
			{
				`let a = "";
				for b in "a".."d" {
					let a = a + b;
				}
				a;`,
				"abc",
			},
			{
				`let a = 0;
				for b in [1,2,3,4] {
					let a = a + b;
				}
				a;`,
				10,
			},
			{
				`let a = 0;
				for b in [1,2,3,4] {
					if (b == 3) {
						break;
					}else{
						let a = a + b;
						continue;
					}
				}
				a;`,
				3,
			},
			{
				`let a = 0;
				for b in [1,2,3,4] {
					let a = a + b;

					if (a > 4) {
						return a;
					}
				}
				a;`,
				6,
			},
		}

		for index, expected := range expecteds {
			Convey(runMessage("Running: %d, Source: %s", index, codeToSingleLine(expected.source)), func() {
				evaluated := testEval(expected.source)

				testLiteralObject(evaluated, expected.result)
			})
		}
	})
}

func TestForEachHashExpression(t *testing.T) {
	Convey("For each hash test", t, func() {
		expecteds := []struct {
			source string
			result interface{}
		}{
			{
				`
				for k, v in { "a": 1, "b": 3 } {
					return k;
				}
				`,
				"a",
			},
			{
				`
				for k, v in { "a": 1, "b": 3 } {
					return v;
				}
				`,
				1,
			},
			{
				`
				for k, v in { 1: "a", 3: 4 } {
					if (k == 3) {
						return v;
					}
				}
				`,
				4,
			},
			{
				`
				let c = 0;
				for k,v in { "a": 1, "b": 2 } {
					let c = c + v;

					if (c == 3) {
						return c;
					}
				}
				c`,
				3,
			},
		}

		for index, expected := range expecteds {
			Convey(runMessage("Running: %d, Source: %s", index, codeToSingleLine(expected.source)), func() {
				evaluated := testEval(expected.source)

				testLiteralObject(evaluated, expected.result)
			})
		}
	})
}

// Statements
func TestLetStatement(t *testing.T) {
	Convey("Let statement test", t, func() {
		expecteds := []struct {
			source string
			value  interface{}
		}{
			{`let a = 5;`, 5},
			{`let b = 5.5;`, 5.5},
			{`let c = "foo";`, "foo"},
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

		testFunctionObject(evaluated, expectedFunction{
			parameterLength: 2,
			blockLength:     1,
		})
	})
}

func TestReturnStatement(t *testing.T) {
	Convey("Return statement test", t, func() {
		expecteds := []struct {
			source string
			result interface{}
		}{
			{"return 10", 10},
			{"return 15.5", 15.5},
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
		expecteds := []expectedFunction{
			{"func myFunc1(a, b, c) { d }", 3, 1},
			{"func myFunc2(a, b) { c; d }", 2, 2},
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
		expecteds := []struct {
			source      string
			returnValue interface{}
		}{
			{"let a = 1;", 1},
			{"let a = 1; let b = 2;", 2},

			{"let a = 1.1;", 1.1},
			{"let a = 1.1; let b = 2.2;", 2.2},

			{`let a = "foo";`, "foo"},
			{`let a = "foo"; let b = "bar"`, "bar"},

			{`let a = "";`, ""},
			{`let a = "foobar"; return a;`, "foobar"},
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
	theLexer := lexer.NewLexer(source)
	theParser := parser.NewParser(theLexer)
	theProgarm := theParser.Parse()

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

func testFunctionObject(obj object.Object, expected expectedFunction) {
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

func testNilObject(obj object.Object, expected string) {
	result, ok := obj.(*object.Nil)
	Convey("Can convert object (nil)", func() {
		So(ok, ShouldBeTrue)
	})

	Convey("Nil was matched", func() {
		So(result.Inspect(), ShouldEqual, expected)
	})
}

func testErrorObject(obj object.Object, expected string) {
	result, ok := obj.(*object.Error)

	Convey("Can convert to object (error)", func() {
		So(ok, ShouldBeTrue)
	})

	Convey("Error message was matched", func() {
		So(result.Message, ShouldEqual, expected)
	})
}

func testArrayObject(obj object.Object, expected expectedArray) {
	array, ok := obj.(*object.Array)
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
}

func testHashObject(obj object.Object, expected expectedHash) {
	hash, ok := obj.(*object.Hash)
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
}

func testBreakObject(obj object.Object, expected string) {
	result, ok := obj.(*object.Break)
	Convey("Can convert object (break)", func() {
		So(ok, ShouldBeTrue)
	})

	Convey("Break was matched", func() {
		So(result.Inspect(), ShouldEqual, expected)
	})
}

func testContinueObject(obj object.Object, expected string) {
	result, ok := obj.(*object.Continue)
	Convey("Can convert object (continue)", func() {
		So(ok, ShouldBeTrue)
	})

	Convey("Continue was matched", func() {
		So(result.Inspect(), ShouldEqual, expected)
	})
}

// Helper functions for common
func runMessage(format string, values ...interface{}) string {
	return fmt.Sprintf(format, values...)
}

func codeToSingleLine(text string) string {
	var result string

	reTab := regexp.MustCompile(`\t+`)
	reLineBreak := regexp.MustCompile((`\n|\r|\r\n`))

	result = reTab.ReplaceAllString(text, "")
	result = reLineBreak.ReplaceAllString(result, " ")

	return result
}
