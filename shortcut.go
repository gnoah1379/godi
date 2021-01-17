package godi

func GetContext() ContainerContext {
	return application.Context()
}

func Register(constructor interface{}, opts ...RegisterOption) error {
	return application.Register(constructor, opts)
}

func Unregister(packetName string, tags ...string) {
	application.Unregister(packetName, tags)
}

func InstanceOf(value interface{}, filters ...string) error {
	return application.InstanceOf(value, filters)
}

func SliceOf(value interface{}, filters ...string) error {
	return application.SliceOf(value, filters)
}

func MapOf(value interface{}, filters ...string) error {
	return application.MapOf(value, filters)
}

func Run() {
	application.Run()
}

func Log() Logger {
	return application.Log()
}

func SetLogger(logger Logger) {
	application.SetLogger(logger)
}

func Trigger(callback func()) error {
	return application.Trigger(callback)
}
