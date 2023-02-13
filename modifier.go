package dstruct

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type FieldMap map[string]*Node[field]

type DynamicStructModifier interface {
	// Instance returns a copy of the struct
	Instance() any

	// New returns a pointer to the struct
	New() any

	// Get gets the value of the struct field `field` and returns an error if the field is not found
	Get(field string) (any, error)

	// Set sets the value of the struct field `field` and returns an error if the field is not found.
	//
	// The program will panic if the type of value does not match the type of the struct field `field`.
	Set(field string, value any) error

	// GetFields returns are map contaning all fields within a struct (including fields subfields)
	GetFields() map[string]field
}

type FieldModifier func(*field)

type DynamicStructModifierImpl struct {
	strct     any
	fieldMap  FieldMap
	fieldData map[string]field
	root      *Node[field]
}

var _ DynamicStructModifier = &DynamicStructModifierImpl{}

func newStruct(strct any, rootNode *Node[field]) *DynamicStructModifierImpl {
	dsm := &DynamicStructModifierImpl{
		strct:     strct,
		fieldMap:  make(FieldMap),
		fieldData: make(map[string]field),
		root:      rootNode,
	}
	dsm.createFieldMap(rootNode)
	return dsm
}

func (dm *DynamicStructModifierImpl) createFieldMap(rootNode *Node[field]) {
	// fmt.Println("F: ", rootNode.data.name, rootNode.children)
	// if rootNode.parent != nil {
	// 	fmt.Println("F@: ", rootNode.parent.data.name)
	// }
	for _, field := range rootNode.children {
		if field.data.name == "Age" {
			fmt.Println("AFL: ", rootNode.data.name, field.data.fqn)
		}
		dm.fieldMap[field.data.fqn] = field
		dm.fieldData[field.data.fqn] = *field.data
		dm.createFieldMap(field)
	}
}

func (dm *DynamicStructModifierImpl) New() any {
	return dm.strct
}
func (dm *DynamicStructModifierImpl) Instance() any {
	return getUnderlyingPointerValue(dm.strct)
}

func (dm *DynamicStructModifierImpl) Get(field string) (any, error) {
	if dm.fieldMap[field] == nil {
		return nil, fmt.Errorf("field %s does not exists in struct", field)
	}
	return dm.fieldMap[field].data.value.Interface(), nil
}

func (dm *DynamicStructModifierImpl) Set(field string, value any) error {
	if dm.fieldMap[field] == nil {
		return fmt.Errorf("field %s does not exists in struct", field)
	}
	fields := strings.Split(field, ".")
	c := fields[len(fields)-1][0]

	if 'a' <= c && c <= 'z' || c == '_' {
		return fmt.Errorf("field %s is not exported", field)
	}

	fieldValue := dm.fieldMap[field].data.value
	if !CanExtend(value) {
		if value == nil {
			if !fieldValue.IsZero() {
				fieldValue.Set(reflect.Zero(fieldValue.Type()))
			}
			return nil
		}
		fieldValue.Set(reflect.ValueOf(value))
		return nil
	}
	fieldValue.Set(reflect.ValueOf(value).Convert(fieldValue.Type()))

	return nil
}

func (dm *DynamicStructModifierImpl) GetFields() map[string]field {
	return dm.fieldData
}

func (dm *DynamicStructModifierImpl) String() string {
	val, _ := json.MarshalIndent(dm.strct, "", "\t")
	return string(val)
}
