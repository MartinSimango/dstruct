package dstruct

import (
	"fmt"
	"reflect"
	"strconv"
)

type structField struct {
	name               string
	tag                reflect.StructTag
	value              reflect.Value
	typ                reflect.Type
	goType             reflect.Type
	pkgPath            string
	anonymous          bool
	jsonName           string
	ptrDepth           int
	fullyQualifiedName string
	structIndex        *int
	numberOfSubFields  *int
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

func (f structField) GetFieldFullyQualifiedName() string {
	return f.fullyQualifiedName
}

func (f structField) GetTag(t string) string {
	return f.tag.Get(t)
}

func (f structField) GetJsonName() string {
	return f.jsonName
}

func (f structField) GetEnumValues() (enumValues map[string]int) {
	enum, ok := f.tag.Lookup("enum")
	if ok {
		numEnums, err := strconv.Atoi(enum)
		if err != nil {
			return
		}
		enumValues = make(map[string]int)
		for i := 1; i <= numEnums; i++ {
			if key := f.tag.Get(fmt.Sprintf("enum_%d", i)); key != "" {
				enumValues[key] = i
			}
		}
	}
	return
}
