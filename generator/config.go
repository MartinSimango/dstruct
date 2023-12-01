package generator

// import (
// 	"time"
// )

// type DateConfig struct {
// 	dateFormat string
// 	dateStart  time.Time
// 	dateEnd    time.Time
// }

// func defaultDateConfig() DateConfig {
// 	return DateConfig{
// 		dateFormat: time.RFC3339,
// 	}
// }

// func NewGenerationConfig() (generationConfig *FunctionGenerationConfig) {
// 	generationConfig = &FunctionGenerationConfig{
// 		GenerationValueConfig: GenerationValueConfig{
// 			valueGenerationType:  UseDefaults,
// 			setNonRequiredFields: false,
// 			recursiveDefinition: recursiveDefinition{
// 				Allow: false,
// 				Count: 1,
// 			},
// 		},
// 		SliceConfig: defaultSliceConfig(),
// 		IntConfig:   defaultIntConfig(),
// 		FloatConfig: defaultFloatConfig(),
// 		DateConfig:  defaultDateConfig(),
// 	}
// 	return
// }

// func (gc *FunctionGenerationConfig) Clone() (config *FunctionGenerationConfig) {
// 	config = &FunctionGenerationConfig{}
// 	*config = *gc
// 	return
// }

// func (gc *FunctionGenerationConfig) AllowRecursion(a bool) *FunctionGenerationConfig {
// 	gc.recursiveDefinition.Allow = a
// 	return gc
// }
// func (gc *FunctionGenerationConfig) SetRecursionCount(r uint) *FunctionGenerationConfig {
// 	gc.recursiveDefinition.Count = r
// 	return gc
// }

// func (gc *FunctionGenerationConfig) SetNonRequiredFields(val bool) *FunctionGenerationConfig {
// 	gc.setNonRequiredFields = val
// 	return gc
// }

// func (gc *FunctionGenerationConfig) SetValueGenerationType(valueGenerationType ValueGenerationType) *FunctionGenerationConfig {
// 	gc.valueGenerationType = valueGenerationType
// 	return gc
// }

// func (gc *FunctionGenerationConfig) SetDateFormat(format string) *FunctionGenerationConfig {

// 	// TODO validate format
// 	gc.dateFormat = format
// 	return gc
// }

// func (gc *FunctionGenerationConfig) SetDateStart(start time.Time) *FunctionGenerationConfig {
// 	if !start.After(gc.dateEnd) {
// 		gc.dateStart = start
// 	}
// 	return gc
// }

// func (gc *FunctionGenerationConfig) SetDateEnd(end time.Time) *FunctionGenerationConfig {
// 	if end.After(gc.dateStart) {
// 		gc.dateStart = end
// 	}
// 	return gc
// }
