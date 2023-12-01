package config

import "github.com/MartinSimango/dstruct/util"

type NumberConfig interface {
	IntRange() (*int, *int)
	SetIntRange(min int, max int) NumberConfig
	Int8Range() (*int8, *int8)
	SetInt8Range(min int8, max int8) NumberConfig
	Int16Range() (*int16, *int16)
	SetInt16Range(min int16, max int16) NumberConfig
	Int32Range() (*int32, *int32)
	SetInt32Range(min int32, max int32) NumberConfig
	Int64Range() (*int64, *int64)
	SetInt64Range(min int64, max int64) NumberConfig
	Copy() NumberConfig
}
type Range[n util.Number] struct {
	min n
	max n
}

type NumberConfigImpl struct {
	intRange   Range[int]
	int8Range  Range[int8]
	int16Range Range[int16]
	int32Range Range[int32]
	int64Range Range[int64]
}

var _ NumberConfig = &NumberConfigImpl{}

// Int16Range implements NumberConfig.
func (nc *NumberConfigImpl) Int16Range() (*int16, *int16) {
	return &nc.int16Range.min, &nc.int16Range.max
}

// Int32Range implements NumberConfig.
func (nc *NumberConfigImpl) Int32Range() (*int32, *int32) {
	return &nc.int32Range.min, &nc.int32Range.max
}

// Int64Range implements NumberConfig.
func (nc *NumberConfigImpl) Int64Range() (*int64, *int64) {
	return &nc.int64Range.min, &nc.int64Range.max

}

// Int8Range implements NumberConfig.
func (nc *NumberConfigImpl) Int8Range() (*int8, *int8) {
	return &nc.int8Range.min, &nc.int8Range.max
}

// SetInt16Range implements NumberConfig.
func (nc *NumberConfigImpl) SetInt16Range(min int16, max int16) NumberConfig {
	setRange(min, max, &nc.int16Range.min, &nc.int16Range.max)
	return nc
}

// SetInt32Range implements NumberConfig.
func (nc *NumberConfigImpl) SetInt32Range(min int32, max int32) NumberConfig {
	setRange(min, max, &nc.int32Range.min, &nc.int32Range.max)
	return nc
}

// SetInt64Range implements NumberConfig.
func (nc *NumberConfigImpl) SetInt64Range(min int64, max int64) NumberConfig {
	setRange(min, max, &nc.int64Range.min, &nc.int64Range.max)
	return nc
}

// SetInt8Range implements NumberConfig.
func (nc *NumberConfigImpl) SetInt8Range(min int8, max int8) NumberConfig {
	setRange(min, max, &nc.int8Range.min, &nc.int8Range.max)
	return nc
}

func (nc *NumberConfigImpl) IntRange() (*int, *int) {
	return &nc.intRange.min, &nc.intRange.max

}

func (nc *NumberConfigImpl) SetIntRange(min, max int) NumberConfig {
	setRange(min, max, &nc.intRange.min, &nc.intRange.max)
	return nc
}

// Copy implements NumberConfig.
func (nc *NumberConfigImpl) Copy() NumberConfig {
	newNumberConfig := &NumberConfigImpl{}
	*newNumberConfig = *nc
	return newNumberConfig
}

func setRange[n util.Number](min, max n, ncMin, ncMax *n) {
	if min > max {
		return
	}
	*ncMin, *ncMax = min, max
}

func NewNumberConfig() *NumberConfigImpl {
	return &NumberConfigImpl{
		intRange:   Range[int]{0, 10},
		int8Range:  Range[int8]{0, 10},
		int16Range: Range[int16]{0, 10},
		int32Range: Range[int32]{0, 10},
		int64Range: Range[int64]{0, 10},
	}
}

// IntRange() (*n, *n)
// 	SetIntRange(min n, max n)

// type NumberConfig[n util.Number] interface {
// 	// TODO pointer values should only be accesible internally move files to same package
// 	MinRange() *n
// 	MaxRange() *n
// 	SetRange(min n, max n)
// 	Copy() NumberConfig[n]
// }

// type NumberConfigImpl[n util.Number] struct {
// 	minRange n
// 	maxRange n
// }

// var _ NumberConfig[int] = &NumberConfigImpl[int]{}

// func NewNumberConfig[n util.Number]() (cfg *NumberConfigImpl[n]) {
// 	return &NumberConfigImpl[n]{
// 		minRange: 0,
// 		maxRange: 10,
// 	}
// }

// func NewNumberConfigWithRange[n util.Number](min, max n) (cfg *NumberConfigImpl[n]) {
// 	return &NumberConfigImpl[n]{
// 		minRange: min,
// 		maxRange: max,
// 	}
// }

// func (nc *NumberConfigImpl[n]) MinRange() *n {
// 	return &nc.minRange

// }

// func (nc *NumberConfigImpl[n]) MaxRange() *n {
// 	return &nc.maxRange
// }

// func (nc *NumberConfigImpl[n]) Copy() NumberConfig[n] {
// 	newNumberConfig := &NumberConfigImpl[n]{}
// 	*newNumberConfig = *nc
// 	return newNumberConfig
// }
