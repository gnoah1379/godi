package godi

import "godi/internal"

var application *Container

func GetApplicationContainer() *Container {
	return application
}

func init() {
	application = &Container{
		ctx:      newContainerContext(),
		triggers: internal.NewAsyncJobQueue(),
	}
}
