package helper

func IsLetter(chr rune) bool {
	return chr >= 'a' && chr <= 'z' || chr >= 'A' && chr <= 'Z' || chr == '_'
}


func IsDigit(chr rune) bool {
	return chr >= '0' && chr <= '9'
}

func IsDot(chr rune) bool {
	return chr == '.'
}
