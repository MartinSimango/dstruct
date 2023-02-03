package dstruct

import (
	"reflect"
)

func getPointerToInterface(str any) any {
	return reflect.New(reflect.ValueOf(str).Type()).Interface()
}

func getUnderlyingPointerValue(ptr any) any {
	return reflect.ValueOf(ptr).Elem().Interface()
}

func getSliceType(value reflect.Value) reflect.Type {
	return reflect.TypeOf(value.Interface()).Elem()
}

func GetPointerToSliceType(sliceType reflect.Type) any {
	return reflect.New(sliceType).Elem().Interface()
}
