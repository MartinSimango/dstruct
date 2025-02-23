package dstruct

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"unsafe"

	"github.com/MartinSimango/dstruct/dreflect"
)

type treeBuilderImpl struct {
	root      *Node[StructField]
	setValues bool
}

var _ Builder = &treeBuilderImpl{}

func createRoot() *Node[StructField] {
	return &Node[StructField]{
		data: &StructField{
			name:          "#root",
			value:         reflect.ValueOf(nil),
			qualifiedName: "#root",
		},
		children: make(map[string]*Node[StructField]),
	}
}

func NewBuilder() *treeBuilderImpl {
	return newBuilderFromNode(createRoot(), false)
}

func canExtend(val any) bool {
	if val == nil {
		return false
	}
	ptrValue, _ := getPtrValue(reflect.ValueOf(val), 0)

	if !ptrValue.IsValid() {
		return false
	}
	return ptrValue.Type().Kind() == reflect.Struct
}

func ExtendStruct(val any) *treeBuilderImpl {
	// TODO check if val is a struct
	b := NewBuilder()
	value := reflect.ValueOf(val)

	if !canExtend(val) {
		panic(fmt.Sprintf("Cannot extend struct value of type %s", value.Type()))
	}

	switch value.Kind() {
	case reflect.Struct:
		b.root.data.dstructType = b.addStructFields(value, b.root, 0)
	case reflect.Ptr:
		b.root.data.dstructType = b.addPtrField(value, b.root)
	}

	return b
}

func newBuilderFromNode(node *Node[StructField], resetFQN bool) *treeBuilderImpl {
	if resetFQN {
		resetNodeFieldsFQN(node)
	}
	return &treeBuilderImpl{
		setValues: true,
		root:      node,
	}
}

func resetNodeFieldsFQN(node *Node[StructField]) *Node[StructField] {
	for _, v := range node.children {
		v.data.qualifiedName = getQualifiedName(node.data.qualifiedName, v.data.name)
		resetNodeFieldsFQN(v)
	}
	return node
}

func (dsb *treeBuilderImpl) AddField(name string, value interface{}, tag string) Builder {
	dsb.addField(false, name, value, tag, dsb.root)
	return dsb
}

func (dsb *treeBuilderImpl) AddEmbeddedField(value interface{}, tag string) Builder {
	ptrValue, _ := getPtrValue(reflect.ValueOf(value), 0)
	name := reflect.TypeOf(ptrValue.Interface()).Name()

	dsb.addField(true, name, value, tag, dsb.root)

	return dsb
}

func (dsb *treeBuilderImpl) addField(
	anonymous bool,
	name string,
	value interface{},
	tag string,
	root *Node[StructField],
) {
	if root.HasChild(name) {
		panic(fmt.Sprintf("Field '%s' already exists", name))
	}

	dsb.root.data.dstructType = dsb.addFieldToTree(
		name,
		value,
		"",
		anonymous,
		reflect.StructTag(tag),
		dsb.root,
	)
}

func (dsb *treeBuilderImpl) RemoveField(name string) Builder {
	fields := strings.Split(name, ".")
	node := dsb.root

	for i := 0; i < len(fields)-1; i++ {
		node = node.GetNode(fields[i])
	}
	if node.children[fields[len(fields)-1]] != nil {
		node.DeleteNode(fields[len(fields)-1])
		*node.data.numberOfSubFields--
	}

	return dsb
}

// GetField implements Builder.GetField
func (dsb *treeBuilderImpl) GetField(field string) Builder {
	if node := dsb.getNode(field); node != nil {
		if node.data.ptrDepth > 0 {
			if node.data.ptrKind != reflect.Struct {
				panic(
					fmt.Sprintf(
						"Cannot get field '%s' because it a points to a non-struct value.",
						field,
					),
				)
			}
		} else if node.data.dstructType.Kind() != reflect.Struct {
			panic(fmt.Sprintf("Cannot get field '%s' because it is not a struct or does not fully derefence to a struct value", field))
		}

		return newBuilderFromNode(node, false)
	}
	panic(fmt.Sprintf("Field '%s' does not exist", field))
}

func (dsb *treeBuilderImpl) NewBuilderFromField(field string) Builder {
	copyNode := dsb.getNode(field).Copy()
	copyNode.data.name = "#root"
	copyNode.data.qualifiedName = "#root"
	return newBuilderFromNode(copyNode, true)
}

func (dsb *treeBuilderImpl) GetFieldCopy(field string) Builder {
	copyNode := dsb.getNode(field).Copy()
	return newBuilderFromNode(copyNode, false)
}

func (dsb *treeBuilderImpl) getNode(field string) *Node[StructField] {
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
	rootCopy := db.root.Copy()
	// Ensure that the current node is treated as root when struct is built
	if db.root.parent != nil {
		rootCopy.data.name = "#root"
		rootCopy.data.qualifiedName = "#root"
		rootCopy = resetNodeFieldsFQN(rootCopy)
	}
	return newStruct(db.buildStruct(rootCopy), rootCopy)
}

func (db *treeBuilderImpl) buildStruct(tree *Node[StructField]) any {
	structValue := reflect.ValueOf(dreflect.GetPointerOfValueType(treeToStruct(tree)))
	tree.data.value = structValue
	// set the value of the struct fields. Currently the tree structure contains the values of the fields
	// so we need to copy the values to the struct fields
	if db.setValues {
		setValues(structValue.Elem(), tree)
	}

	return structValue.Interface()
}

func setValues(value reflect.Value, tree *Node[StructField]) {
	if value.Kind() == reflect.Ptr {
		setPointerFieldValue(value, tree)
	} else {
		setStructFieldValues(value, tree)
	}
}

