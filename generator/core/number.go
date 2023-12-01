package core

import (
	"math/rand"

	"github.com/MartinSimango/dstruct/generator"
	"github.com/MartinSimango/dstruct/generator/config"
	"github.com/MartinSimango/dstruct/util"
)

func generateNum[n util.Number](min, max n) n {
	return min + (n(rand.Float64() * float64(max+1-min)))
}

func GenerateNumberFunc[n util.Number](cfg config.NumberConfig) generator.GenerationFunction {
	min, max := getNumberRange[n](cfg)
	return &coreGenerationFunction{
		_func: func(parameters ...any) any {
			return generateNum(*min, *max)
		},
		args: []any{min, max},
	}

}

func getNumberRange[n util.Number](cfg config.NumberConfig) (*n, *n) {
	var min, max any
	switch any(*new(n)).(type) {
	case int:
		min, max = cfg.IntRange()

	case int8:
		min, max = cfg.Int8Range()
	case int32:
		min, max = cfg.Int32Range()
	case int64:
		min, max = cfg.Int64Range()
	}
	return any(min).(*n), any(max).(*n)
}

func NewGenerateNumberFunctionHolder[N util.Number](numberConfig config.NumberConfig) *NumberFunctionHolder {
	return NewNumberFunctionHolder(GenerateNumberFunc[N], numberConfig)
}
