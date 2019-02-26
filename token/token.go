package token

// Kind will represent each token type
type Kind string

// Token will store the input value
type Token struct {
	Type 		Kind
	Literal 	string
	LineNumber 	int
}

var keywords = map[string]Kind{
	"func":   FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"for":	  FOR,
}

// LookupKeywordType will return keyword type or plain ident
func LookupKeywordType(ident string) Kind {
	if keyword, ok := keywords[ident]; ok {
		return keyword
	}

	return IDENT
}
