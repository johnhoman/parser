package object

type Function struct{}

func (f *Function) Type() Type {
	panic("implement me")
}

func (f *Function) Inspect() string {
	panic("implement me")
}

var _ Object = &Function{}
