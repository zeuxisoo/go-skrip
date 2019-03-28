package lexer

import (
	"strings"

	"github.com/zeuxisoo/go-skrip/token"
	"github.com/zeuxisoo/go-skrip/pkg/helper"
)

type Lexer struct {
	source			string
	currentChar		rune	// current character
	currentPosition	int		// position of current character
	nextPosition	int 	// position after current character (greater than 1)
	currentLine		int		// position of current line
}

func NewLexer(source string) *Lexer {
	lexer := &Lexer{
		source: source,
		currentLine: 1,
	}

	lexer.readChar()

	return lexer
}

//
func (l *Lexer) NextToken() token.Token {
	var theToken token.Token

	l.skipWhitespace()

	if l.currentChar == '/' && l.nextChar() == '/' {
		l.skipSingleLineComment()

		return l.NextToken()
	}

	if l.currentChar == '/' && l.nextChar() == '*' {
		l.skipMultiLineComment()

		return l.NextToken()
	}

	switch l.currentChar {
	case '=':
		// if next char is '=', it should be "==" operator
		// otherwise, it should be "=" assign operator
		if l.nextChar() == '=' {
			oldCurrentChar := l.currentChar

			l.readChar()

			theToken = token.Token{
				Type   : token.EQ,
				Literal: string(oldCurrentChar) + string(l.currentChar), // text: ==
			}
		}else{
			theToken = l.newToken(token.ASSIGN)
		}
	case '+':
		theToken = l.newToken(token.PLUS)
	case ',':
		theToken = l.newToken(token.COMMA)
	case ';':
		theToken = l.newToken(token.SEMICOLON)
	case '(':
		theToken = l.newToken(token.LEFT_PARENTHESIS)
	case ')':
		theToken = l.newToken(token.RIGHT_PARENTHESIS)
	case '{':
		theToken = l.newToken(token.LEFT_BRACE)
	case '}':
		theToken = l.newToken(token.RIGHT_BRACE)
	case '[':
		theToken = l.newToken(token.LEFT_BRACKET)
	case ']':
		theToken = l.newToken(token.RIGHT_BRACKET)
	case '!':
		if l.nextChar() == '=' {
			oldCurrentChar := l.currentChar

			l.readChar()

			theToken = token.Token{
				Type   : token.NOT_EQ,
				Literal: string(oldCurrentChar) + string(l.currentChar), // text: !=
			}
		}else{
			theToken = l.newToken(token.BANG)
		}
	case '-':
		theToken = l.newToken(token.MINUS)
	case '/':
		theToken = l.newToken(token.SLASH)
	case '*':
		theToken = l.newToken(token.ASTERISK)
	case '<':
		if l.nextChar() == '=' {
			oldCurrentChar := l.currentChar

			l.readChar()

			theToken = token.Token{
				Type: token.LTEQ,
				Literal: string(oldCurrentChar) + string(l.currentChar), // text: <=
			}
		}else{
			theToken = l.newToken(token.LT)
		}
	case '>':
		if l.nextChar() == '=' {
			oldCurrentChar := l.currentChar

			l.readChar()

			theToken = token.Token{
				Type: token.GTEQ,
				Literal: string(oldCurrentChar) + string(l.currentChar), // text: >=
			}
		}else{
			theToken = l.newToken(token.GT)
		}
	case '"':
		theToken = token.Token{
			Type   : token.STRING,
			Literal: l.readString(),
		}
	case '&':
		if l.nextChar() == '&' {
			oldCurrentChar := l.currentChar

			l.readChar()

			theToken = token.Token{
				Type   : token.AND,
				Literal: string(oldCurrentChar) + string(l.currentChar), // text: &&
			}
		}else{
			theToken = l.newToken(token.ILLEGAL)
		}
	case '|':
		if l.nextChar() == '|' {
			oldCurrentChar := l.currentChar

			l.readChar()

			theToken = token.Token{
				Type   : token.OR,
				Literal: string(oldCurrentChar) + string(l.currentChar), // text: ||
			}
		}else{
			theToken = l.newToken(token.ILLEGAL)
		}
	case ':':
		theToken = l.newToken(token.COLON)
	case '.':
		if l.nextChar() == '.' {
			oldCurrentChar := l.currentChar

			l.readChar()

			theToken = token.Token{
				Type   : token.RANGE,
				Literal: string(oldCurrentChar) + string(l.currentChar),     // text: ..
			}
		}else{
			theToken = token.Token{
				Type   : token.DOT,
				Literal: l.readString(),
			}
		}
	case 0:
		theToken.Literal = ""
		theToken.Type    = token.EOF
	default:
		if helper.IsLetter(l.currentChar) {
			lineNumber := l.currentLine
			identifier := l.readIdentifier()

			return token.Token{
				Type      : token.FindKeywordType(identifier),
				Literal   : identifier,
				LineNumber: lineNumber,
			}
		}

		if helper.IsDigit(l.currentChar) {
			theToken.Literal    = l.readNumber()
			theToken.LineNumber = l.currentLine

			switch len(strings.Split(theToken.Literal, ".")) {
			case 1:	// e.g. 12, 13
				theToken.Type = token.INT
			case 2: // e.g. 12.00, 13.77
				theToken.Type = token.FLOAT
			default:
				return l.newIllegalToken(theToken.Literal)
			}

			return theToken
		}

		theToken = l.newIllegalToken(string(l.currentChar))
	}

	l.readChar()

	theToken.LineNumber = l.currentLine

	return theToken
}

