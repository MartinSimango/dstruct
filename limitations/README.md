Limitations:
* You cannot extend structs with unexported embedded fields.
* If a struct pointer cannot be fully dereferenced then the struct's subfields won't be added to the dynamic struct.
This is done mainly to avoid self referencing structs as these will create infinite node trees.
For example
```go

type NilStructPointer struct {
	Field  int
}

type Struct struct {
	Field string
	NilField **NilStructPointer
}

func main() {
	// This will create a struct with fields  
	// - Field
	// - NilField 
	// because we cannot fully dereference a.NilField (**a.NilField will cause the program to panic)
	var n *NilStructPointer
	a := Struct{NilField: &n}
	dstruct.ExtendStruct(a).Build()

	// However this will create a struct with fields 
	// - Field
	// - NilField 
	// - NilField.Field
	n = &NilStructPointer{}
	a = Struct{NilField: &n}
	dstruct.ExtendStruct(a).Build()

}

```
* Dynamic structs with struct fields of type `any (interface {})` cannot be created. If you try
extend or merge structs which have struct fields of type `any` their value must be set to a concrete type. 

For example

```go
type A struct {
	Field any
}

func main() {
	// This will be fine
	dstruct.ExtendStruct(A{Field: 2})

	// This will panic as Field's type cannot be determined
	dstruct.ExtendStruct(A{})

	// This will also panic as Field's type cannot be determined
	dstruct.NewBuilder().AddField("Field", nil, `json:"Field"`)

}

* Generated structs cannot generate slice or array fields whose types have recursive definition.

```