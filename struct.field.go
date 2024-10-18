package dstruct

import (
	"fmt"
	"reflect"
	"strconv"
)

type StructField struct {
	name              string
	tag               reflect.StructTag
	value             reflect.Value
	dstructType       reflect.Type
	goType            reflect.Type
	pkgPath           string
	anonymous         bool
	jsonName          string
	qualifiedName     string
	ptrDepth          int
	ptrKind           reflect.Kind
	structIndex       *int
	numberOfSubFields *int
}

// GetFieldName returns the name of the field.
func (f StructField) GetFieldName() string {
	return f.name
}

// GetValue returns the value of the field.
func (f StructField) GetValue() any {
	return f.value.Interface()
}

// GetType returns the Dstruct type of the field which can be different from the Go type (reflect.TypeOf(val))
func (f StructField) GetType() reflect.Type {
	return f.dstructType
}

// GetQualifiedName returns the fully qualified name of the field.
func (f StructField) GetQualifiedName() string {
	return f.qualifiedName
}

// GetTag returns the tag of the field.
func (f StructField) GetTag(t string) string {
	return f.tag.Get(t)
}

// GetJsonName returns the json name of the field.
func (f StructField) GetJsonName() string {
	return f.jsonName
}

// TODO: shoudl this even be here as it is not used ?
func (f StructField) GetEnumValues() (enumValues map[string]int) {
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
