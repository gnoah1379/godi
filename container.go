package godi

import (
	"github.com/gnoah1379/godi/internal"
	"reflect"
	"sync"
)

type Container struct {
	ctx      *containerContext
	packets  sync.Map
	triggers internal.AsyncJobQueue
}

func (c *Container) Context() ContainerContext {
	return c.ctx
}

func (c *Container) SetLogger(logger Logger) {
	if c.ctx.State() == Initializing {
		c.ctx.setLog(logger)
	}
}

func (c *Container) Log() Logger {
	return c.ctx.Log()
}

func (c *Container) Run() {
	wg := sync.WaitGroup{}
	c.startContainer()
	c.executeTrigger(&wg)
	c.startServices(&wg)
	<-c.ctx.Done()
	wg.Wait()
	c.shutdownContainer()
}

func (c *Container) Trigger(callback func()) error {
	if c.ctx.State() == ContainerRunning {
		return ContainerInRunning
	}
	c.triggers.Push(callback)
	return nil
}

func (c *Container) Register(constructor interface{}, opts []RegisterOption) error {
	options := newRegisterOption()
	for _, opt := range opts {
		opt(&options)
	}
	return c.register(constructor, options.scope, options.types, options.tag, options.filters)
}

// unregister all if tags = nil
func (c *Container) Unregister(packetName string, tags []string) {
	if tags == nil {
		c.packets.Delete(packetName)
	} else {
		pack, ok := c.packets.Load(packetName)
		if ok {
			pack.(*packet).remove(tags)
		}
	}
}

func (c *Container) InstanceOf(value interface{}, filters []string) error {
	v := reflect.ValueOf(value)
	if err := c.wantRetrieve(v); err != nil {
		return err
	}
	return c.instanceOf(v, filters)
}

func (c *Container) SliceOf(value interface{}, filters []string) error {
	v := reflect.ValueOf(value)
	if err := c.wantRetrieve(v); err != nil {
		return err
	}
	return c.sliceOf(v, filters)
}

func (c *Container) MapOf(value interface{}, filters []string) error {
	v := reflect.ValueOf(value)
	if err := c.wantRetrieve(v); err != nil {
		return err
	}
	return c.mapOf(v, filters)
}
