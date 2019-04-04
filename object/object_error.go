package object

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJECT
}

func (e *Error) Inspect() string {
	return "[Error] " + e.Message
}
