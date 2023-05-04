package dstruct

import "reflect"

type structField struct {
	name              string
	tag               reflect.StructTag
	value             reflect.Value
	typ               reflect.Type
	pkgPath           string
	anonymous         bool
	jsonName          string
	ptrDepth          int
	fqn               string
	structIndex       *int
	numberOfSubFields *int
}

func (f structField) GetFieldName() string {
	return f.name
}

func (f structField) GetValue() any {
	return f.value.Interface()
}

func (f structField) GetType() reflect.Type {
	return f.typ
}

func (f structField) GetFieldFQName() string {
	return f.fqn
}

func (f structField) GetTag(t string) string {
	return f.tag.Get(t)
}
