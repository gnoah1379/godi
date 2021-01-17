package godi

import (
	"godi/internal"
	"reflect"
	"sync"
)

type packet struct {
	ctn       *Container
	prototype reflect.Type
	providers sync.Map
}

func newPacket(ctn *Container, iType reflect.Type) *packet {
	return &packet{
		ctn:       ctn,
		prototype: iType,
	}
}

func (p *packet) retrieveInstance(value reflect.Value, tags []string) error {
	prv, err := p.firstProviderMatched(tags)
	if err != nil {
		return err
	}
	instance, err := prv.getInstance(p.ctn)
	if err != nil {
		return err
	}
	value.Elem().Set(instance)
	return nil
}

func (p *packet) retrieveSlice(value reflect.Value, tags []string) error {
	matchedInstances := reflect.MakeSlice(value.Type().Elem(), 0, 0)
	for _, prv := range p.getProviders(tags) {
		instance, err := prv.getInstance(p.ctn)
		if err != nil {
			return err
		}
		matchedInstances = reflect.Append(matchedInstances, instance)
	}
	value.Elem().Set(matchedInstances)
	return nil
}

func (p *packet) retrieveMap(value reflect.Value, tags []string) error {
	matchedInstances := reflect.MakeMap(value.Type().Elem())
	for tag, prv := range p.getProviders(tags) {
		instance, err := prv.getInstance(p.ctn)
		if err != nil {
			return err
		}
		matchedInstances.SetMapIndex(reflect.ValueOf(tag), instance)
	}
	value.Elem().Set(matchedInstances)
	return nil
}

func (p *packet) provide(tag string, src *provider) {
	p.providers.Store(tag, src)
}

func (p *packet) remove(tags []string) {
	for _, tag := range tags {
		p.providers.Delete(tag)
	}
}

func (p *packet) getProviders(tags []string) map[string]*provider {
	result := make(map[string]*provider)
	getAll := tags == nil
	p.providers.Range(func(key, value interface{}) bool {
		keyStr := key.(string)
		if getAll || internal.KeyInStrings(keyStr, tags) {
			result[keyStr] = value.(*provider)
		}
		return true
	})
	return result
}

func (p *packet) getProvider(tag string) (*provider, error) {
	prv, ok := p.providers.Load(tag)
	if !ok {
		return nil, ProviderNotFound
	}
	return prv.(*provider), nil
}

func (p *packet) firstProviderMatched(tags []string) (*provider, error) {
	for _, tag := range tags {
		prv, err := p.getProvider(tag)
		if err == nil {
			return prv, nil
		}
	}
	return nil, ProviderNotFound
}

func (p *packet) canIsService() bool {
	implPrototype := p.prototype
	if implPrototype.Kind() == reflect.Struct {
		implPrototype = reflect.PtrTo(implPrototype)
	}
	return implPrototype.Implements(reflect.TypeOf(new(PreLoader)).Elem()) ||
		implPrototype.Implements(reflect.TypeOf(new(Runner)).Elem()) ||
		implPrototype.Implements(reflect.TypeOf(new(ShutDowner)).Elem())
}

func (p *packet) getServices() (services map[string]interface{}, err error) {
	services = make(map[string]interface{})
	p.providers.Range(func(key, value interface{}) bool {
		if err == nil {
			prv, tag := value.(*provider), key.(string)
			if prv.scope == Singleton && prv.types&Service > 0 {
				var instance reflect.Value
				instance, err = prv.getInstance(p.ctn)
				if err != nil {
					return false
				}
				services[tag] = instance.Interface()
			}
		}
		return true
	})
	return services, err
}
