package godi

import (
	"godi/internal"
	"reflect"

	"github.com/pkg/errors"
)

func (c *Container) wantRetrieve(v reflect.Value) error {
	if !c.ctx.Retrievable() {
		return ContainerIsNotRunning
	}
	if v.Type().Kind() != reflect.Ptr {
		return errors.Wrapf(VariableIsNotPtr, "'%s'", v.Type().Name())
	}
	return nil
}

func (c *Container) instanceOf(value reflect.Value, filters []string) error {
	packetType := value.Type().Elem()
	name, ok := internal.GetGodiPacketName(packetType)
	if !ok {
		return errors.Wrapf(ValueIsNotPtrOrInterface, "instance of '%s'", value.String())
	}
	packets, ok := c.packets.Load(name)
	if !ok {
		return errors.Wrapf(PacketNotFound, "'%s'", name)
	}

	if len(filters) == 0 {
		filters = []string{DefaultTag}
	}
	return packets.(*packet).retrieveInstance(value, filters)
}

func (c *Container) sliceOf(value reflect.Value, filters []string) error {
	sliceType := value.Type().Elem()
	if sliceType.Kind() != reflect.Slice {
		return errors.Wrapf(VariableIsNotSlice, "slice of %s", sliceType.String())
	}
	packetType := sliceType.Elem()
	name, ok := internal.GetGodiPacketName(packetType)
	if !ok {
		return errors.Wrapf(ValueIsNotPtrOrInterface, "slice elements of '%s'", sliceType.String())
	}
	cpn, ok := c.packets.Load(name)
	if !ok {
		return errors.Wrapf(PacketNotFound, "'%s'", name)
	}
	return cpn.(*packet).retrieveSlice(value, filters)
}

func (c *Container) mapOf(value reflect.Value, filters []string) error {
	mapType := value.Type().Elem()
	if mapType.Kind() != reflect.Map {
		return errors.Wrapf(VariableIsNotMap, "map of '%s'", mapType.String())
	}
	packetType := mapType.Elem()
	name, ok := internal.GetGodiPacketName(packetType)
	if !ok {
		return errors.Wrapf(ValueIsNotPtrOrInterface, "map elements of '%s'", mapType.String())
	}
	cpn, ok := c.packets.Load(name)
	if !ok {
		return errors.Wrapf(PacketNotFound, "'%s'", name)
	}
	return cpn.(*packet).retrieveMap(value, filters)
}
