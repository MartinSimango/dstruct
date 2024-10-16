package config

// Number is a type that represents a number.
type Number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64 |
		uint | uint8 | uint16 | uint32 | uint64 | uintptr
}

// numRange represents the range of a number.
type NumberRange[n Number] struct {
	min n
	max n
}

func (nr *NumberRange[n]) Range() (n, n) {
	return nr.min, nr.max
}

func (nr *NumberRange[n]) Max() n {
	return nr.max
}

func (nr *NumberRange[n]) Min() n {
	return nr.min
}

func (nr *NumberRange[n]) SetRange(min n, max n) {
	if min < max {
		nr.min = min
		nr.max = max
	}
}

func (nr *NumberRange[n]) SetMax(max n) {
	nr.SetRange(nr.min, max)
}

func (nr *NumberRange[n]) SetMin(min n) {
	nr.SetRange(min, nr.max)
}

func (nr *NumberRange[n]) RangeRef() (*n, *n) {
	return &nr.min, &nr.max
}

func numberRange[N Number](min, max N) NumberRange[N] {
	return NumberRange[N]{min, max}
}

// NumberRangeConfig represents the configuration for a number range within a dynamic struct.
type NumberRangeConfig interface {
	// Int returns the configuration for an int number range.
	Int() *NumberRange[int]

	// Int8 returns the configuration for an int8 number range.
	Int8() *NumberRange[int8]

	// Int16 returns the configuration for an int16 number range.
	Int16() *NumberRange[int16]

	// Int32 returns the configuration for an int32 number range.
	Int32() *NumberRange[int32]

	// Int64 returns the configuration for an int64 number range.
	Int64() *NumberRange[int64]

	// Float32 returns the configuration for a float32 number range.
	Float32() *NumberRange[float32]

	// Float64 returns the configuration for a float64 number range.
	Float64() *NumberRange[float64]

	// UInt returns the configuration for a uint number range.
	UInt() *NumberRange[uint]

	// UInt8 returns the configuration for a uint8 number range.
	UInt8() *NumberRange[uint8]

	// UInt16 returns the configuration for a uint16 number range.
	UInt16() *NumberRange[uint16]

	// Uint32 returns the configuration for a uint32 number range.
	UInt32() *NumberRange[uint32]

	// UInt64 returns the configuration for a uint64 number range.
	UInt64() *NumberRange[uint64]

	// UIntPtr returns the configuration for a uintptr number range.
	UIntPtr() *NumberRange[uintptr]

	// Copy returns a copy of the configuration.
	Copy() NumberRangeConfig

	// SetFrom sets the configuration from another configuration.
	SetFrom(cfg NumberRangeConfig)
}
