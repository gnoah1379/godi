package godi

import "github.com/pkg/errors"

var (
	PacketNotFound           = errors.New("packet not found")
	VariableIsNotPtr         = errors.New("retrieve variable must be pointer")
	VariableIsNotSlice       = errors.New("retrieve variable must be slice")
	VariableIsNotMap         = errors.New("retrieve variable must be map")
	ValueIsNotPtrOrInterface = errors.New("value must be pointer or interface")
	InvalidConstructor       = errors.New("invalid godi constructor")
	InvalidDependencies      = errors.New("invalid dependencies")
	ContainerIsNotRunning    = errors.New("Container not running")
	ContainerInRunning       = errors.New("Container running")
	ProviderNotFound         = errors.New("provider not found")
)
