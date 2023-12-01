package config

type ConfigType uint

type ValueGenerationType uint8

const (
	Generate     ValueGenerationType = iota // will generate all field
	GenerateOnce                            // will generate the fields once
	UseDefaults
)

type RecursiveDefinition struct {
	Allow bool
	Depth uint
}

type GenerationValueConfig struct {
	ValueGenerationType  ValueGenerationType
	SetNonRequiredFields bool
	RecursiveDefinition  RecursiveDefinition
}

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

func NewConfig() *ConfigImpl {
	return NewConfigBuilder().
		WithNumberConfig(NewNumberConfig()).
		WithSliceConfig(NewSliceConfig()).
		Build()
}

type Config interface {
	ConfigBuilder
	Number() NumberConfig
	Slice() SliceConfig
	SetSliceLength(min, max int) Config
	Copy() Config
}

type ConfigImpl struct {
	ConfigBuilderImpl
	SliceConfig  SliceConfig
	NumberConfig NumberConfig
}

var _ Config = &ConfigImpl{}

func (c *ConfigImpl) Copy() Config {
	newConfig := &ConfigImpl{}
	if c.SliceConfig != nil {
		newConfig.SliceConfig = c.SliceConfig.Copy()
	}

	if c.NumberConfig != nil {
		newConfig.NumberConfig = c.NumberConfig.Copy()
	}

	return newConfig

}

func (c *ConfigImpl) Slice() SliceConfig {
	return c.SliceConfig

}

func (c *ConfigImpl) Number() NumberConfig {
	return c.NumberConfig
}

func (c *ConfigImpl) SetSliceLength(min, max int) Config {
	c.SliceConfig.SetLengthRange(min, max)
	return c
}
