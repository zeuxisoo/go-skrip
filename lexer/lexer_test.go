package lexer

import (
	"fmt"
	"testing"

	"github.com/zeuxisoo/go-skriplang/token"

	. "github.com/smartystreets/goconvey/convey"
)

//
type expectedToken struct{
	Type 	token.Type
	Literal string
}

//
func compareToken(theLexer *Lexer, expectedTokens []expectedToken) {
	for index, currentExpectedToken := range expectedTokens {
		lexerToken := theLexer.NextToken()

		messageType  := fmt.Sprintf(
			"Running %d, Got: %s, %s, Expected: %s, %s",
			index,
			lexerToken.Type, lexerToken.Literal,
			currentExpectedToken.Type, currentExpectedToken.Literal,
		)

		Convey(messageType, func() {
			So(lexerToken.Type, ShouldEqual, currentExpectedToken.Type)
			So(lexerToken.Literal, ShouldEqual, currentExpectedToken.Literal)
		})
	}
}

//
func TestLexerAssign(t *testing.T) {
	Convey("Basic assign testing", t, func() {
		source := `
			let five = 5;
			let ten = 10;
			let five_float = 5.00;
			let hello_world = "Hello world";

			let add = func(x, y) {
				x + y;
			};

			let result = add(five, ten);

			let key_value = { "foo": "bar" };

			let array = [1, 2];
		`;

		expectedTokens := []expectedToken{
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
			{ token.IDENTIFIER, "five_float" },
			{ token.ASSIGN, "=" },
			{ token.FLOAT, "5.00" },
			{ token.SEMICOLON, ";" },

			{ token.LET, "let" },
			{ token.IDENTIFIER, "hello_world" },
			{ token.ASSIGN, "=" },
			{ token.STRING, "Hello world" },
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

			{ token.LET, "let" },
			{ token.IDENTIFIER, "key_value" },
			{ token.ASSIGN, "=" },
			{ token.LEFT_BRACE, "{" },
			{ token.STRING, "foo" },
			{ token.COLON, ":" },
			{ token.STRING, "bar" },
			{ token.RIGHT_BRACE, "}" },
			{ token.SEMICOLON, ";" },

			{ token.LET, "let" },
			{ token.IDENTIFIER, "array" },
			{ token.ASSIGN, "=" },
			{ token.LEFT_BRACKET, "[" },
			{ token.INT, "1" },
			{ token.COMMA, "," },
			{ token.INT, "2" },
			{ token.RIGHT_BRACKET, "]" },
			{ token.SEMICOLON, ";" },

			{ token.EOF, "" },
		}

		compareToken(NewLexer(source), expectedTokens)
	})
}

func TestLexerOperator(t *testing.T) {
	Convey("Operator testing", t, func() {
		source := `
			!+-=/*5;

			5 != 10;

			5 < 10 > 5;

			5 <= 10 >= 5;

			5 && 5;

			5 || 5;
		`

		expectedTokens := []expectedToken{
			{ token.BANG, "!" },
			{ token.PLUS, "+" },
			{ token.MINUS, "-" },
			{ token.ASSIGN, "=" },
			{ token.SLASH, "/" },
			{ token.ASTERISK, "*" },
			{ token.INT, "5" },
			{ token.SEMICOLON, ";" },

			{ token.INT, "5" },
			{ token.NOT_EQ, "!=" },
			{ token.INT, "10" },
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

			{ token.INT, "5" },
			{ token.AND, "&&" },
			{ token.INT, "5" },
			{ token.SEMICOLON, ";" },

			{ token.INT, "5" },
			{ token.OR, "||" },
			{ token.INT, "5" },
			{ token.SEMICOLON, ";" },

			{ token.EOF, "" },
		}

		compareToken(NewLexer(source), expectedTokens)
	})
}

func TestLexerKeywords(t *testing.T) {
	Convey("Keywords testing", t, func() {
		source := `
			let name = func(x, y) {
			}

			if (x > 5) {
				return true;
			}else{
				return false;
			}

			for (index, value) in array_data {
			}
		`

		expectedTokens := []expectedToken{
			{ token.LET, "let" },
			{ token.IDENTIFIER, "name" },
			{ token.ASSIGN, "=" },
			{ token.FUNCTION, "func" },
			{ token.LEFT_PARENTHESIS, "(" },
			{ token.IDENTIFIER, "x" },
			{ token.COMMA, "," },
			{ token.IDENTIFIER, "y" },
			{ token.RIGHT_PARENTHESIS, ")" },
			{ token.LEFT_BRACE, "{" },
			{ token.RIGHT_BRACE, "}" },

			{ token.IF, "if" },
			{ token.LEFT_PARENTHESIS, "(" },
			{ token.IDENTIFIER, "x" },
			{ token.GT, ">" },
			{ token.INT, "5" },
			{ token.RIGHT_PARENTHESIS, ")" },
			{ token.LEFT_BRACE, "{" },
			{ token.RETURN, "return" },
			{ token.TRUE, "true" },
			{ token.SEMICOLON, ";" },
			{ token.RIGHT_BRACE, "}" },
			{ token.ELSE, "else" },
			{ token.LEFT_BRACE, "{" },
			{ token.RETURN, "return" },
			{ token.FALSE, "false" },
			{ token.SEMICOLON, ";" },
			{ token.RIGHT_BRACE, "}" },

			{ token.FOR, "for" },
			{ token.LEFT_PARENTHESIS, "(" },
			{ token.IDENTIFIER, "index" },
			{ token.COMMA, "," },
			{ token.IDENTIFIER, "value" },
			{ token.RIGHT_PARENTHESIS, ")" },
			{ token.IN, "in" },
			{ token.IDENTIFIER, "array_data" },
			{ token.LEFT_BRACE, "{" },
			{ token.RIGHT_BRACE, "}" },

			{ token.EOF, "" },
		}

		compareToken(NewLexer(source), expectedTokens)
	})
}

func TestStringEscapeQuote(t *testing.T) {
	Convey("String escape quote", t, func() {
		source := `
			"this is a \"quote\" string"
		`

		theLexer := NewLexer(source)
		theToken := theLexer.NextToken()

		So(theToken.Type, ShouldEqual, token.STRING)
		So(theToken.Literal, ShouldEqual, `this is a "quote" string`)
	})
}

func TestSkipComment(t *testing.T) {
	Convey("Skip comment", t, func() {
		source := `
			// this is single line comment 1
			// this is single line comment 2

			let a = 5;
		`

		expectedTokens := []expectedToken{
			{ token.LET, "let" },
			{ token.IDENTIFIER, "a" },
			{ token.ASSIGN, "=" },
			{ token.INT, "5" },
			{ token.SEMICOLON, ";" },
		}

		compareToken(NewLexer(source), expectedTokens)
	})
}
