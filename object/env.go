package object

import "sync"

type Env struct{ objects sync.Map }

func (env *Env) Get(name string) (Object, bool) {
	v, ok := env.objects.Load(name)
	if !ok {
		return nil, false
	}
	obj := v.(Object)
	return obj, true
}

func (env *Env) Set(name string, obj Object) Object {
	env.objects.Store(name, obj)
	return obj
}

func NewEnv() *Env {
	return &Env{sync.Map{}}
}
