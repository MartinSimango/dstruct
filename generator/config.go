package generator

import "time"

type IntConfig struct {
	int64Min int64
	int64Max int64
	int32Min int32
	int32Max int32
	intMin   int
	intMax   int
}

type FloatConfig struct {
	float64Min float64
	float64Max float64
	float32Min float32
	float32Max float32
}

type SliceConfig struct {
	sliceMinLength int
	sliceMaxLength int
}

type DateConfig struct {
	dateFormat string
	dateStart  time.Time
	dateEnd    time.Time
}

func defaultIntConfig() IntConfig {
	return IntConfig{
		int64Min: 0,
		int64Max: 100,
		int32Min: 0,
		int32Max: 100,
		intMin:   0,
		intMax:   100,
	}
}

func defaultFloatConfig() FloatConfig {
	return FloatConfig{
		float64Min: 0,
		float64Max: 100,
		float32Min: 0,
		float32Max: 100,
	}
}

func defaultSliceConfig() SliceConfig {
	return SliceConfig{
		sliceMinLength: 0,
		sliceMaxLength: 10,
	}
}

func defaultDateConfig() DateConfig {
	return DateConfig{
		dateFormat: time.RFC3339,
	}
}

type ValueGenerationType uint8

const (
	Generate     ValueGenerationType = iota // will generate all field
	GenerateOnce                            // will generate all the fields
	UseDefaults
)

type GenerationValueConfig struct {
	valueGenerationType  ValueGenerationType
	setNonRequiredFields bool
}

type GenerationConfig struct {
	GenerationValueConfig
	DefaultGenerationFunctions GenerationDefaults
	SliceConfig
	IntConfig
	FloatConfig
	DateConfig
}

func (gc *GenerationConfig) SetNonRequiredFields(val bool) *GenerationConfig {
	gc.setNonRequiredFields = val
	return gc
}

func (gc *GenerationConfig) SetValueGenerationType(valueGenerationType ValueGenerationType) *GenerationConfig {
	gc.valueGenerationType = valueGenerationType
	return gc
}

func (gc *GenerationConfig) SetIntMax(max int) *GenerationConfig {
	if max >= gc.intMin {
		gc.intMax = max
	}
	return gc
}

func (gc *GenerationConfig) SetIntMin(min int) *GenerationConfig {
	if min <= gc.intMin {
		gc.intMin = min
	}
	return gc
}

func (gc *GenerationConfig) SetInt64Max(max int64) *GenerationConfig {
	if max >= gc.int64Min {
		gc.int64Max = max
	}
	return gc
}

func (gc *GenerationConfig) SetInt64Min(min int64) *GenerationConfig {
	if min <= gc.int64Max {
		gc.int64Min = min
	}
	return gc
}

func (gc *GenerationConfig) SetInt32Max(max int32) *GenerationConfig {
	if max >= gc.int32Min {
		gc.int32Max = max
	}
	return gc
}

func (gc *GenerationConfig) SetInt32Min(min int32) *GenerationConfig {
	if min <= gc.int32Max {
		gc.int32Min = min
	}
	return gc
}

func (gc *GenerationConfig) SetFloat64Max(max float64) *GenerationConfig {
	if max >= gc.float64Min {
		gc.float64Max = max
	}
	return gc
}

func (gc *GenerationConfig) SetFloat64Min(min float64) *GenerationConfig {
	if min <= gc.float64Max {
		gc.float64Min = min
	}
	return gc
}

func (gc *GenerationConfig) SetFloat32Max(max float32) *GenerationConfig {
	if max >= gc.float32Min {
		gc.float32Max = max
	}
	return gc
}

func (gc *GenerationConfig) SetFloat32Min(min float32) *GenerationConfig {
	if min <= gc.float32Max {
		gc.float32Min = min
	}
	return gc
}

func (gc *GenerationConfig) SetSliceLengthMax(max int) *GenerationConfig {
	if max >= gc.sliceMinLength {
		gc.sliceMaxLength = max
	}
	return gc
}

func (gc *GenerationConfig) SetSliceLengthMin(min int) *GenerationConfig {
	if min <= gc.sliceMaxLength {
		gc.sliceMinLength = min
	}
	return gc
}

func (gc *GenerationConfig) SetDateFormat(format string) *GenerationConfig {

	// TODO validate format
	gc.dateFormat = format
	return gc
}

func (gc *GenerationConfig) SetDateStart(start time.Time) *GenerationConfig {
	if !start.After(gc.dateEnd) {
		gc.dateStart = start
	}
	return gc
}

func (gc *GenerationConfig) SetDateEnd(end time.Time) *GenerationConfig {
	if end.After(gc.dateStart) {
		gc.dateStart = end
	}
	return gc
}
