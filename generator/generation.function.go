package generator

type GenerationFunction interface {
	Generate() any
	// Copy copies the generation config and returns a copy of it with the same config
	// as the origin generation config
	// Copy(*GenerationConfig) GenerationFunction
	// GetGenerationConfig() *GenerationConfig
	// SetGenerationConfig(*GenerationConfig) GenerationFunction
}
