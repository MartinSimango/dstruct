package config

// ConfigBuilder is an interface of a DstructConfig builder.
type DstructConfigBuilder interface {

	// WithNumberConfig sets the NumberConfig in the builder.
	WithNumberConfig(numberConfig NumberRangeConfig) DstructConfigBuilder

	// WithSliceConfig sets the SliceConfig in the builder.
	WithSliceConfig(sliceConfig SliceConfig) DstructConfigBuilder

	// WithDateRangeConfig sets the DateConfig in the builder.
	WithDateRangeConfig(DateRangeConfig DateRangeConfig) DstructConfigBuilder

	// Build builds the configuration object
	Build() *DstructConfig
}

type dstructConfigBuilder struct {
	config *DstructConfig
}

var _ DstructConfigBuilder = &dstructConfigBuilder{}

// NewDstructConfigBuilder is a constructor for DstructConfigBuilder.
func NewDstructConfigBuilder() *dstructConfigBuilder {
	return &dstructConfigBuilder{
		config: &DstructConfig{},
	}
}

// WithNumberConfig implements ConfigBuilder.WithNumberConfig.
func (dcb *dstructConfigBuilder) WithNumberConfig(nc NumberRangeConfig) DstructConfigBuilder {
	dcb.config.NumberConfig = nc
	return dcb
}

// WithSliceConfig implements ConfigBuilder.WithSliceConfig.
func (dcb *dstructConfigBuilder) WithSliceConfig(sc SliceConfig) DstructConfigBuilder {
	dcb.config.SliceConfig = sc
	return dcb
}

// WithDateRangeConfig implements ConfigBuilder.WithDateRangeConfig.
func (dcb *dstructConfigBuilder) WithDateRangeConfig(dc DateRangeConfig) DstructConfigBuilder {
	dcb.config.DateConfig = dc
	return dcb
}

// Build implements ConfigBuilder.Build.
func (dcb *dstructConfigBuilder) Build() *DstructConfig {
	return dcb.config
}
