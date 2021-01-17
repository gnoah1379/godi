package internal

import (
	"reflect"

	"github.com/pkg/errors"
)

var (
	IsNotFunction = errors.New("is not function")
)

func InputOf(f interface{}) []reflect.Type {
	t := reflect.TypeOf(f)
	numIn := t.NumIn()
	result := make([]reflect.Type, numIn)
	for i := 0; i < numIn; i++ {
		result[i] = t.In(i)
	}
	return result
}

func OutputOf(f interface{}) []reflect.Type {
	t := reflect.TypeOf(f)
	numOut := t.NumOut()
	result := make([]reflect.Type, numOut)
	for i := 0; i < numOut; i++ {
		result[i] = t.Out(i)
	}
	return result
}

func GetOutputFieldType(f interface{}, idx int) (reflect.Type, error) {
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		return nil, IsNotFunction
	}
	numOut := t.NumOut()
	if idx < 0 || idx >= numOut {
		return nil, errors.New("out of range")
	}
	return t.Out(idx), nil
}

func IsPtrStructOrInterface(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Interface:
		return true
	case reflect.Ptr:
		return t.Elem().Kind() == reflect.Struct
	default:
		return false
	}
}

func IsError(t reflect.Type) bool {
	var err error
	return t.Implements(reflect.TypeOf(&err).Elem())
}
