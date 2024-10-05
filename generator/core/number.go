package core

import (
	"fmt"
	"math/rand"
	"reflect"

	"github.com/MartinSimango/dstruct/generator"
	"github.com/MartinSimango/dstruct/generator/config"
)

func generateNum[n config.Number](min, max n) n {
	return min + (n(rand.Float64() * float64(max+1-min)))
}

func GenerateNumberFunc[n config.Number](
	cfg config.NumberRangeConfig,
) generator.GenerationFunction {
	// get reference to min and max so when the function is called, it will use the current value
	min, max := getNumberRange[n](cfg)
	return &coreGenerationFunction{
		_func: func(parameters ...any) any {
			return generateNum(*min, *max)
		},
		kind: reflect.ValueOf(*new(n)).Kind(),
	}
}

func GenerateSequential(seed int) generator.GenerationFunction {
	return &coreGenerationFunction{
		_func: func(parameters ...any) any {
			s := parameters[0].(*int)
			*s++
			return *s
		},
		kind: reflect.Int,
		args: []any{&seed},
	}
}

func getNumberRange[n config.Number](cfg config.NumberRangeConfig) (*n, *n) {
	var min, max any
	v := any(*new(n))
	switch v.(type) {
	case int:
		min, max = cfg.Int().RangeRef()
	case int8:
		min, max = cfg.Int8().RangeRef()
	case int16:
		min, max = cfg.Int16().RangeRef()
	case int32:
		min, max = cfg.Int32().RangeRef()
	case int64:
		min, max = cfg.Int64().RangeRef()
	case uint:
		min, max = cfg.UInt().RangeRef()
	case uint8:
		min, max = cfg.UInt8().RangeRef()
	case uint16:
		min, max = cfg.UInt16().RangeRef()
	case uint32:
		min, max = cfg.UInt32().RangeRef()
	case uint64:
		min, max = cfg.UInt64().RangeRef()
	case float32:
		min, max = cfg.Float32().RangeRef()
	case float64:
		min, max = cfg.Float64().RangeRef()
	case uintptr:
		min, max = cfg.UIntPtr().RangeRef()
	default:
		panic(fmt.Sprintf("Type not supported for getNumberRange: %s", reflect.TypeOf(v)))
	}
	return any(min).(*n), any(max).(*n)
}

func NewGenerateNumberFunctionHolder[N config.Number](
	numberConfig config.NumberRangeConfig,
) *NumberFunctionHolder {
	return NewNumberFunctionHolder(GenerateNumberFunc[N], numberConfig)
}
