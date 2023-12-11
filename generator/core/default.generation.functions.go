package core

import (
	"reflect"

	"github.com/MartinSimango/dstruct/generator/config"
)

type DefaultGenerationFunctions map[reflect.Kind]FunctionHolder

func (d DefaultGenerationFunctions) Copy(kind reflect.Kind) (dgf DefaultGenerationFunctions) {
	// look at the kind and only return what needs to be copied
	dgf = make(DefaultGenerationFunctions)
	switch kind {
	case reflect.Struct, reflect.Slice, reflect.Ptr:
		for k, v := range d {
			dgf[k] = v.Copy()
		}
	default:
		dgf[kind] = d[kind].Copy()
	}
	return dgf

}

func NewDefaultGenerationFunctions(cfg config.Config) DefaultGenerationFunctions {

	defaultGenerationFunctions := make(DefaultGenerationFunctions)
	defaultGenerationFunctions[reflect.String] = NewFixedFunctionHolder(GenerateStringFromRegexFunc, "^[a-zA-Z]{3}$")
	defaultGenerationFunctions[reflect.Int] = NewGenerateNumberFunctionHolder[int](cfg.Number())
	// NewFixedFunctionHolder(GenerateSequential, 0)
	// NewGenerateNumberFunctionHolder[int](cfg.Number())
	defaultGenerationFunctions[reflect.Int8] = NewGenerateNumberFunctionHolder[int8](cfg.Number())
	defaultGenerationFunctions[reflect.Int16] = NewGenerateNumberFunctionHolder[int16](cfg.Number())
	defaultGenerationFunctions[reflect.Int32] = NewGenerateNumberFunctionHolder[int32](cfg.Number())
	defaultGenerationFunctions[reflect.Int64] = NewGenerateNumberFunctionHolder[int64](cfg.Number())

	defaultGenerationFunctions[reflect.Uint] = NewGenerateNumberFunctionHolder[uint](cfg.Number())
	defaultGenerationFunctions[reflect.Uint8] = NewGenerateNumberFunctionHolder[uint8](cfg.Number())
	defaultGenerationFunctions[reflect.Uint16] = NewGenerateNumberFunctionHolder[uint16](cfg.Number())
	defaultGenerationFunctions[reflect.Uint32] = NewGenerateNumberFunctionHolder[uint32](cfg.Number())
	defaultGenerationFunctions[reflect.Uint64] = NewGenerateNumberFunctionHolder[uint64](cfg.Number())

	defaultGenerationFunctions[reflect.Float32] = NewGenerateNumberFunctionHolder[float32](cfg.Number())
	defaultGenerationFunctions[reflect.Float64] = NewGenerateNumberFunctionHolder[float64](cfg.Number())
	defaultGenerationFunctions[reflect.Bool] = NewFunctionHolderNoArgs(GenerateBoolFunc())
	defaultGenerationFunctions[reflect.Ptr] = NewFunctionHolderNoArgs(GenerateNilValueFunc())
	return defaultGenerationFunctions

}
