package generator

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
	gu.CurrentFunction = gu.Field.getGenerationFunction()
	return gu
}

func (gu *GenerationUnit) Generate() any {
	// check if important fields have changed and then regenerate the currentfunction
	if gu.configChanged(gu.PreviousValueConfig) || gu.UpdateCurrentFunction {
		gu.CurrentFunction = gu.Field.getGenerationFunction()
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