func (dsb *treeBuilderImpl) addFieldToTree(
	name string,
	typ interface{},
	pkgPath string,
	anonymous bool,
	tag reflect.StructTag,
	root *Node[StructField],
) reflect.Type {
	value := reflect.ValueOf(typ)
	if !value.IsValid() {
		panic(fmt.Sprintf("Cannot determine type of field '%s'", name))
	}
	if root.data.numberOfSubFields == nil {
		root.data.numberOfSubFields = new(int)
	} else {
		*root.data.numberOfSubFields++
	}

	field := &StructField{
		name:        name,
		value:       value, // this will initally be unaddressable until the struct is built
		tag:         tag,
		dstructType: value.Type(),
		typeHash:    dreflect.GetTypeHash(typ),
		pkgPath:     pkgPath,
		anonymous:   anonymous,
		jsonName:    strings.Split(tag.Get("json"), ",")[0],
	}

	field.structIndex = new(int)
	*field.structIndex = *root.data.numberOfSubFields
	field.qualifiedName = getQualifiedName(root.data.GetQualifiedName(), field.name)

	root.AddNode(name, field)

	// don't add struct fields if special kind
	switch value.Kind() {
	case reflect.Struct:
		field.dstructType = dsb.addStructFields(value, root.children[name], 0)

	case reflect.Ptr:
		field.dstructType = dsb.addPtrField(value, root.children[name])
	}

	// return field.dstructType
	// We need to rebuild the stryct to get it's new type as a result of adding a new field
	return reflect.TypeOf(dsb.Build().Instance())
}

func sortKeys(root *Node[StructField]) (keys []string) {
	for key := range root.children {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return *root.children[keys[i]].data.structIndex < *root.children[keys[j]].data.structIndex
	})

	return
}

// this only allocates memory for the struct and its fields and does not set any values
// so the returned value will be an uninitialized struct
func treeToStruct(root *Node[StructField]) any {
	var structFields []reflect.StructField

	// sort the keys to ensure type  of struct produced is always the same
	keys := sortKeys(root)

	for _, fieldName := range keys {
		var typ reflect.Type
		field := root.GetNode(fieldName)
		g, _ := getPtrValue(field.data.value, 0)
		if g.Kind() == reflect.Struct {
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

func setStructFieldValues(strct reflect.Value, root *Node[StructField]) {
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
		//  make the tree node value point to the struct field - ensuring that when the node value changes the struct field value changes
		// also update the node value to actually have a memory address
		currentNode.data.value = field

		// if currentNode.data.anonymous {
		// 	db.addAnonymousSubfields(currentNode)
		// }

	}
}

func setPointerFieldValue(field reflect.Value, currentNode *Node[StructField]) {
	if currentNode.data.value.IsNil() {
		return
	}

	f := field
	for i := 0; i < currentNode.data.ptrDepth; i++ {
		// We don't use f.Set(reflect.New(f.Type().Elem())) because it panics when the field is unexported
		// se we need to access the memory address of the field and set the value which bypasses the panic
		reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).
			Elem().
			Set(reflect.New(f.Type().Elem()))
		//			f.Set(reflect.New(f.Type().Elem()))
		f = f.Elem()
	}

	switch f.Kind() {
	case reflect.Struct:
		setStructFieldValues(f, currentNode)
	default:
		field.Set(currentNode.data.value)
	}
	currentNode.data.value = field
}

func (dsb *treeBuilderImpl) addStructFields(
	strct reflect.Value,
	root *Node[StructField],
	ptrDepth int,
) reflect.Type {
	var structFields []reflect.StructField

	// Create pointer to struct to allow unexported field values to be read in order
	// to obtain their types
	pointerToStruct := reflect.New(reflect.TypeOf(strct.Interface()))
	pointerToStruct.Elem().Set(strct)

	for i := 0; i < strct.NumField(); i++ {
		fieldName := strct.Type().Field(i).Name
		fieldTag := strct.Type().Field(i).Tag

		var fieldValue any
		if strct.Type().Field(i).IsExported() {
			fieldValue = strct.Field(i).Interface()
		} else {
			f := pointerToStruct.Elem().Field(i)
			fieldValue = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem().Interface()
		}

		pkgPath := strct.Type().Field(i).PkgPath
		anonymous := strct.Type().Field(i).Anonymous
		fieldType := dsb.addFieldToTree(fieldName, fieldValue, pkgPath, anonymous, fieldTag, root)

		if anonymous {
			pkgPath = ""
		}

		structFields = append(structFields, reflect.StructField{
			Name:      fieldName,
			PkgPath:   pkgPath,
			Type:      fieldType,
			Tag:       fieldTag,
			Anonymous: anonymous,
		})

	}

	// If the struct is a pointer, we need to create a pointer to the struct

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

func (dsb *treeBuilderImpl) addPtrField(
	value reflect.Value,
	node *Node[StructField],
) reflect.Type {
	if value.IsNil() {
		return reflect.TypeOf(value.Interface())
	}

	ptrValue, ptrDepth := getPtrValue(value, 0)

	node.data.ptrDepth = ptrDepth
	node.data.ptrKind = ptrValue.Kind()

	switch ptrValue.Kind() {
	case reflect.Struct:
		return dsb.addStructFields(ptrValue, node, ptrDepth)
	}
	return reflect.TypeOf(value.Interface())
}

func getQualifiedName(root, name string) string {
	if root != "#root" && root != "" {
		return fmt.Sprintf("%s.%s", root, name)
	}
	return name
}
