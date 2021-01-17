package godi

import (
	"github.com/gnoah1379/godi/internal"
	"reflect"
	"sync"

	"github.com/pkg/errors"
)

type provider struct {
	constructor interface{}
	filters     map[string][]string
	inTypes     []reflect.Type
	scope       provideScope
	types       provideType
	instance    reflect.Value
	once        sync.Once
}

func newProvider(constructor interface{}, scope provideScope, types provideType, filters map[string][]string) *provider {
	return &provider{
		constructor: constructor,
		filters:     filters,
		scope:       scope,
		types:       types,
	}
}

func (p *provider) getInstance(ctn *Container) (reflect.Value, error) {
	if p.scope == Singleton {
		var err error
		p.once.Do(func() {
			p.inTypes = internal.InputOf(p.constructor)
			p.instance, err = p.call(ctn)
		})
		return p.instance, err
	} else {
		return p.call(ctn)
	}
}

func (p *provider) inject(ctn *Container) ([]reflect.Value, error) {
	typesLen := len(p.inTypes)
	inParam := make([]reflect.Value, typesLen)
	var err error
	for i, inType := range p.inTypes {
		v := reflect.New(inType)
		switch inType.Kind() {
		case reflect.Slice:
			err = ctn.sliceOf(v, p.filters[inType.Elem().String()])
		case reflect.Map:
			err = ctn.mapOf(v, p.filters[inType.Elem().String()])
		case reflect.Interface:
			err = ctn.instanceOf(v, p.filters[inType.String()])
		case reflect.Ptr:
			err = ctn.instanceOf(v, p.filters[inType.Elem().String()])
		default:
			err = InvalidDependencies
		}
		if err != nil {
			return inParam, errors.Wrapf(err, "inject %s failed", v.Type().String())
		}
		inParam[i] = v.Elem()
	}
	return inParam, nil
}

func (p *provider) call(ctn *Container) (reflect.Value, error) {
	injectParam, err := p.inject(ctn)
	if err != nil {
		return reflect.Value{}, err
	}
	constructorResults := reflect.ValueOf(p.constructor).Call(injectParam)
	result := constructorResults[0]
	err, ok := constructorResults[1].Interface().(error)
	if !ok {
		err = nil
	}
	if err != nil {
		comp, _ := internal.GetOutputFieldType(p.constructor, 0)
		return result, errors.Wrapf(err, "initialize %s failed", comp.String())
	}
	return result, nil
}
