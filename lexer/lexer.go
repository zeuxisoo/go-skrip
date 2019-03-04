package lexer

import (
	"strings"

	"github.com/zeuxisoo/go-skriplang/token"
	"github.com/zeuxisoo/go-skriplang/pkg/helper"
)

type Lexer struct {
	source			string
	currentChar		byte	// current character
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

	switch l.currentChar {
	case 0:
		theToken.Literal = ""
		theToken.Type    = token.EOF
	default:
		if helper.IsLetter(l.currentChar) {
			theToken.Literal = l.readIdentifier()
			theToken.Type    = token.FindKeywordType(theToken.Literal)

			return theToken
		}

		if helper.IsDigit(l.currentChar) {
			theToken.Literal = l.readNumber()

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
		l.currentChar = l.source[l.nextPosition]
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
		l.readChar()
	}

	return l.source[startPosition:l.currentPosition]
}

func (l *Lexer) newIllegalToken(literal string) token.Token {
	return token.Token{
		Type      : token.ILLEGAL,
		Literal   : literal,
		LineNumber: l.currentLine,
	}
}
