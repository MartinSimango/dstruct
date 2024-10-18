package dstruct

type Builder interface {
	// AddField adds a field to the struct.
	AddField(name string, value interface{}, tag string) Builder

	// AddEmbeddedFields adds an embedded field to the struct.
	AddEmbeddedField(value interface{}, tag string) Builder

	// Build returns a DynamicStructModifier instance.
	Build() DynamicStructModifier

	// GetField returns a builder instance of the subfield of the struct that is currently being built. It will panic if the
	// field is not a struct or does not fully dereference to a struct value. It returns nil if field does not exist.
	GetField(name string) Builder

	// GetNewBuilderFromField returns a new builder instance where the subfield of the struct "field" is the root of the struct.
	NewBuilderFromField(field string) Builder

	// RemoveField removes a field from the struct. If the field is a subfield of a nil struct it will not be removed.
	RemoveField(field string) Builder
}
