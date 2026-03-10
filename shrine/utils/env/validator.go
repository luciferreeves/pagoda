package env

import (
	"errors"
	"reflect"
	"shrine/messages"
)

func validateConfigInput(config any) (reflect.Value, reflect.Type, error) {
	v := reflect.ValueOf(config)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return reflect.Value{}, nil, errors.New(messages.ConfigMustBePointer)
	}
	elem := v.Elem()
	return elem, elem.Type(), nil
}
