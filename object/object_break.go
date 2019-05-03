package object

type Break struct {
}

func (b *Break) Type() ObjectType {
    return BREAK_OBJECT
}

func (b *Break) Inspect() string {
    return "break"
}
