package token

const (
	//
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers and literals
	IDENTIFIER = "IDENTIFIER"	// function name, variable name, etc
	INT        = "INT"          // 12345
	FLOAT      = "FLOAT"        // 12.345
	STRING     = "STRING"       // "text"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT   = "<"
	LTEQ = "<="
	GT   = ">"
	GTEQ = ">="

	EQ      = "=="
	NOT_EQ  = "!="
	AND     = "&&"
	OR      = "||"

	DOT     = "."
	RANGE   = ".."

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"


	LEFT_PARENTHESIS  = "("
	RIGHT_PARENTHESIS = ")"
	LEFT_BRACE        = "{"
	RIGHT_BRACE       = "}"
	LEFT_BRACKET      = "["
	RIGHT_BRACKET     = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	FOR      = "FOR"
	IN       = "IN"
	NIL      = "NIL"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
)
