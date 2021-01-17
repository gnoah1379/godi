package godi

type PreLoader interface {
	Preload(ctx Context)
}

type Runner interface {
	Run(ctx Context) error
}

type ShutDowner interface {
	ShutDown(ctx Context)
}

//func isService(t reflect.Type) bool {
//	return t.Implements(reflect.TypeOf(new(Service)).Elem())
//}
