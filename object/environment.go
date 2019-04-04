package object

type Environment struct {
	store  map[string]Object
	parent *Environment
}

func NewEnvironment() *Environment {
	return &Environment{
		store : make(map[string]Object),
		parent: nil,
	}
}

func (env *Environment) Get(name string) (Object, bool) {
	obj, ok := env.store[name]

	// Try get from parent store when current store is not available
	if ok == false && env.parent != nil {
		obj, ok = env.parent.Get(name)
	}

	return obj, ok
}

func (env *Environment) Set(name string, value Object) Object {
	env.store[name] = value

	return value
}
