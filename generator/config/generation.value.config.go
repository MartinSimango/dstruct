package config

type ConfigType uint

type ValueGenerationType uint8

const (
	Generate     ValueGenerationType = iota // will generate all field
	GenerateOnce                            // will generate the fields once
	UseDefaults
)

// RecursiveDefinition defines the configuration for recursive definitions within a dynamic struct.
type RecursiveDefinition struct {
	// Allow defines if recursive definitions are allowed.
	Allow bool
	// Depth defines the depth of the recursive definition if recursion is allowed.
	Depth uint
}

// GenerationValueConfig defines how values are generated within a dynamic struct.
type GenerationValueConfig struct {
	ValueGenerationType  ValueGenerationType
	SetNonRequiredFields bool
	RecursiveDefinition  RecursiveDefinition
}

// DefaultGenerationValueConfig returns a default configuration for generation values.
func DefaultGenerationValueConfig() GenerationValueConfig {
	return GenerationValueConfig{
		ValueGenerationType:  UseDefaults,
		SetNonRequiredFields: false,
		RecursiveDefinition: RecursiveDefinition{
			Allow: false,
			Depth: 1,
		},
	}
}
