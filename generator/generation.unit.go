package generator

import (
	"reflect"
)

type GenerationUnit struct {
	GenerationConfig    *GenerationConfig
	PreviousValueConfig GenerationValueConfig
	CurrentFunction     GenerationFunction
	count               int
	latestValue         any
	field               Field
}

func NewGenerationUnit(generationConfig *GenerationConfig,
	field Field,
) *GenerationUnit {
	gu := &GenerationUnit{
		GenerationConfig: generationConfig,
		field:            field,
	}
	gu.PreviousValueConfig = gu.GenerationConfig.GenerationValueConfig
	gu.CurrentFunction = gu.getGenerationFunction() //GetGenerationFunction(gu.field, *gu.GenerationConfig)

	return gu
}

func (gu *GenerationUnit) Generate() any {
	// check if important fields have changed and then regenerate the currentfunction
	if gu.configChanged(gu.PreviousValueConfig) {
		gu.CurrentFunction = gu.getGenerationFunction()
	}
	if gu.GenerationConfig.valueGenerationType == GenerateOnce && gu.count > 0 {
		return gu.latestValue
	}

	gu.latestValue = gu.CurrentFunction.Generate()
	gu.PreviousValueConfig = gu.GenerationConfig.GenerationValueConfig
	gu.count++
	return gu.latestValue
}

func (gu *GenerationUnit) configChanged(previousConfig GenerationValueConfig) bool {

	return gu.GenerationConfig.valueGenerationType != previousConfig.valueGenerationType ||
		gu.GenerationConfig.setNonRequiredFields != previousConfig.setNonRequiredFields
}

func (gu *GenerationUnit) SetGenerationDefaultFunctionForKind(kind reflect.Kind, function GenerationFunction) {
	if gu.GenerationConfig.DefaultGenerationFunctions[kind] == nil {
		return
	}
	gu.GenerationConfig.DefaultGenerationFunctions[kind] = function
	gu.CurrentFunction = gu.getGenerationFunction()

}

func (gu *GenerationUnit) getGenerationFunction() GenerationFunction {

	switch gu.field.Value.Kind() {
	case reflect.Slice:
		return GenerateSliceFunc(gu.GenerationConfig, gu.field.Value)
	case reflect.Struct:
		return GenerateStructFunc(gu.GenerationConfig, gu.field.Value)
	case reflect.Ptr:
		if gu.GenerationConfig.setNonRequiredFields {
			return GeneratePointerValueFunc(gu.GenerationConfig, gu.field.Value, gu.field.Tag)
		} else {
			return GenerateNilValueFunc()
		}
	}
	return generationFunctionFromTags(gu.field.Value.Kind(), gu.field.Tag, gu.GenerationConfig)

}
