package dstruct

type Builder interface {
	// AddField adds a field to the struct.
	AddField(name string, value interface{}, tag string) Builder

	// AddEmbeddedFields adds an embedded field to the struct.
	AddEmbeddedField(value interface{}, tag string) Builder

	// Build returns a DynamicStructModifier instance.
	Build() DynamicStructModifier

	// GetField returns a builder instance of the subfield of the struct that is currently being built.
	GetField(name string) Builder

	// GetFieldCopy returns a copy of a builder instance of the subfield of the struct that is currently being built.
	//
	// Deprecated: this method will be removed use NewBuilderFromField instead.
	GetFieldCopy(name string) Builder

	// GetNewBuilderFromField returns a new builder instance where the subfield of the struct "field" is the root of the struct.
	NewBuilderFromField(field string) Builder

	// RemoveField removes a field from the struct.
	RemoveField(name string) Builder
}
