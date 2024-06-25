// Package config provides the configuration for the generation of dynamic structs.
package config

// Config represents the generation config for a dynamic struct.
type Config interface {

	// Number returns the number configuration.
	Number() NumberRangeConfig

	// Slice returns the slice configuration.
	Slice() SliceConfig

	// Date returns the date configuration.
	Date() DateRangeConfig

	// SetSliceLength sets the length range for the slice.
	SetSliceLength(min, max int) Config

	// SetIntRange sets the range for int values.
	SetIntRange(min, max int) Config

	// SetInt8Range sets the range for int8 values.
	SetInt8Range(min, max int8) Config

	// SetInt16Range sets the range for int16 values.
	SetInt16Range(min, max int16) Config

	// SetInt32Range sets the range for int32 values.
	SetInt32Range(min, max int32) Config

	// SetInt64Range sets the range for int64 values.
	SetInt64Range(min, max int64) Config

	// SetFloat32Range sets the range for float32 values.
	SetFloat32Range(min, max float32) Config

	// SetFloat64Range sets the range for float64 values.
	SetFloat64Range(min, max float64) Config

	// SetUIntRange sets the range for uint values.
	SetUIntRange(min, max uint) Config

	// SetUInt8Range sets the range for uint8 values.
	SetUInt8Range(min, max uint8) Config

	// SetUInt16Range sets the range for uint16 values.
	SetUInt16Range(min, max uint16) Config

	// SetUInt32Range sets the range for uint32 values.
	SetUInt32Range(min, max uint32) Config

	// SetUInt64Range sets the range for the uint64 value.
	SetUInt64Range(min, max uint64) Config

	// SetUIntPtrRange sets the range for the uintptr value.
	SetUIntPtr(min, max uintptr) Config

	// Copy returns a copy of the configuration.
	Copy() Config
}
