package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnviroment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func NewClosedEnvironments(outer *Environment) *Environment {
	env := NewEnviroment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, value Object) Object {
	e.store[name] = value
	return value
}
