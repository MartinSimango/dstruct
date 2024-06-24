package core

import (
	"reflect"
)

var (
	customKind map[reflect.Type]*reflect.Kind = make(map[reflect.Type]*reflect.Kind)
	latestKind                                = reflect.UnsafePointer
)

// NewKind returns a new reflect.Kind for a custom type
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
