package internal

import (
	"reflect"
)

func IsGodiConstructor(constructor interface{}) bool {
	t := reflect.TypeOf(constructor)
	if t.Kind() != reflect.Func {
		return false
	}

	numIn := t.NumIn()
	numOut := t.NumOut()

	for i := 0; i < numIn; i++ {
		typeIn := t.In(i)
		typeInValid := false
		switch typeIn.Kind() {
		case reflect.Slice:
			typeInValid = IsPtrStructOrInterface(typeIn.Elem())

		case reflect.Map:
			typeInValid = IsPtrStructOrInterface(typeIn.Elem())

		default:
			typeInValid = IsPtrStructOrInterface(typeIn)
		}
		if !typeInValid {
			return false
		}
	}

	if numOut != 2 {
		return false
	}

	tOut := t.Out(0)

	typeOutValid := false
	switch tOut.Kind() {
	//case reflect.Slice:
	//	typeOutValid = IsPtrStructOrInterface(tOut.Elem())
	//
	//case reflect.Map:
	//	typeOutValid = IsPtrStructOrInterface(tOut.Elem())

	default:
		typeOutValid = IsPtrStructOrInterface(tOut)
	}
	if !typeOutValid {
		return false
	}
	return IsError(t.Out(1))
}

func GetGodiPacketName(t reflect.Type) (string, bool) {
	k := t.Kind()
	switch k {
	case reflect.Ptr:
		return t.Elem().String(), true
	case reflect.Interface:
		return t.String(), true
	default:
		return "", false
	}
}
