package dstruct

import "reflect"

type Field struct {
	Name        string
	Tag         reflect.StructTag
	Value       reflect.Value
	Type        reflect.Type
	jsonName    string
	ptrDepth    int
	fqn         string
	StructIndex int
	SubFields   int
}

func (f Field) GetFieldName() string {
	if f.jsonName == "" {
		return f.Name
	}
	return f.jsonName
}

func (f Field) GetFieldFQName() string {
	return f.fqn
}
