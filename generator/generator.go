package generator

import (
	"reflect"
)

type GenerationFunction interface {
	Generate() any
	Copy(*GenerationConfig) GenerationFunction
	GetGenerationConfig() *GenerationConfig
	SetGenerationConfig(*GenerationConfig) GenerationFunction
}

type GenerationFunctionImpl struct {
	*GenerationConfig
	basicGenerationFunction
}

func (bsf basicGenerationFunction) Generate() any {
	return bsf._func(bsf.args...)
}

func (bsf basicGenerationFunction) Copy(*GenerationConfig) GenerationFunction {
	return bsf
}

func (bsf basicGenerationFunction) GetGenerationConfig() *GenerationConfig {
	return nil
}

func (bsf basicGenerationFunction) SetGenerationConfig(*GenerationConfig) GenerationFunction {
	return bsf
}

type DefaultGenerationFunctionType map[reflect.Kind]GenerationFunction

type Generator struct {
	GenerationConfig           *GenerationConfig
	DefaultGenerationFunctions DefaultGenerationFunctionType
}

func NewGenerator(gc *GenerationConfig) *Generator {

	defaultGenerationFunctions := make(DefaultGenerationFunctionType)
	defaultGenerationFunctions[reflect.String] = GenerateFixedValueFunc("string")
	defaultGenerationFunctions[reflect.Ptr] = GenerateNilValueFunc()
	defaultGenerationFunctions[reflect.Int] = GenerateNumberFunc(gc.intMin, gc.intMax).SetGenerationConfig(gc)
	defaultGenerationFunctions[reflect.Int64] = GenerateNumberFunc(gc.int64Min, gc.int64Max).SetGenerationConfig(gc)
	defaultGenerationFunctions[reflect.Int32] = GenerateNumberFunc(gc.int32Min, gc.int32Max).SetGenerationConfig(gc)
	defaultGenerationFunctions[reflect.Float64] = GenerateNumberFunc(gc.float64Min, gc.float64Max).SetGenerationConfig(gc)
	defaultGenerationFunctions[reflect.Bool] = GenerateBoolFunc()

	return &Generator{
		GenerationConfig:           gc,
		DefaultGenerationFunctions: defaultGenerationFunctions,
	}
}

func (gd *Generator) Copy() (generationDefaults *Generator) {
	generationDefaults = &Generator{
		DefaultGenerationFunctions: make(DefaultGenerationFunctionType),
		GenerationConfig:           gd.GenerationConfig.Clone(),
	}
	for k, v := range gd.DefaultGenerationFunctions {
		generationDefaults.DefaultGenerationFunctions[k] = v.Copy(generationDefaults.GenerationConfig)
	}
	return
}
