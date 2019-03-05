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
			let five = 5;
			let ten = 10;

			let add = func(x, y) {
				x + y;
			};
		`;

		testTokens := []struct{
			expectedType 	token.Type
			expectedLiteral string
		}{
			{ token.LET, "let" },
			{ token.IDENTIFIER, "five" },
			{ token.ASSIGN, "=" },
			{ token.INT, "5" },
			{ token.SEMICOLON, ";" },

			{ token.LET, "let" },
			{ token.IDENTIFIER, "ten" },
			{ token.ASSIGN, "=" },
			{ token.INT, "10" },
			{ token.SEMICOLON, ";" },

			{ token.LET, "let" },
			{ token.IDENTIFIER, "add" },
			{ token.ASSIGN, "=" },
			{ token.FUNCTION, "func" },
			{ token.LEFT_PARENTHESIS, "(" },
			{ token.IDENTIFIER, "x" },
			{ token.COMMA, "," },
			{ token.IDENTIFIER, "y" },
			{ token.RIGHT_PARENTHESIS, ")" },
			{ token.LEFT_BRACE, "{" },
			{ token.IDENTIFIER, "x" },
			{ token.PLUS, "+" },
			{ token.IDENTIFIER, "y" },
			{ token.SEMICOLON, ";" },
			{ token.RIGHT_BRACE, "}" },
			{ token.SEMICOLON, ";" },
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
