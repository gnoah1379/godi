package godi

import "reflect"

func TypeOf(value interface{}) string {
	return reflect.TypeOf(value).String()
}

func TypeOfPtr(value interface{}) string {
	return reflect.TypeOf(value).Elem().String()
}
