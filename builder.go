package dstruct

type Builder interface {
	// AddField adds a field to the struct.
	AddField(name string, value interface{}, tag string) Builder

	// Build returns a DynamicStructModifier instance.
	Build() DynamicStructModifier

	// GetField returns a builder instance of the subfield of the struct that is currently being built.
	GetField(name string) Builder

	// GetFieldCopy returns a copy of a builder instance of the subfield of the struct that is currently being built.
	GetFieldCopy(field string) Builder

	// RemoveField removes a field from the struct.
	RemoveField(name string) Builder
}