//
func (l *Lexer) readChar() {
	// Reset to 0 when next position greater than source length (for EOF char)
	// Otherwise set next position to current position
	if l.nextPosition >= len(l.source) {
		l.currentChar = 0
	}else{
		l.currentChar = rune(l.source[l.nextPosition])
	}

	// Increase the current line no when encountering a newline
	if l.currentChar == '\n' {
		l.currentLine++
	}

	l.currentPosition = l.nextPosition

	l.nextPosition++
}

func (l *Lexer) skipWhitespace() {
	for l.currentChar == ' ' || l.currentChar == '\t' || l.currentChar == '\n' || l.currentChar == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipSingleLineComment() {
	for l.currentChar != '\n' && l.currentChar != 0 {
		l.readChar()
	}

	l.skipWhitespace()
}

func (l *Lexer) skipMultiLineComment() {
	stop := false

	for !stop {
		// stop when got end of file
		if l.currentChar == 0 {
			stop = true
		}

		// stop when found */
		if l.currentChar == '*' && l.nextChar() == '/' {
			stop = true

			// Read once for set the current char from "*" to "/"
			l.readChar()
		}

		l.readChar()
	}

	l.skipWhitespace()
}

func (l *Lexer) readIdentifier() string {
	startPosition := l.currentPosition

	for helper.IsLetter(l.currentChar) || helper.IsDigit(l.currentChar) {
		l.readChar()
	}

	return l.source[startPosition:l.currentPosition]
}

func (l *Lexer) readNumber() string {
	startPosition := l.currentPosition

	for helper.IsDigit(l.currentChar) || helper.IsDot(l.currentChar) {
		// When meet "." it may range "..", if range "..", return current read value
		if l.currentChar == '.' && l.nextChar() == '.' {
			return l.source[startPosition:l.currentPosition]
		}

		l.readChar()
	}

	return l.source[startPosition:l.currentPosition]
}

func (l *Lexer) readString() string {
	startPosition := l.currentPosition + 1 // skip start "

	for l.currentChar != 0 {
		l.readChar()

		// Continue when meet quote escape \"
		if l.currentChar == '\\' && l.nextChar() == '"' {
			l.readChar()
			l.readChar()
		}

		// Break when meet end of "
		if l.currentChar == '"' {
			break
		}
	}

	text := l.source[startPosition:l.currentPosition]

	// No limit for replacements \" to "
	return strings.Replace(text, "\\\"", "\"", -1)
}

func (l *Lexer) newToken(tokenType token.Type) token.Token {
	return token.Token{
		Type      : tokenType,
		Literal   : string(l.currentChar),
		LineNumber: l.currentLine,
	}
}

func (l *Lexer) newIllegalToken(literal string) token.Token {
	return token.Token{
		Type      : token.ILLEGAL,
		Literal   : literal,
		LineNumber: l.currentLine,
	}
}

func (l *Lexer) nextChar() rune {
	// e.g. End of file will return 0
	if l.nextPosition > len(l.source) {
		return 0
	}

	return rune(l.source[l.nextPosition])
}
