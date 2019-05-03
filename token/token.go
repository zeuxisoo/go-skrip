package token

// Type will represent each token type
type Type string

// Token will store the input value
type Token struct {
	Type 		Type
	Literal 	string
	LineNumber 	int
}

var keywords = map[string]Type{
	"func"    : FUNCTION,
	"let"     : LET,
	"true"    : TRUE,
	"false"   : FALSE,
	"if"      : IF,
	"else"    : ELSE,
	"return"  : RETURN,
	"for"     : FOR,
	"in"      : IN,
	"nil"     : NIL,
	"break"   : BREAK,
	"continue": CONTINUE,
}

// FindKeywordType will return keyword type
func FindKeywordType(literal string) Type {
	if keyword, ok := keywords[literal]; ok {
		return keyword
	}

	return IDENTIFIER
}
