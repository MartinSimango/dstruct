package config

import (
	"fmt"
)

// DstructNumberRangeConfig implements NumberRangeConfig.
type DstructNumberRangeConfig struct {
	IntConfig     NumberRange[int]
	Int8Config    NumberRange[int8]
	Int16Config   NumberRange[int16]
	Int32Config   NumberRange[int32]
	Int64Config   NumberRange[int64]
	Float32Config NumberRange[float32]
	Float64Config NumberRange[float64]
	UIntConfig    NumberRange[uint]
	UInt8Config   NumberRange[uint8]
	UInt16Config  NumberRange[uint16]
	UInt32Config  NumberRange[uint32]
	UInt64Config  NumberRange[uint64]
	UIntPtrConfig NumberRange[uintptr]
}

// NewNumberRangeConfig is a constructor for DstructNumberRangeConfig.
func NewNumberRangeConfig() *DstructNumberRangeConfig {
	return &DstructNumberRangeConfig{
		IntConfig:     numberRange[int](0, 10),
		Int8Config:    numberRange[int8](0, 10),
		Int16Config:   numberRange[int16](0, 10),
		Int32Config:   numberRange[int32](0, 10),
		Int64Config:   numberRange[int64](0, 10),
		Float32Config: numberRange[float32](0, 10),
		Float64Config: numberRange[float64](0, 10),
		UIntConfig:    numberRange[uint](0, 10),
		UInt8Config:   numberRange[uint8](0, 10),
		UInt16Config:  numberRange[uint16](0, 10),
		UInt32Config:  numberRange[uint32](0, 10),
		UInt64Config:  numberRange[uint64](0, 10),
		UIntPtrConfig: numberRange[uintptr](0, 10),
	}
}

var _ NumberRangeConfig = &DstructNumberRangeConfig{}

// Float32 implements NumberConfig.Float32.
func (nc *DstructNumberRangeConfig) Float32() *NumberRange[float32] {
	return &nc.Float32Config
}

// Float64 implements NumberConfig.Float64.
func (nc *DstructNumberRangeConfig) Float64() *NumberRange[float64] {
	return &nc.Float64Config
}

// Int implements NumberConfig.Int.
func (nc *DstructNumberRangeConfig) Int() *NumberRange[int] {
	return &nc.IntConfig
}

// Int16 implements NumberConfig.Int16.
func (nc *DstructNumberRangeConfig) Int16() *NumberRange[int16] {
	return &nc.Int16Config
}

// Int32 implements NumberConfig.Int32.
func (nc *DstructNumberRangeConfig) Int32() *NumberRange[int32] {
	return &nc.Int32Config
}

// Int64 implements NumberConfig.Int64.
func (nc *DstructNumberRangeConfig) Int64() *NumberRange[int64] {
	return &nc.Int64Config
}

// Int8 implements NumberConfig.Int8.
func (nc *DstructNumberRangeConfig) Int8() *NumberRange[int8] {
	return &nc.Int8Config
}

// UInt implements NumberConfig.UInt.
func (nc *DstructNumberRangeConfig) UInt() *NumberRange[uint] {
	return &nc.UIntConfig
}

// UInt16 implements NumberConfig.UInt16.
func (nc *DstructNumberRangeConfig) UInt16() *NumberRange[uint16] {
	return &nc.UInt16Config
}

// UInt32 implements NumberConfig.UInt32,
func (nc *DstructNumberRangeConfig) UInt32() *NumberRange[uint32] {
	return &nc.UInt32Config
}

// UInt64 implements NumberConfig.UInt64.
func (nc *DstructNumberRangeConfig) UInt64() *NumberRange[uint64] {
	return &nc.UInt64Config
}

// UInt8 implements NumberConfig.UInt8.
func (nc *DstructNumberRangeConfig) UInt8() *NumberRange[uint8] {
	return &nc.UInt8Config
}

// UIntptr implements NumberConfig.UIntPtr.
func (nc *DstructNumberRangeConfig) UIntPtr() *NumberRange[uintptr] {
	return &nc.UIntPtrConfig
}

// Copy implements NumberConfig.Copy.
func (nc *DstructNumberRangeConfig) Copy() NumberRangeConfig {
	fmt.Println("Copying DstructNumberRangeConfig ")
	fmt.Println(nc.Int().RangeRef())
	numberConfigCopy := new(DstructNumberRangeConfig)
	*numberConfigCopy = *nc
	return numberConfigCopy
}
