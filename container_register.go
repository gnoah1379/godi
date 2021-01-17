package godi

import (
	"godi/internal"
	"reflect"

	"github.com/pkg/errors"
)

func (c *Container) register(constructor interface{}, scope provideScope, types provideType, tag string, filters map[string][]string) error {
	if !internal.IsGodiConstructor(constructor) {
		return errors.Wrapf(InvalidConstructor, "'%s'", reflect.TypeOf(constructor).String())
	}
	prototype, _ := internal.GetOutputFieldType(constructor, 0)
	if !types.isValidPrototype(prototype) {
		return errors.Wrapf(InvalidConstructor, "'%s' not impliments '%s'", prototype.String(), types.Name())
	}
	if prototype.Kind() == reflect.Ptr {
		prototype = prototype.Elem()
	}
	prototypeName := prototype.String()

	cpn, ok := c.packets.Load(prototypeName)
	if !ok {
		cpn = newPacket(c, prototype)
		c.packets.Store(prototype.String(), cpn)
	}
	cpn.(*packet).provide(tag, newProvider(constructor, scope, types, filters))
	return nil
}
