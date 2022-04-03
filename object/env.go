package object

import "sync"

type Env struct {
	objects sync.Map
	outer   *Env
}

func (env *Env) Push() *Env {
	e := NewEnv()
	e.outer = env
	return e
}

func (env *Env) Pop() *Env {
	return env.outer
}

func (env *Env) Get(name string) (Object, bool) {
	v, ok := env.objects.Load(name)
	if !ok {
		if env.outer == nil {
			return nil, false
		}
		v, ok = env.outer.Get(name)
		if !ok {
			return nil, false
		}
	}
	obj := v.(Object)
	return obj, true
}

func (env *Env) Set(name string, obj Object) Object {
	env.objects.Store(name, obj)
	return obj
}

func (env *Env) Delete(name string) { env.objects.Delete(name) }

func NewEnv() *Env {
	return &Env{objects: sync.Map{}, outer: nil}
}
