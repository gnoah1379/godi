package godi

import (
	"reflect"
	"strings"
)

type provideScope uint8
type provideType uint8
type ContainerState uint
type ServiceState uint

func (p provideScope) Name() string {
	switch p {
	case Singleton:
		return "singleton"
	case Prototype:
		return "prototype"
	default:
		return "unknown"
	}
}

func (c ContainerState) Name() string {
	switch c {
	case Initializing:
		return "initialize"
	case ContainerRunning:
		return "Container running"
	//case ServicesPreloading:
	//	return "preloading"
	//case ServicesStarting:
	//	return "service starting"
	//case ServicesRunning:
	//	return "service running"
	//case ServicesShutDown:
	//	return "service shutDown"
	case ContainerShutdown:
		return "Container shutdown"
	default:
		return "unknown"
	}
}

func (p provideType) Name() string {
	var name = make([]string, 0)
	if p&Preloadable != 0 {
		name = append(name, "PreLoader")
	}

	if p&Shutdownable != 0 {
		name = append(name, "ShutDowner")
	}

	if p&Runnable != 0 {
		name = append(name, "Runner")
	}
	lenName := len(name)
	if lenName == 0 {
		return "Normal"
	} else {
		return strings.Join(name, ", ")
	}
}
func (p provideType) isValidPrototype(prototype reflect.Type) bool {
	if p&Preloadable > 0 && !prototype.Implements(reflect.TypeOf(new(PreLoader)).Elem()) {
		return false
	}

	if p&Shutdownable > 0 && !prototype.Implements(reflect.TypeOf(new(ShutDowner)).Elem()) {
		return false
	}

	if p&Runnable > 0 && !prototype.Implements(reflect.TypeOf(new(Runner)).Elem()) {
		return false
	}
	return true
}
