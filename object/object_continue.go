package object

type Continue struct {
}

func (b *Continue) Type() ObjectType {
    return CONTINUE_OBJECT
}

func (b *Continue) Inspect() string {
    return "continue"
}
