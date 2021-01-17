package godi

import "github.com/gnoah1379/godi/internal"

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
