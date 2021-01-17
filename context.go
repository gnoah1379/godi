package godi

import (
	"context"
)

type Context interface {
	context.Context
	Container() ContainerContext
	State() ServiceState
	Log() Logger
}

type serviceContext struct {
	context.Context
	ctx    *containerContext
	cancel context.CancelFunc
	state  ServiceState
}

func (c *serviceContext) Container() ContainerContext {
	return c.ctx
}

func (c *serviceContext) State() ServiceState {
	return c.state
}

func (c *serviceContext) Log() Logger {
	return c.ctx.Log()
}

func (c *serviceContext) setState(state ServiceState) {
	c.state = state
}

func newServiceContext(ctx *containerContext) *serviceContext {
	c := &serviceContext{
		ctx:   ctx,
		state: ServiceInitialize,
	}
	c.Context, c.cancel = context.WithCancel(ctx)
	return c
}
