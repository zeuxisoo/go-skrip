package evaluator

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/zeuxisoo/go-skrip/lexer"
	"github.com/zeuxisoo/go-skrip/parser"
	"github.com/zeuxisoo/go-skrip/object"
)

func TestReturnStatement(t *testing.T) {
	Convey("Return statement test", t, func() {
		expecteds := []struct{
			source string
			result int64
		}{
			{"return 10", 10},
		}

		for index, expected := range expecteds {
			Convey(runMessage("Running: %d, source: %s", index, expected.source), func() {
				evaluated := testEval(expected.source)

				testDecimalObject(evaluated, expected.result)
			})
		}
	})
}

//
func testEval(source string) object.Object {
	theLexer       := lexer.NewLexer(source)
	theParser      := parser.NewParser(theLexer)
	theProgarm 	   := theParser.Parse()
	theEnvironment := object.NewEnvironment()

	return Eval(theProgarm, theEnvironment)
}

func testDecimalObject(obj object.Object, expected interface{}) {
	switch expected := expected.(type) {
	case int64:
		testIntegerObject(obj, expected)
	case float64:
		testFloatObject(obj, expected)
	}
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

}

// Helper functions for common
func runMessage(format string, values ...interface{}) string {
	return fmt.Sprintf(format, values...)
}
