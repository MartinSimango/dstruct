package core

import (
	"github.com/MartinSimango/dstruct/generator"
	"github.com/MartinSimango/dstruct/generator/config"
)

type GenerationUnit struct {
	PreviousValueConfig   config.GenerationValueConfig
	CurrentFunction       generator.GenerationFunction
	UpdateCurrentFunction bool
	*GeneratedField
	count            int
	latestValue      any
	generationConfig config.GenerationValueConfig
}

func NewGenerationUnit(field *GeneratedField) *GenerationUnit {
	gu := &GenerationUnit{
		GeneratedField:   field,
		generationConfig: field.GeneratedFieldConfig.GenerationValueConfig,
	}
	gu.PreviousValueConfig = gu.generationConfig
	gu.CurrentFunction = gu.getGenerationFunction()
	return gu
}

func (gu *GenerationUnit) Generate() any {
	// check if important fields have changed and then regenerate the currentfunction
	gu.CurrentFunction = gu.getGenerationFunction()
	if gu.configChanged(gu.PreviousValueConfig) || gu.UpdateCurrentFunction {
		gu.CurrentFunction = gu.getGenerationFunction()
		gu.UpdateCurrentFunction = false
	}

	if gu.generationConfig.ValueGenerationType == config.GenerateOnce && gu.count > 0 {
		return gu.latestValue
	}
	gu.latestValue = gu.CurrentFunction.Generate()
	gu.PreviousValueConfig = gu.generationConfig
	gu.count++
	return gu.latestValue
}

func (gu *GenerationUnit) configChanged(previousConfig config.GenerationValueConfig) bool {

	return gu.generationConfig.ValueGenerationType != previousConfig.ValueGenerationType ||
		gu.generationConfig.SetNonRequiredFields != previousConfig.SetNonRequiredFields
}
