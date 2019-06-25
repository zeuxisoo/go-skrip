package object

type Continue struct {
}

func (c *Continue) Type() ObjectType {
	return CONTINUE_OBJECT
}

func (c *Continue) Inspect() string {
	return "continue"
}
