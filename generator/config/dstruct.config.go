package config

// DstructConfig defines the configuration for a dynamic struct.
type DstructConfig struct {
	SliceConfig  SliceConfig
	NumberConfig NumberRangeConfig
	DateConfig   DateRangeConfig
	children     []Config
}

// NewDstructConfig is a constructor for DstructConfig.
func NewDstructConfig() *DstructConfig {
	return NewDstructConfigBuilder().
		WithNumberConfig(NewNumberRangeConfig()).
		WithSliceConfig(NewSliceConfig()).
		WithDateRangeConfig(NewDstructDateRangeConfig()).
		Build()
}

var _ Config = &DstructConfig{}

// Copy implements Config.
// TODO: implement this
func (c *DstructConfig) Copy() Config {
	newConfig := &DstructConfig{}
	if c.SliceConfig != nil {
		newConfig.SliceConfig = c.SliceConfig.Copy()
	}

	if c.NumberConfig != nil {
		newConfig.NumberConfig = c.NumberConfig.Copy()
	}

	return newConfig
}

// Slice implements Config.Slice.
func (c *DstructConfig) Slice() SliceConfig {
	return c.SliceConfig
}

// Date implements Config.Date.
func (c *DstructConfig) Date() DateRangeConfig {
	return c.DateConfig
}

// Number implements Config.Number.
func (c *DstructConfig) Number() NumberRangeConfig {
	return c.NumberConfig
}

// SetSliceLength implements Config.SetSliceLength.
func (c *DstructConfig) SetSliceLength(min, max int) Config {
	c.SliceConfig.SetLengthRange(min, max)
	for _, sub := range c.children {
		sub.SetSliceLength(min, max)
	}
	return c
}

// SetIntRange implements Config.SetIntRange.
func (c *DstructConfig) SetIntRange(min, max int) Config {
	c.NumberConfig.Int().SetRange(min, max)
	return c
}

// SetInt8Range implements Config.SetInt8Range.
func (c *DstructConfig) SetInt8Range(min, max int8) Config {
	c.NumberConfig.Int8().SetRange(min, max)
	return c
}

// SetInt16Range implements Config.SetInt16Range.
func (c *DstructConfig) SetInt16Range(min, max int16) Config {
	c.NumberConfig.Int16().SetRange(min, max)
	return c
}

// SetInt32Range implements Config.SetInt32Range.
func (c *DstructConfig) SetInt32Range(min, max int32) Config {
	c.NumberConfig.Int32().SetRange(min, max)
	return c
}

// SetInt64Range implements Config.SetInt64Range.
func (c *DstructConfig) SetInt64Range(min, max int64) Config {
	c.NumberConfig.Int64().SetRange(min, max)
	return c
}

// SetFloat32Range implements Config.SetFloat32Range.
func (c *DstructConfig) SetFloat32Range(min, max float32) Config {
	c.NumberConfig.Float32().SetRange(min, max)
	return c
}

// SetFloat64Range implements Config.SetFloat64Range.
func (c *DstructConfig) SetFloat64Range(min, max float64) Config {
	c.NumberConfig.Float64().SetRange(min, max)
	return c
}

// SetUIntRange implements Config.SetUIntRange.
func (c *DstructConfig) SetUIntRange(min, max uint) Config {
	c.NumberConfig.UInt().SetRange(min, max)
	return c
}

// SetUInt8Range implements Config.SetUInt8Range.
func (c *DstructConfig) SetUInt8Range(min, max uint8) Config {
	c.NumberConfig.UInt8().SetRange(min, max)
	return c
}

// SetUInt16Range implements Config.SetUInt16Range.
func (c *DstructConfig) SetUInt16Range(min, max uint16) Config {
	c.NumberConfig.UInt16().SetRange(min, max)
	return c
}

// SetUInt32Range implements Config.SetUInt32Range.
func (c *DstructConfig) SetUInt32Range(min, max uint32) Config {
	c.NumberConfig.UInt32().SetRange(min, max)
	return c
}

// SetUInt64Range implements Config.SetUInt64Range.
func (c *DstructConfig) SetUInt64Range(min, max uint64) Config {
	c.NumberConfig.UInt64().SetRange(min, max)
	return c
}

// SetUIntPtrRange implements Config.SetUIntPtrRange.
func (c *DstructConfig) SetUIntPtr(min, max uintptr) Config {
	c.NumberConfig.UIntPtr().SetRange(min, max)
	return c
}
