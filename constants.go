package godi

const (
	DefaultTag = "default"
)

const (
	Singleton provideScope = iota
	Prototype
)

const (
	Initializing ContainerState = iota
	ContainerRunning
	ContainerShutdown
)

const (
	Normal       provideType = 0
	Preloadable  provideType = 1
	Shutdownable provideType = 2
	Runnable     provideType = 4
	Service                  = Preloadable | Runnable | Shutdownable
)

const (
	ServiceInitialize ServiceState = 0
	ServicePreloading ServiceState = 1
	ServiceRunning    ServiceState = 2
	ServiceShutDown   ServiceState = 4
)
