package utils

import "reflect"

func KeySet(m interface{}) []reflect.Value {
	return reflect.ValueOf(m).MapKeys()
}