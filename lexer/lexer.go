package lexer

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

	return lexer
}
