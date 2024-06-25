package config

// DstructSliceConfig implements SliceConfig.
type DstructSliceConfig struct {
	minLength int
	maxLength int
}

var _ SliceConfig = &DstructSliceConfig{}

// SetLengthRange implements SliceConfig.SetLengthRange.
func (sc *DstructSliceConfig) SetLengthRange(min, max int) SliceConfig {
	if min > max {
		return sc
	}
	sc.minLength, sc.maxLength = min, max
	return sc

}

// MinLength implements SliceConfig.MinLength.
func (sc *DstructSliceConfig) MinLength() int {
	return sc.minLength

}

// MaxLength implements SliceConfig.MaxLength.
func (sc *DstructSliceConfig) MaxLength() int {
	return sc.maxLength
}

// Copy implements SliceConfig.Copy.
func (sc *DstructSliceConfig) Copy() SliceConfig {
	s := *sc
	return &s
}

// NewSliceConfig is a constructor for DstructSliceConfig.
func NewSliceConfig() *DstructSliceConfig {
	return &DstructSliceConfig{0, 10}
}
