package core

import "reflect"

var customKind map[reflect.Type]*reflect.Kind = make(map[reflect.Type]*reflect.Kind)
var latestKind = reflect.UnsafePointer

func NewKind(val any) reflect.Kind {
	customKindType := reflect.TypeOf(val)
	if customKind[customKindType] == nil {
		customKind[customKindType] = new(reflect.Kind)
		latestKind++
		*customKind[customKindType] = latestKind
		return latestKind
	}
	return *customKind[customKindType]
}
