package dstruct

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/MartinSimango/dstruct/dreflect"
)

type FieldNode map[string]*Node[StructField]

type FieldData map[string]StructField

type DynamicStructModifier interface {
	// Instance returns a copy of the struct
	Instance() any

	// New returns a pointer to the struct
	New() any

	// Get gets the value of the struct field `field` and returns an error if the field is not found
	Get(field string) (any, error)

	// Get_ gets the value of the struct field `field` and panics if the field is not found
	Get_(field string) any

	// Set sets the value of the struct field `field` and returns an error if the field is not found.
	//
	// The program will panic if the type of value does not match the type of the struct field `field`.
	Set(field string, value any) error

	// GetFields returns a map containing all fields within a struct
	GetFields() map[string]StructField

	// Update updates the struct's underlying tree to represent that of the strct's value.
	// The structs underlying tree can change if new fields are added due to fields within the struct changing from
	// nil to become not nil. This can lead to new additional fields being introduced within the struct
	Update()

	// Apply is a combination of Set and Update. Update is not called if Apply fails.
	Apply(field string, value any) error
}

type FieldModifier func(*StructField)

type DynamicStructModifierImpl struct {
	strct        any
	fieldNodeMap FieldNode
	fieldData    map[string]StructField
	root         *Node[StructField]
}

var _ DynamicStructModifier = &DynamicStructModifierImpl{}

func newStruct(strct any, rootNode *Node[StructField]) *DynamicStructModifierImpl {
	dsm := &DynamicStructModifierImpl{
		strct:        strct,
		fieldNodeMap: make(FieldNode),
		fieldData:    make(map[string]StructField),
		root:         rootNode,
	}
	dsm.createFieldToNodeMappings(rootNode)
	return dsm
}

func (dm *DynamicStructModifierImpl) createFieldToNodeMappings(rootNode *Node[StructField]) {
	for _, field := range rootNode.children {
		dm.fieldNodeMap[field.data.qualifiedName] = field
		dm.fieldData[field.data.qualifiedName] = *field.data
		dm.createFieldToNodeMappings(field)
	}
}

func (dm *DynamicStructModifierImpl) New() any {
	return dm.strct
}

func (dm *DynamicStructModifierImpl) Instance() any {
	return dreflect.GetUnderlyingPointerValue(dm.strct)
}

func (dm *DynamicStructModifierImpl) get(field string) (n *Node[StructField]) {
	return dm.fieldNodeMap[field]
}

func (dm *DynamicStructModifierImpl) Get(field string) (any, error) {
	if f := dm.get(field); f == nil {
		return nil, fmt.Errorf("field %s does not exists in struct", field)
	} else {
		return f.data.value.Interface(), nil
	}
}

func (dm *DynamicStructModifierImpl) Get_(field string) any {
	if f := dm.get(field); f == nil {
		panic(fmt.Errorf("field %s does not exists in struct", field))
	} else {
		return f.data.value.Interface()
	}
}

func isFieldExported(field string) bool {
	fields := strings.Split(field, ".")

	for _, f := range fields {
		c := f[0]

		if 'a' <= c && c <= 'z' || c == '_' {
			return false
		}
	}
	return true
}

func (dm *DynamicStructModifierImpl) Set(field string, value any) error {
	var f *Node[StructField]
	if f = dm.get(field); f == nil {
		return fmt.Errorf("field %s does not exists in struct", field)
	}

	if !isFieldExported(field) {
		return fmt.Errorf("field %s is not exported", field)
	}
	fieldValue := f.data.value
	if !canExtend(value) { // if value is not a struct
		if value == nil {
			if !fieldValue.IsZero() {
				fieldValue.Set(reflect.Zero(fieldValue.Type()))
			}
			return nil
		}
		fieldValue.Set(dreflect.Convert(reflect.ValueOf(value), fieldValue.Type()))

		return nil
	}

	fieldValue.Set(dreflect.Convert(reflect.ValueOf(value), fieldValue.Type()))
	dm.update(f)

	return nil
}

// update updates the struct's underlying node to represent that of the strct's value.
func (dm *DynamicStructModifierImpl) update(node *Node[StructField]) {
	nodeValue := node.data.value
	if canExtend(nodeValue.Interface()) {
		node.children = ExtendStruct(nodeValue.Interface()).root.children
		setValues(nodeValue, resetNodeFieldsFQN(node))
		dm.createFieldToNodeMappings(node)
	}
}

func (dm *DynamicStructModifierImpl) GetFields() map[string]StructField {
	return dm.fieldData
}

func (dm *DynamicStructModifierImpl) String() string {
	val, _ := json.MarshalIndent(dm.strct, "", "\t")
	return string(val)
}

func (dm *DynamicStructModifierImpl) Update() {
	*dm = *ExtendStruct(dm.Instance()).Build().(*DynamicStructModifierImpl)
}

func (dm *DynamicStructModifierImpl) Apply(field string, value any) error {
	if err := dm.Set(field, value); err != nil {
		return err
	}
	dm.Update()
	return nil
}
