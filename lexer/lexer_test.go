package lexer

import (
	"fmt"
	"testing"

	"github.com/zeuxisoo/go-skriplang/token"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLexerAssign(t *testing.T) {
	Convey("Basic assign testing", t, func() {
		source := `
			let five = 5;
			let ten = 10;

			let add = func(x, y) {
				x + y;
			};

			let result = add(five, ten);
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

			{ token.LET, "let" },
			{ token.IDENTIFIER, "result" },
			{ token.ASSIGN, "=" },
			{ token.IDENTIFIER, "add" },
			{ token.LEFT_PARENTHESIS, "(" },
			{ token.IDENTIFIER, "five" },
			{ token.COMMA, "," },
			{ token.IDENTIFIER, "ten" },
			{ token.RIGHT_PARENTHESIS, ")" },
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

func TestLexerOperator(t *testing.T) {
	Convey("Operator testing", t, func() {
		source := `
			!-/*5;
			5 < 10 > 5;
			5 <= 10 >= 5;
		`

		testTokens := []struct{
			expectedType 	token.Type
			expectedLiteral string
		}{
			{ token.BANG, "!" },
			{ token.MINUS, "-" },
			{ token.SLASH, "/" },
			{ token.ASTERISK, "*" },
			{ token.INT, "5" },
			{ token.SEMICOLON, ";" },

			{ token.INT, "5" },
			{ token.LT, "<" },
			{ token.INT, "10" },
			{ token.GT, ">" },
			{ token.INT, "5" },
			{ token.SEMICOLON, ";" },

			{ token.INT, "5" },
			{ token.LTEQ, "<=" },
			{ token.INT, "10" },
			{ token.GTEQ, ">=" },
			{ token.INT, "5" },
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
