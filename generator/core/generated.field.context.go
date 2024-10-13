package core

import (
	"github.com/MartinSimango/dstruct/generator"
	"github.com/MartinSimango/dstruct/generator/config"
)

type GeneratedFieldContext struct {
	PreviousGenerationSettings config.GenerationSettings
	CurrentFunction            generator.GenerationFunction
	GeneratedField             *GeneratedField

	count              int
	latestValue        any
	generationSettings config.GenerationSettings
}

func NewGeneratedFieldContext(field *GeneratedField) *GeneratedFieldContext {
	gfc := &GeneratedFieldContext{
		GeneratedField:     field,
		generationSettings: field.Config.GenerationSettings,
	}
	gfc.PreviousGenerationSettings = gfc.generationSettings
	// gu.CurrentFunction = gu.getGenerationFunction()
	return gfc
}

func (gfc *GeneratedFieldContext) Generate() any {
	//
	gfc.CurrentFunction = gfc.GeneratedField.getGenerationFunction()

	if gfc.configChanged(gfc.PreviousGenerationSettings) {
		gfc.count = 0
	}

	if gfc.generationSettings.ValueGenerationType == config.GenerateOnce && gfc.count > 0 {
		return gfc.latestValue
	}

	gfc.latestValue = gfc.CurrentFunction.Generate()
	gfc.PreviousGenerationSettings = gfc.generationSettings
	gfc.count++
	return gfc.latestValue
}

func (gfc *GeneratedFieldContext) SetGeneratedFieldConfig(cfg config.Config) {
	gfc.GeneratedField.SetConfig(cfg)
}

func (gfc *GeneratedFieldContext) SetGeneratedFieldSettings(settings config.GenerationSettings) {
	gfc.GeneratedField.SetGenerationSettings(settings)
}

// TODO: when does config need to be changed?

func (gu *GeneratedFieldContext) configChanged(previousConfig config.GenerationSettings) bool {
	return gu.generationSettings.ValueGenerationType != previousConfig.ValueGenerationType ||
		gu.generationSettings.SetNonRequiredFields != previousConfig.SetNonRequiredFields
}
