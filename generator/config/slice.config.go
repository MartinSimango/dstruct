package config

type SliceConfig interface {

	// SetLengthRange sets the length range for the slice.
	SetLengthRange(min, max int) SliceConfig

	// MinLength returns the minimum length of the slice.
	MinLength() int

	// MaxLength returns the maximum length of the slice.
	MaxLength() int

	// Copy returns a copy of the slice configuration.
	Copy() SliceConfig
}
