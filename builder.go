package dstruct

type Builder interface {
	AddField(name string, typ interface{}, tag string) Builder
	Build() DynamicStructModifier
	BuildWithFieldModifier(fieldModifier FieldModifier) DynamicStructModifier
	GetField(name string) Builder
	GetFieldCopy(field string) Builder
	RemoveField(name string) Builder
}
