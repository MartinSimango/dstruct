package dstruct

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type treeBuilderImpl struct {
	root      *Node[Field]
	setValues bool
}

var _ Builder = &treeBuilderImpl{}

func NewBuilder() *treeBuilderImpl {
	return newBuilderFromNode(&Node[Field]{
		data: &Field{
			Name:  "#root",
			Value: reflect.ValueOf(nil),
		},
		children: make(map[string]*Node[Field]),
	},
	)
}

func CanExtend(val any) bool {
	if val == nil {
		return false
	}
	ptrValue, _ := getPtrValue(reflect.ValueOf(val), 0)
	return ptrValue.Type().Kind() == reflect.Struct
}

func ExtendStruct(val any) *treeBuilderImpl {
	// TODO check if val is a struct
	b := NewBuilder()
	value := reflect.ValueOf(val)

	if !CanExtend(val) {
		panic(fmt.Sprintf("Cannot extend struct value of type %s", value.Type()))
	}

	switch value.Kind() {
	case reflect.Struct:
		b.addStructFields(value, b.root, 0)
	case reflect.Ptr:
		b.addPtrField(value, b.root)
	}

	return b

}
func newBuilderFromNode(node *Node[Field]) *treeBuilderImpl {

	return &treeBuilderImpl{
		setValues: true,
		root:      node,
	}

}

func (dsb *treeBuilderImpl) AddField(name string, typ interface{}, tag string) Builder {
	dsb.addFieldToTree(name, typ, reflect.StructTag(tag), dsb.root)
	return dsb
}

func (dsb *treeBuilderImpl) RemoveField(name string) Builder {
	fields := strings.Split(name, ".")
	node := dsb.root

	for i := 0; i < len(fields)-1; i++ {
		node = node.GetNode(fields[i])
	}
	node.DeleteNode(fields[len(fields)-1])

	return dsb
}

func (dsb *treeBuilderImpl) GetField(field string) Builder {

	node := dsb.getNode(field)
	// TODO validate node exists
	return newBuilderFromNode(node)
}

func (dsb *treeBuilderImpl) GetFieldCopy(field string) Builder {
	copyNode := dsb.getNode(field).Copy()
	return newBuilderFromNode(copyNode)
}

func (dsb *treeBuilderImpl) getNode(field string) *Node[Field] {

	fields := strings.Split(field, ".")
	node := dsb.root

	for i := 0; i < len(fields); i++ {
		if node = node.GetNode(fields[i]); node == nil {
			return nil
		}
	}
	return node

}

func (db *treeBuilderImpl) Build() DynamicStructModifier {
	return db.BuildWithFieldModifier(nil)
}

func (db *treeBuilderImpl) BuildWithFieldModifier(fieldModifier FieldModifier) DynamicStructModifier {
	return newStruct(db.buildStruct(db.root), db.root.Copy(), fieldModifier)
}

func (db *treeBuilderImpl) buildStruct(tree *Node[Field]) any {
	strctValue := reflect.ValueOf(getPointerToInterface(treeToStruct(tree)))
	tree.data.Value = strctValue
	if db.setValues {
		if strctValue.Elem().Kind() == reflect.Ptr {
			setPointerFieldValue(strctValue.Elem(), tree)
		} else {
			setStructFieldValues(strctValue.Elem(), tree)
		}
	}

	return strctValue.Interface()
}

func (dsb *treeBuilderImpl) addFieldToTree(name string, typ interface{}, tag reflect.StructTag, root *Node[Field]) reflect.Type {
	value := reflect.ValueOf(typ)
	if !value.IsValid() {
		panic(fmt.Sprintf("Cannot determine type of %s", name))
	}

	field := &Field{
		Name:        name,
		Value:       value,
		Tag:         tag,
		Type:        reflect.TypeOf(value.Interface()),
		jsonName:    strings.Split(tag.Get("json"), ",")[0],
		StructIndex: root.data.SubFields,
	}
	field.fqn = getFQN(root.data.GetFieldName(), field.GetFieldName())

	root.AddNode(name, field)
	root.data.SubFields++

	switch value.Kind() {
	case reflect.Struct:
		field.Type = dsb.addStructFields(value, root.children[name], 0)
	case reflect.Ptr:
		field.Type = dsb.addPtrField(value, root.children[name])
	}

	return field.Type
}

