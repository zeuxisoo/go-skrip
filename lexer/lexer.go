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

	lexer.readChar()

	return lexer
}

func (l *Lexer) readChar() {
	// Reset to 0 when next position greater than source length
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
