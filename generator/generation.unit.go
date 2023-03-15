package generator

import (
	"reflect"
)

type GenerationUnit struct {
	PreviousValueConfig   GenerationValueConfig
	CurrentFunction       GenerationFunction
	UpdateCurrentFunction bool
	Field                 *GeneratedField
	count                 int
	latestValue           any
	generationConfig      *GenerationConfig
}

func NewGenerationUnit(field *GeneratedField) *GenerationUnit {
	gu := &GenerationUnit{
		Field:            field,
		generationConfig: field.Generator.GenerationConfig,
	}
	gu.PreviousValueConfig = gu.generationConfig.GenerationValueConfig
	gu.CurrentFunction = gu.getGenerationFunction()
	return gu
}

func (gu *GenerationUnit) Generate() any {
	// check if important fields have changed and then regenerate the currentfunction
	if gu.configChanged(gu.PreviousValueConfig) || gu.UpdateCurrentFunction {
		gu.CurrentFunction = gu.getGenerationFunction()

		gu.UpdateCurrentFunction = false
	}
	if gu.generationConfig.valueGenerationType == GenerateOnce && gu.count > 0 {
		return gu.latestValue
	}
	gu.latestValue = gu.CurrentFunction.Generate()
	gu.PreviousValueConfig = gu.generationConfig.GenerationValueConfig
	gu.count++
	return gu.latestValue
}

func (gu *GenerationUnit) configChanged(previousConfig GenerationValueConfig) bool {

	return gu.generationConfig.valueGenerationType != previousConfig.valueGenerationType ||
		gu.generationConfig.setNonRequiredFields != previousConfig.setNonRequiredFields
}

func (gu *GenerationUnit) SetGenerationDefaultFunctionForKind(kind reflect.Kind, function GenerationFunction) {
	if gu.Field.Generator.DefaultGenerationFunctions[kind] == nil {
		return
	}
	gu.Field.Generator.DefaultGenerationFunctions[kind] = function
	gu.CurrentFunction = gu.getGenerationFunction()

}

func (gu *GenerationUnit) getGenerationFunction() GenerationFunction {

	switch gu.Field.Value.Kind() {
	case reflect.Slice:
		return GenerateSliceFunc(gu.Field)
	case reflect.Struct:
		return GenerateStructFunc(gu.Field)
	case reflect.Ptr:
		if gu.generationConfig.setNonRequiredFields {
			return GeneratePointerValueFunc(gu.Field)
		} else {
			return GenerateNilValueFunc()
		}
	}
	return gu.Field.getGenerationFunction()

}
