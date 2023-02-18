package dstruct

import "reflect"

type field struct {
	name              string
	tag               reflect.StructTag
	value             reflect.Value
	typ               reflect.Type
	pkgPath           string
	anonymous         bool
	jsonName          string
	ptrDepth          int
	fqn               string
	structIndex       int
	numberOfSubFields int
}

func (f field) GetFieldName() string {
	return f.name
}

func (f field) GetValue() any {
	return f.value.Interface()
}

func (f field) GetType() reflect.Type {
	return f.typ
}

func (f field) GetFieldFQName() string {
	return f.fqn
}

func (f field) GetTag(t string) string {
	return f.tag.Get(t)
}
