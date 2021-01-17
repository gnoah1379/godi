package godi

import (
	"github.com/gnoah1379/godi/internal"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

func (c *Container) startServices(group *sync.WaitGroup) {
	specialPackets := c.getSpecialPackets()
	for name, pack := range specialPackets {
		if c.ctx.Retrievable() && pack.canIsService() {
			services, err := pack.getServices()
			if err != nil {
				c.Log().Errorf("Can't get runner from packet %s because error: %v", name, err)
			} else {
				for tag, serv := range services {
					group.Add(1)
					go c.runService(group, serv, name, tag)
				}
			}
		}
	}
}

func (c *Container) startContainer() {
	c.ctx.start()
	c.ctx.setState(ContainerRunning)
	c.Log().Infof("Godi version: %s", internal.Version())
	c.Log().Infof("Starting %s using %s on %s with pid %d",
		filepath.Base(os.Args[0]),
		runtime.Version(),
		os.Getenv("USER"),
		os.Getpid())
}

func (c *Container) shutdownContainer() {
	c.ctx.setState(ContainerShutdown)
	c.Log().Infof("Godi Container shutdown finished")
}

func (c *Container) getSpecialPackets() map[string]*packet {
	specialPackets := make(map[string]*packet)
	c.packets.Range(func(key, value interface{}) bool {
		pack, name := value.(*packet), key.(string)
		if pack.canIsService() {
			specialPackets[name] = pack
		}
		return true
	})
	return specialPackets
}

func (c *Container) runService(group *sync.WaitGroup, service interface{}, name, tag string) {
	defer func() {
		if r := recover(); r != nil {
			c.Log().Warnf("Godi service %s with tag %s recovered exception: %v", name, tag, r)
		}
		group.Done()
	}()
	var err error
	ctx := newServiceContext(c.ctx)
	group.Add(1)
	go func() {
		defer group.Done()
		<-ctx.Done()
		shutDowner, ok := service.(ShutDowner)
		ctx.setState(ServiceShutDown)
		if ok {
			c.Log().Debugf("Godi service '%s' with tag '%s' shutting down", name, tag)
			shutDowner.ShutDown(ctx)
		}
	}()

	preloader, ok := service.(PreLoader)
	if ok {
		ctx.setState(ServicePreloading)
		c.Log().Debugf("Godi service '%s' with tag '%s' preloading", name, tag)
		preloader.Preload(ctx)
	}

	runner, ok := service.(Runner)
	if ok {
		c.Log().Debugf("Godi service '%s' with tag '%s' running", name, tag)
		ctx.setState(ServiceRunning)
		err = runner.Run(ctx)
		if err != nil {
			c.Log().Warnf("Godi service '%s' with tag '%s' exited: %v", name, tag, err)
		}
		ctx.cancel()
	}
}

func (c *Container) executeTrigger(group *sync.WaitGroup) {
	group.Add(1)
	go c.triggers.Execute(group)
}
