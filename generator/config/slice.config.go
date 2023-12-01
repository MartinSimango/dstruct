package config

type SliceConfig interface {
	SetLengthRange(min, max int) SliceConfig
	MinLength() int
	MaxLength() int
	Copy() SliceConfig
}

type SliceConfigImpl struct {
	minLength int
	maxLength int
}

var _ SliceConfig = &SliceConfigImpl{}

func (sc *SliceConfigImpl) SetLengthRange(min, max int) SliceConfig {
	if min > max {
		return sc
	}
	sc.minLength, sc.maxLength = min, max
	return sc

}

// MinLength implements SliceConfig.
func (sc *SliceConfigImpl) MinLength() int {
	return sc.minLength

}

// MaxLength implements SliceConfig.
func (sc *SliceConfigImpl) MaxLength() int {
	return sc.maxLength
}

func (sc *SliceConfigImpl) Copy() SliceConfig {
	s := *sc
	return &s
}

func NewSliceConfig() *SliceConfigImpl {
	return &SliceConfigImpl{0, 10}
}