func sortKeys(root *Node[Field]) (keys []string) {
	for key := range root.children {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return root.children[keys[i]].data.StructIndex < root.children[keys[j]].data.StructIndex
	})

	return
}

func treeToStruct(root *Node[Field]) any {
	var structFields []reflect.StructField

	//sort the keys to ensure type  of struct produced is always the same
	var keys []string = sortKeys(root)

	for _, fieldName := range keys {
		var typ reflect.Type
		field := root.GetNode(fieldName)
		if field.HasChildren() {
			typ = reflect.TypeOf(treeToStruct(field))
		} else {
			typ = field.data.Value.Type()
		}

		structFields = append(structFields, reflect.StructField{
			Name:      fieldName,
			PkgPath:   "",
			Type:      typ,
			Tag:       field.data.Tag,
			Anonymous: false,
		})

	}

	strct := reflect.New(reflect.StructOf(structFields)).Elem()
	for i := 0; i < root.data.ptrDepth; i++ {
		strct = reflect.New(reflect.TypeOf(strct.Interface()))
	}

	return strct.Interface()
}

func setStructFieldValues(strct reflect.Value, root *Node[Field]) {
	for i := 0; i < strct.NumField(); i++ {
		field := strct.Field(i)
		fieldName := strct.Type().Field(i).Name
		currentNode := root.GetNode(fieldName)
		switch field.Kind() {
		case reflect.Struct:
			setStructFieldValues(field, currentNode)
		case reflect.Ptr:
			setPointerFieldValue(field, currentNode)
		default:
			field.Set(currentNode.data.Value)
		}
		currentNode.data.Value = field

	}

}

func setPointerFieldValue(field reflect.Value, currentNode *Node[Field]) {
	if currentNode.data.Value.IsNil() {
		return
	}

	f := field
	if currentNode.HasChildren() { // node is a struct that needs to be derefenced
		for i := 0; i < currentNode.data.ptrDepth; i++ {
			f.Set(reflect.New(f.Type().Elem()))
			f = f.Elem()
		}
	}

	switch f.Kind() {
	case reflect.Struct:
		setStructFieldValues(f, currentNode)
	default:
		field.Set(currentNode.data.Value)
	}
	currentNode.data.Value = field

}

func (dsb *treeBuilderImpl) addStructFields(strct reflect.Value, root *Node[Field], ptrDepth int) reflect.Type {
	var structFields []reflect.StructField

	for i := 0; i < strct.NumField(); i++ {
		fieldName := strct.Type().Field(i).Name
		fieldTag := strct.Type().Field(i).Tag
		fieldType := dsb.addFieldToTree(fieldName, strct.Field(i).Interface(), fieldTag, root)

		structFields = append(structFields, reflect.StructField{
			Name:      fieldName,
			PkgPath:   "",
			Type:      fieldType,
			Tag:       fieldTag,
			Anonymous: false,
		})

	}
	retStruct := reflect.New(reflect.StructOf(structFields)).Elem()
	for i := 0; i < ptrDepth; i++ {
		retStruct = reflect.New(retStruct.Type())
	}
	return reflect.TypeOf(retStruct.Interface())
}

func getPtrValue(value reflect.Value, ptrDepth int) (reflect.Value, int) {
	switch value.Kind() {
	case reflect.Ptr:
		return getPtrValue(value.Elem(), ptrDepth+1)
	}
	return value, ptrDepth
}

func (dsb *treeBuilderImpl) addPtrField(value reflect.Value, node *Node[Field]) reflect.Type {

	if value.IsNil() {
		return reflect.TypeOf(value.Interface())
	}

	ptrValue, ptrDepth := getPtrValue(value, 0)
	node.data.ptrDepth = ptrDepth

	switch ptrValue.Kind() {
	case reflect.Struct:
		return dsb.addStructFields(ptrValue, node, ptrDepth)
	}
	return reflect.TypeOf(value.Interface())
}

func getFQN(root, name string) string {
	if root != "#root" {
		return fmt.Sprintf("%s.%s", root, name)
	}
	return name
}
