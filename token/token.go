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
	"func":   FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"for":	  FOR,
}

// FindKeywordType will return keyword type or plain ident
func FindKeywordType(ident string) Type {
	if keyword, ok := keywords[ident]; ok {
		return keyword
	}

	return IDENTIFIER
}
