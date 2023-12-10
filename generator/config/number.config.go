package config

type Number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64 |
		uint | uint8 | uint16 | uint32 | uint64 | uintptr
}

type NumRange[n Number] struct {
	min n
	max n
}
type NumberRangeImpl[n Number] struct {
	NumRange[n]
}

func (nr *NumberRangeImpl[n]) Range() (n, n) {
	return nr.min, nr.max
}

func (nr *NumberRangeImpl[n]) Max() n {
	return nr.max

}
func (nr *NumberRangeImpl[n]) Min() n {
	return nr.min

}

func (nr *NumberRangeImpl[n]) SetRange(min n, max n) {
	if min < max {
		nr.min = min
		nr.max = max
	}

}

func (nr *NumberRangeImpl[n]) SetMax(max n) {
	nr.SetRange(nr.min, max)

}
func (nr *NumberRangeImpl[n]) SetMin(min n) {
	nr.SetRange(min, nr.max)

}

func (nr *NumberRangeImpl[n]) RangeRef() (*n, *n) {
	return &nr.min, &nr.max

}

type NumberConfig interface {
	Int() *NumberRangeImpl[int]
	Int8() *NumberRangeImpl[int8]
	Int16() *NumberRangeImpl[int16]
	Int32() *NumberRangeImpl[int32]
	Int64() *NumberRangeImpl[int64]
	Float32() *NumberRangeImpl[float32]
	Float64() *NumberRangeImpl[float64]
	UInt() *NumberRangeImpl[uint]
	UInt8() *NumberRangeImpl[uint8]
	UInt16() *NumberRangeImpl[uint16]
	UInt32() *NumberRangeImpl[uint32]
	UInt64() *NumberRangeImpl[uint64]
	UIntPtr() *NumberRangeImpl[uintptr]
	Copy() NumberConfig
}

type IntConfig NumberRangeImpl[int]
type NumberConfigImpl struct {
	IntConfig     NumberRangeImpl[int]
	Int8Config    NumberRangeImpl[int8]
	Int16Config   NumberRangeImpl[int16]
	Int32Config   NumberRangeImpl[int32]
	Int64Config   NumberRangeImpl[int64]
	Float32Config NumberRangeImpl[float32]
	Float64Config NumberRangeImpl[float64]
	UIntConfig    NumberRangeImpl[uint]
	UInt8Config   NumberRangeImpl[uint8]
	UInt16Config  NumberRangeImpl[uint16]
	UInt32Config  NumberRangeImpl[uint32]
	UInt64Config  NumberRangeImpl[uint64]
	UIntPtrConfig NumberRangeImpl[uintptr]
}

func NumberRangeConfig[N Number](min, max N) NumberRangeImpl[N] {
	return NumberRangeImpl[N]{NumRange[N]{min, max}}
}

func NewNumberConfig() *NumberConfigImpl {
	return &NumberConfigImpl{
		IntConfig:     NumberRangeConfig[int](0, 10),
		Int8Config:    NumberRangeConfig[int8](0, 10),
		Int16Config:   NumberRangeConfig[int16](0, 10),
		Int32Config:   NumberRangeConfig[int32](0, 10),
		Int64Config:   NumberRangeConfig[int64](0, 10),
		Float32Config: NumberRangeConfig[float32](0, 10),
		Float64Config: NumberRangeConfig[float64](0, 10),
		UIntConfig:    NumberRangeConfig[uint](0, 10),
		UInt8Config:   NumberRangeConfig[uint8](0, 10),
		UInt16Config:  NumberRangeConfig[uint16](0, 10),
		UInt32Config:  NumberRangeConfig[uint32](0, 10),
		UInt64Config:  NumberRangeConfig[uint64](0, 10),
		UIntPtrConfig: NumberRangeConfig[uintptr](0, 10),
	}

}

var _ NumberConfig = &NumberConfigImpl{}

// Float32 implements NumberConfig.
func (nc *NumberConfigImpl) Float32() *NumberRangeImpl[float32] {
	return &nc.Float32Config
}

// Float64 implements NumberConfig.
func (nc *NumberConfigImpl) Float64() *NumberRangeImpl[float64] {
	return &nc.Float64Config
}

// Int implements NumberConfig.
func (nc *NumberConfigImpl) Int() *NumberRangeImpl[int] {
	return &nc.IntConfig
}

// Int16 implements NumberConfig.
func (nc *NumberConfigImpl) Int16() *NumberRangeImpl[int16] {
	return &nc.Int16Config
}

// Int32 implements NumberConfig.
func (nc *NumberConfigImpl) Int32() *NumberRangeImpl[int32] {
	return &nc.Int32Config
}

// Int64 implements NumberConfig.
func (nc *NumberConfigImpl) Int64() *NumberRangeImpl[int64] {
	return &nc.Int64Config
}

// Int8 implements NumberConfig.
func (nc *NumberConfigImpl) Int8() *NumberRangeImpl[int8] {
	return &nc.Int8Config

}

// UInt implements NumberConfig.
func (nc *NumberConfigImpl) UInt() *NumberRangeImpl[uint] {
	return &nc.UIntConfig
}

// UInt16 implements NumberConfig.
func (nc *NumberConfigImpl) UInt16() *NumberRangeImpl[uint16] {
	return &nc.UInt16Config
}

// UInt32 implements NumberConfig.
func (nc *NumberConfigImpl) UInt32() *NumberRangeImpl[uint32] {
	return &nc.UInt32Config

}

// UInt64 implements NumberConfig.
func (nc *NumberConfigImpl) UInt64() *NumberRangeImpl[uint64] {
	return &nc.UInt64Config
}

// UInt8 implements NumberConfig.
func (nc *NumberConfigImpl) UInt8() *NumberRangeImpl[uint8] {
	return &nc.UInt8Config
}

// UIntptr implements NumberConfig.
func (nc *NumberConfigImpl) UIntPtr() *NumberRangeImpl[uintptr] {
	return &nc.UIntPtrConfig
}

// UIntptr implements NumberConfig.
func (nc *NumberConfigImpl) Copy() NumberConfig {
	return simplePointerCopy(nc)
	// return &NumberConfigImpl{
	// 	IntConfig:     nc.IntConfig,
	// 	Int8Config:    nc.Int8Config,
	// 	Int16Config:   nc.Int16Config,
	// 	Int32Config:   nc.Int32Config,
	// 	Int64Config:   nc.Int64Config,
	// 	Float32Config: nc.Float32Config,
	// 	Float64Config: nc.Float64Config,
	// 	UIntConfig:    nc.UIntConfig,
	// 	UInt8Config:   nc.UInt8Config,
	// 	UInt16Config:  nc.UInt16Config,
	// 	UInt32Config:  nc.UInt32Config,
	// 	UInt64Config:  nc.UInt64Config,
	// 	UIntPtrConfig: nc.UIntPtrConfig,
	// }
}

func simplePointerCopy[T any](p *T) *T {
	v := new(T)
	*v = *p
	return v

}
