package dstruct

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"unsafe"
)

type treeBuilderImpl struct {
	root      *Node[field]
	setValues bool
}

var _ Builder = &treeBuilderImpl{}

func creatRoot() *Node[field] {
	return &Node[field]{
		data: &field{
			name:  "#root",
			value: reflect.ValueOf(nil),
			fqn:   "#root",
		},
		children: make(map[string]*Node[field]),
	}
}
func NewBuilder() Builder {
	return newBuilderFromNode(creatRoot(), false)
}

func CanExtend(val any) bool {
	if val == nil {
		return false
	}
	ptrValue, _ := getPtrValue(reflect.ValueOf(val), 0)
	return ptrValue.Type().Kind() == reflect.Struct
}

func ExtendStruct(val any) Builder {
	// TODO check if val is a struct
	b := NewBuilder().(*treeBuilderImpl)
	value := reflect.ValueOf(val)

	if !CanExtend(val) {
		panic(fmt.Sprintf("Cannot extend struct value of type %s", value.Type()))
	}

	switch value.Kind() {
	case reflect.Struct:
		b.addStructFields(value, b.root, 0, false)
	case reflect.Ptr:
		b.addPtrField(value, b.root, false)
	}

	return b

}
func newBuilderFromNode(node *Node[field], resetFQN bool) *treeBuilderImpl {

	if resetFQN {
		resetNodeFieldsFQN(node)
	}
	return &treeBuilderImpl{
		setValues: true,
		root:      node,
	}
}

func resetNodeFieldsFQN(node *Node[field]) {
	for _, v := range node.children {
		v.data.fqn = getFQN(node.data.name, v.data.name)
		resetNodeFieldsFQN(v)
	}
}

func (dsb *treeBuilderImpl) AddField(name string, value interface{}, tag string) Builder {
	dsb.addFieldToTree(name, value, "", false, reflect.StructTag(tag), dsb.root)
	return dsb
}

func (dsb *treeBuilderImpl) AddEmbeddedField(value interface{}, tag string) Builder {
	ptrValue, _ := getPtrValue(reflect.ValueOf(value), 0)
	name := reflect.TypeOf(ptrValue.Interface()).Name()
	dsb.addFieldToTree(name, value, "", true, reflect.StructTag(tag), dsb.root)

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
	if node == nil {
		return nil
	}
	return newBuilderFromNode(node, false)
}

func (dsb *treeBuilderImpl) NewBuilderFromField(field string) Builder {
	copyNode := dsb.getNode(field).Copy()
	copyNode.data.name = "#root"
	copyNode.data.fqn = "#root"
	return newBuilderFromNode(copyNode, true)
}

func (dsb *treeBuilderImpl) GetFieldCopy(field string) Builder {
	copyNode := dsb.getNode(field).Copy()
	return newBuilderFromNode(copyNode, false)
}

func (dsb *treeBuilderImpl) getNode(field string) *Node[field] {

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
	return newStruct(db.buildStruct(db.root), db.root.Copy())
}

func (db *treeBuilderImpl) buildStruct(tree *Node[field]) any {
	structValue := reflect.ValueOf(getPointerToInterface(treeToStruct(tree)))
	tree.data.value = structValue
	if db.setValues {
		if structValue.Elem().Kind() == reflect.Ptr {
			setPointerFieldValue(structValue.Elem(), tree)
		} else {
			setStructFieldValues(structValue.Elem(), tree)
		}
	}

	return structValue.Interface()
}

func (dsb *treeBuilderImpl) addFieldToTree(name string, typ interface{}, pkgPath string, anonymous bool, tag reflect.StructTag, root *Node[field]) reflect.Type {
	value := reflect.ValueOf(typ)
	if !value.IsValid() {
		panic(fmt.Sprintf("Cannot determine type of %s", name))
	}

	field := &field{
		name:      name,
		value:     value,
		tag:       tag,
		typ:       reflect.TypeOf(value.Interface()),
		pkgPath:   pkgPath,
		anonymous: anonymous,

		jsonName:    strings.Split(tag.Get("json"), ",")[0],
		structIndex: root.data.numberOfSubFields,
	}
	field.fqn = getFQN(root.data.GetFieldFQName(), field.name)

	root.AddNode(name, field)
	root.data.numberOfSubFields++

	switch value.Kind() {
	case reflect.Struct:
		field.typ = dsb.addStructFields(value, root.children[name], 0, anonymous)
	case reflect.Ptr:
		field.typ = dsb.addPtrField(value, root.children[name], anonymous)
	}

	return field.typ
}

func sortKeys(root *Node[field]) (keys []string) {
	for key := range root.children {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return root.children[keys[i]].data.structIndex < root.children[keys[j]].data.structIndex
	})

	return
}

