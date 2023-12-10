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
		WithDateConfig(NewDateConfig()).
		Build()
}

type Config interface {
	ConfigBuilder
	Number() NumberConfig
	Slice() SliceConfig
	Date() DateConfig
	SetSliceLength(min, max int) Config
	SetIntRange(min, max int) Config
	SetInt8Range(min, max int8) Config
	SetInt16Range(min, max int16) Config
	SetInt32Range(min, max int32) Config
	SetInt64Range(min, max int64) Config
	SetFloat32Range(min, max float32) Config
	SetFloat64Range(min, max float64) Config
	SetUIntRange(min, max uint) Config
	SetUInt8Range(min, max uint8) Config
	SetUInt16Range(min, max uint16) Config
	SetUInt32Range(min, max uint32) Config
	SetUInt64Range(min, max uint64) Config
	SetUIntPtr(min, max uintptr) Config

	Copy() Config
}

type ConfigImpl struct {
	ConfigBuilderImpl
	SliceConfig  SliceConfig
	NumberConfig NumberConfig
	DateConfig   DateConfig
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

func (c *ConfigImpl) Date() DateConfig {
	return c.DateConfig

}

func (c *ConfigImpl) Number() NumberConfig {
	return c.NumberConfig
}

func (c *ConfigImpl) SetSliceLength(min, max int) Config {
	c.SliceConfig.SetLengthRange(min, max)
	return c
}

func (c *ConfigImpl) SetIntRange(min, max int) Config {
	c.NumberConfig.Int().SetRange(min, max)
	return c
}

func (c *ConfigImpl) SetInt8Range(min, max int8) Config {
	c.NumberConfig.Int8().SetRange(min, max)
	return c
}

func (c *ConfigImpl) SetInt16Range(min, max int16) Config {
	c.NumberConfig.Int16().SetRange(min, max)
	return c
}

func (c *ConfigImpl) SetInt32Range(min, max int32) Config {
	c.NumberConfig.Int32().SetRange(min, max)
	return c
}

func (c *ConfigImpl) SetInt64Range(min, max int64) Config {
	c.NumberConfig.Int64().SetRange(min, max)
	return c
}

func (c *ConfigImpl) SetFloat32Range(min, max float32) Config {
	c.NumberConfig.Float32().SetRange(min, max)
	return c
}
func (c *ConfigImpl) SetFloat64Range(min, max float64) Config {
	c.NumberConfig.Float64().SetRange(min, max)
	return c
}

func (c *ConfigImpl) SetUIntRange(min, max uint) Config {
	c.NumberConfig.UInt().SetRange(min, max)
	return c
}

func (c *ConfigImpl) SetUInt8Range(min, max uint8) Config {
	c.NumberConfig.UInt8().SetRange(min, max)
	return c
}

func (c *ConfigImpl) SetUInt16Range(min, max uint16) Config {
	c.NumberConfig.UInt16().SetRange(min, max)
	return c
}

func (c *ConfigImpl) SetUInt32Range(min, max uint32) Config {
	c.NumberConfig.UInt32().SetRange(min, max)
	return c
}

func (c *ConfigImpl) SetUInt64Range(min, max uint64) Config {
	c.NumberConfig.UInt64().SetRange(min, max)
	return c
}

func (c *ConfigImpl) SetUIntPtr(min, max uintptr) Config {
	c.NumberConfig.UIntPtr().SetRange(min, max)
	return c
}
