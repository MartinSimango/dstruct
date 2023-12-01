package config

type ConfigBuilder interface {
	WithNumberConfig(NumberConfig) ConfigBuilder
	WithSliceConfig(SliceConfig) ConfigBuilder
	Build() *ConfigImpl
}

type ConfigBuilderImpl struct {
	config *ConfigImpl
}

var _ ConfigBuilder = &ConfigBuilderImpl{}

func NewConfigBuilder() *ConfigBuilderImpl {
	return &ConfigBuilderImpl{
		config: &ConfigImpl{},
	}
}

// WithSliceConfig implements ConfigBuilder.
func (cb *ConfigBuilderImpl) WithNumberConfig(nc NumberConfig) ConfigBuilder {
	cb.config.NumberConfig = nc
	return cb
}

// WithSliceConfig implements ConfigBuilder.
func (cb *ConfigBuilderImpl) WithSliceConfig(sc SliceConfig) ConfigBuilder {
	cb.config.SliceConfig = sc
	return cb
}

func (cb *ConfigBuilderImpl) Build() *ConfigImpl {
	return cb.config
}