func treeToStruct(root *Node[field]) any {
	var structFields []reflect.StructField

	//sort the keys to ensure type  of struct produced is always the same
	var keys []string = sortKeys(root)

	for _, fieldName := range keys {
		var typ reflect.Type
		field := root.GetNode(fieldName)
		if field.HasChildren() {
			typ = reflect.TypeOf(treeToStruct(field))
		} else {
			typ = field.data.value.Type()
		}

		structFields = append(structFields, reflect.StructField{
			Name:      fieldName,
			PkgPath:   field.data.pkgPath,
			Type:      typ,
			Tag:       field.data.tag,
			Anonymous: field.data.anonymous,
		})

	}

	strct := reflect.New(reflect.StructOf(structFields)).Elem()
	for i := 0; i < root.data.ptrDepth; i++ {
		strct = reflect.New(reflect.TypeOf(strct.Interface()))
	}

	return strct.Interface()
}

func setStructFieldValues(strct reflect.Value, root *Node[field]) {
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
			// field.Set(currentNode.data.value)
			// Support for setting values from non exported fields
			reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
				Elem().
				Set(currentNode.data.value)
		}
		currentNode.data.value = field

		if currentNode.data.anonymous {
			addAnonymousSubfields(currentNode)
		}

	}

}

func addAnonymousSubfields(anonymousNode *Node[field]) {

	parent := anonymousNode.parent
	for parent != nil {
		// Add anonymous node to parent
		if parent.children[anonymousNode.data.name] == nil {
			copyNode := anonymousNode.Copy()
			copyNode.data.fqn = getFQN(parent.data.name, copyNode.data.name)
			resetNodeFieldsFQN(copyNode)
			parent.children[anonymousNode.data.name] = copyNode
		}
		// Add anonymous node children to parent
		for k, v := range anonymousNode.children {
			copyNode := v.Copy()
			copyNode.data.fqn = getFQN(parent.data.name, copyNode.data.name)
			resetNodeFieldsFQN(copyNode)
			parent.children[k] = copyNode
		}

		parent = parent.parent
	}

}

func setPointerFieldValue(field reflect.Value, currentNode *Node[field]) {
	if currentNode.data.value.IsNil() {
		return
	}

	f := field
	if currentNode.HasChildren() { // node is a struct that needs to be dereferenced
		for i := 0; i < currentNode.data.ptrDepth; i++ {
			f.Set(reflect.New(f.Type().Elem()))
			f = f.Elem()
		}
	}

	switch f.Kind() {
	case reflect.Struct:
		setStructFieldValues(f, currentNode)
	default:
		field.Set(currentNode.data.value)
	}
	currentNode.data.value = field
	if currentNode.data.anonymous {
		addAnonymousSubfields(currentNode)
	}

}

func (dsb *treeBuilderImpl) addStructFields(strct reflect.Value, root *Node[field], ptrDepth int, anon bool) reflect.Type {
	var structFields []reflect.StructField

	for i := 0; i < strct.NumField(); i++ {
		fieldName := strct.Type().Field(i).Name
		fieldTag := strct.Type().Field(i).Tag
		var fieldValue any
		if strct.Type().Field(i).IsExported() {
			fieldValue = strct.Field(i).Interface()
		} else {
			// fieldValue = reflect.NewAt(strct.Field(i).Type(), unsafe.Pointer(root.data.value.Field(i).UnsafeAddr())).Elem().Interface()
			// reflect.New(strct.Field(i).Type()).Elem().Interface()
			fieldValue = reflect.New(strct.Field(i).Type()).Elem().Interface()

		}
		pkgPath := strct.Type().Field(i).PkgPath
		anonymous := strct.Type().Field(i).Anonymous
		fieldType := dsb.addFieldToTree(fieldName, fieldValue, pkgPath, anonymous, fieldTag, root)

		structFields = append(structFields, reflect.StructField{
			Name:      fieldName,
			PkgPath:   pkgPath,
			Type:      fieldType,
			Tag:       fieldTag,
			Anonymous: anonymous,
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

func (dsb *treeBuilderImpl) addPtrField(value reflect.Value, node *Node[field], anonymous bool) reflect.Type {

	if value.IsNil() {
		return reflect.TypeOf(value.Interface())
	}

	ptrValue, ptrDepth := getPtrValue(value, 0)
	node.data.ptrDepth = ptrDepth

	switch ptrValue.Kind() {
	case reflect.Struct:
		return dsb.addStructFields(ptrValue, node, ptrDepth, anonymous)
	}
	return reflect.TypeOf(value.Interface())
}

func getFQN(root, name string) string {
	if root != "#root" {
		return fmt.Sprintf("%s.%s", root, name)
	}
	return name
}
