package config

type ConfigType uint

type ValueGenerationType uint8

const (
	// Generate will generated all fields.
	Generate ValueGenerationType = iota
	// GenerateOnce will generate the fields once.
	GenerateOnce
	// UseDefaults will use the default value of the field based on the default tag. If the default tag is not set, the
	// field will use the field type's default generation function to generate a value for the field.
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


func (gs *GenerationSettings) WithValueGenerationType(
	valueGenerationType ValueGenerationType,
) GenerationSettings {
	gs.ValueGenerationType = valueGenerationType
	return *gs
}

func (gs *GenerationSettings) WithRecursiveDefinition(
	recursiveDefinition RecursiveDefinition,
) GenerationSettings {
	gs.RecursiveDefinition = recursiveDefinition
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
