package generator

import (
	"reflect"
)

type GenerationFunction interface {
	Generate() any
}

type GenerationFunc struct {
	_func func(...any) any
	args  []any
}

var _ GenerationFunction = &GenerationFunc{}

func (f GenerationFunc) Generate() any {
	return f._func(f.args...)
}

type GenerationDefaults map[reflect.Kind]GenerationFunction

func NewGenerationConfig() (generationConfig *GenerationConfig) {
	generationConfig = &GenerationConfig{
		GenerationValueConfig: GenerationValueConfig{
			valueGenerationType:  UseDefaults,
			setNonRequiredFields: false,
		},
		DefaultGenerationFunctions: make(GenerationDefaults),
		SliceConfig:                defaultSliceConfig(),
		IntConfig:                  defaultIntConfig(),
		FloatConfig:                defaultFloatConfig(),
		DateConfig:                 defaultDateConfig(),
	}
	generationConfig.initGenerationFunctionDefaults()
	return
}

func (gc *GenerationConfig) initGenerationFunctionDefaults() {
	gc.DefaultGenerationFunctions[reflect.String] = GenerateFixedValueFunc("string")
	gc.DefaultGenerationFunctions[reflect.Ptr] = GenerateNilValueFunc()
	gc.DefaultGenerationFunctions[reflect.Int64] = GenerateNumberFunc(&gc.int64Min, &gc.int64Max)
	gc.DefaultGenerationFunctions[reflect.Int32] = GenerateNumberFunc(&gc.int32Min, &gc.int32Max)
	gc.DefaultGenerationFunctions[reflect.Int] = GenerateNumberFunc(&gc.intMin, &gc.intMax)
	gc.DefaultGenerationFunctions[reflect.Float64] = GenerateNumberFunc(&gc.float64Min, &gc.float64Max)
	gc.DefaultGenerationFunctions[reflect.Bool] = GenerateBoolFunc()
}
