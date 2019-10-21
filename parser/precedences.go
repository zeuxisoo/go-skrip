package parser

import (
	"github.com/zeuxisoo/go-skrip/token"
)

const (
	_           int = iota
	LOWEST          //
	ANDOR           // || or &&
	ASSIGN          // =
	EQUALS          // ==
	LESSGREATER     // > or <
	SUM             // +
	PRODUCT         // *
	RANGE           // ..
	PREFIX          // -X or !X
	CALL            // func(X)
	INDEX           // array[index]
	DOT             // any.function() or any.property
)

var precedences = map[token.Type]int{
	token.AND:              ANDOR,
	token.OR:               ANDOR,
	token.ASSIGN:           ASSIGN,
	token.EQ:               EQUALS,
	token.NOT_EQ:           EQUALS,
	token.LT:               LESSGREATER,
	token.LTEQ:             LESSGREATER,
	token.GT:               LESSGREATER,
	token.GTEQ:             LESSGREATER,
	token.PLUS:             SUM,
	token.MINUS:            SUM,
	token.SLASH:            PRODUCT,
	token.ASTERISK:         PRODUCT,
	token.RANGE:            RANGE,
	token.LEFT_PARENTHESIS: CALL,
	token.LEFT_BRACKET:     INDEX,
	token.DOT:              DOT,
}
