package dstruct

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type FieldMap map[string]*Node[Field]

type DynamicStructModifier interface {
	Instance() any
	New() any
	Get(field string) (any, error)
	Set(field string, value any) error
}

type FieldModifier func(*Field)

type DynamicStructModifierImpl struct {
	strct         any
	fieldMap      FieldMap
	fieldModifier FieldModifier
	root          *Node[Field]
}

var _ DynamicStructModifier = &DynamicStructModifierImpl{}

func newStruct(strct any, rootNode *Node[Field], fieldModifer FieldModifier) *DynamicStructModifierImpl {
	dsm := &DynamicStructModifierImpl{
		strct:         strct,
		fieldMap:      make(FieldMap),
		fieldModifier: fieldModifer,
		root:          rootNode,
	}
	dsm.createFieldMap(rootNode)
	return dsm
}

func (dm *DynamicStructModifierImpl) createFieldMap(rootNode *Node[Field]) {
	for _, field := range rootNode.children {
		dm.fieldMap[field.data.fqn] = field
		dm.createFieldMap(field)
		if !field.HasChildren() {
			if dm.fieldModifier != nil {
				dm.fieldModifier(field.data)
			}
		}
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
	return dm.fieldMap[field].data.Value.Interface(), nil
}

func (dm *DynamicStructModifierImpl) Set(field string, value any) error {
	if dm.fieldMap[field] == nil {
		return fmt.Errorf("field %s does not exists in struct", field)
	}

	fieldValue := dm.fieldMap[field].data.Value
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

func (dm *DynamicStructModifierImpl) String() string {
	val, _ := json.MarshalIndent(dm.strct, "", "\t")
	return string(val)
}
