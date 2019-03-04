package helper

func IsLetter(chr byte) bool {
	return chr >= 'a' && chr <= 'z' || chr >= 'A' && chr <= 'Z' || chr == '_'
}


func IsDigit(chr byte) bool {
	return chr >= '0' && chr <= '9'
}

func IsDot(chr byte) bool {
	return chr == '.'
}
