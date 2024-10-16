package core

import (
	"fmt"
	"reflect"

	"github.com/MartinSimango/dstruct/generator/config"
)

type DefaultGenerationFunctions map[reflect.Kind]FunctionHolder

func (d DefaultGenerationFunctions) Copy(
	cfg config.Config,
	kind reflect.Kind,
) (dgf DefaultGenerationFunctions) {
	// look at the kind and only return what needs to be copied
	dgf = make(DefaultGenerationFunctions)
	if d[kind] == nil {
		// Copy ll the functions if the kind is not found
		// TODO: check if this can ever happen
		fmt.Println("Kind not found", kind)
		for k, v := range d {
			dgf[k] = v.Copy(cfg)
		}
		return
	}
	dgf[kind] = d[kind].Copy(cfg)
	return
}

func NewDefaultGenerationFunctions(cfg config.Config) DefaultGenerationFunctions {
	defaultGenerationFunctions := make(DefaultGenerationFunctions)
	defaultGenerationFunctions[reflect.String] = NewFixedFunctionHolder(
		GenerateStringFromRegexFunc,
		"^[a-zA-Z]{3}$",
	)
	defaultGenerationFunctions[reflect.Int] = NewGenerateNumberFunctionHolder[int](cfg.Number())
	// NewFixedFunctionHolder(GenerateSequential, 0)
	// NewGenerateNumberFunctionHolder[int](cfg.Number())
	defaultGenerationFunctions[reflect.Int8] = NewGenerateNumberFunctionHolder[int8](cfg.Number())
	defaultGenerationFunctions[reflect.Int16] = NewGenerateNumberFunctionHolder[int16](cfg.Number())
	defaultGenerationFunctions[reflect.Int32] = NewGenerateNumberFunctionHolder[int32](cfg.Number())
	defaultGenerationFunctions[reflect.Int64] = NewGenerateNumberFunctionHolder[int64](cfg.Number())

	defaultGenerationFunctions[reflect.Uint] = NewGenerateNumberFunctionHolder[uint](cfg.Number())
	defaultGenerationFunctions[reflect.Uint8] = NewGenerateNumberFunctionHolder[uint8](cfg.Number())
	defaultGenerationFunctions[reflect.Uint16] = NewGenerateNumberFunctionHolder[uint16](
		cfg.Number(),
	)
	defaultGenerationFunctions[reflect.Uint32] = NewGenerateNumberFunctionHolder[uint32](
		cfg.Number(),
	)
	defaultGenerationFunctions[reflect.Uint64] = NewGenerateNumberFunctionHolder[uint64](
		cfg.Number(),
	)

	defaultGenerationFunctions[reflect.Float32] = NewGenerateNumberFunctionHolder[float32](
		cfg.Number(),
	)
	defaultGenerationFunctions[reflect.Float64] = NewGenerateNumberFunctionHolder[float64](
		cfg.Number(),
	)
	defaultGenerationFunctions[reflect.Bool] = NewFunctionHolderNoArgs(GenerateBoolFunc())
	defaultGenerationFunctions[reflect.Ptr] = NewFunctionHolderNoArgs(GenerateNilValueFunc())
	return defaultGenerationFunctions
}
