package object

type Environment struct {
	store map[string]Object
}

func NewEnviroment() *Environment {
	s := make(map[string]Object)
	return &Environment{store:s}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

func (e *Environment) Set(name string, value Object) Object {
	e.store[name] = value
	return value
}
