package godi

import (
	"context"
	"github.com/gnoah1379/godi/internal"
	"sync"
)

type ContainerContext interface {
	context.Context
	Log() Logger
	State() ContainerState
	Retrievable() bool
	ShutDown()
}

type containerContext struct {
	context.Context
	cancel context.CancelFunc
	state  ContainerState
	logger Logger
	locker sync.Locker
}

func (c *containerContext) Log() Logger {
	c.locker.Lock()
	log := c.logger
	c.locker.Unlock()
	return log
}

func (c *containerContext) State() ContainerState {
	c.locker.Lock()
	state := c.state
	c.locker.Unlock()
	return state
}

func (c *containerContext) ShutDown() {
	c.cancel()
}

func (c *containerContext) Retrievable() bool {
	return c.State() == ContainerRunning
}

func (c *containerContext) setState(phase ContainerState) {
	c.locker.Lock()
	c.state = phase
	c.locker.Unlock()
}

func (c *containerContext) setLog(log Logger) {
	c.locker.Lock()
	c.logger = log
	c.locker.Unlock()
}

func (c *containerContext) start() {
	c.Context, c.cancel = context.WithCancel(context.Background())
}

func newContainerContext() *containerContext {
	c := &containerContext{
		logger: defaultLogger,
		locker: internal.SpinLock(),
	}
	c.setState(Initializing)
	return c
}
