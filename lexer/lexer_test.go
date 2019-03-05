package lexer

import (
	"fmt"
	"testing"

	"github.com/zeuxisoo/go-skriplang/token"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBasicLexer(t *testing.T) {
	Convey("Basic Lexer testing", t, func() {
		source := `
			let five = 5
		`;

		testTokens := []struct{
			expectedType 	token.Type
			expectedLiteral string
		}{
			{ token.LET, "let" },
			{ token.IDENTIFIER, "five" },
			{ token.ASSIGN, "=" },
			{ token.INT, "5" },
		}

		theLexer := NewLexer(source)

		for index, testToken := range testTokens {
			lexerToken := theLexer.NextToken()

			got      := lexerToken.Type
			expected := testToken.expectedType
			message  := fmt.Sprintf("Running %d, got: %s, expected: %s", index, got, expected)

			Convey(message, func() {
				So(lexerToken.Type, ShouldEqual, testToken.expectedType)
			})
		}
	})
}
