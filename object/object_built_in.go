package object

type BuiltInFunction func(env *Environment, arguments ...Object) Object

type BuiltIn struct {
	Function BuiltInFunction
}

func (b *BuiltIn) Type() ObjectType {
	return BUILTIN_OBJECT
}

func (b *BuiltIn) Inspect() string {
	return "built-in function"
}
