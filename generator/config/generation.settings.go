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

// GenerationSettings defines how values are generated within a dynamic struct.
type GenerationSettings struct {
	ValueGenerationType  ValueGenerationType
	SetNonRequiredFields bool
	RecursiveDefinition  RecursiveDefinition
}

func (gs *GenerationSettings) WithNonRequiredFields(required bool) GenerationSettings {
	gs.SetNonRequiredFields = required
	return *gs
}

// DefaultGenerationSettings returns a default configuration for generation values.
func DefaultGenerationSettings() GenerationSettings {
	return GenerationSettings{
		ValueGenerationType:  UseDefaults,
		SetNonRequiredFields: false,
		RecursiveDefinition: RecursiveDefinition{
			Allow: false,
			Depth: 1,
		},
	}
}
