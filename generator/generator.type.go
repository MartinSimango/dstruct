package generator

// type GeneratorType interface {
// 	GetGenerationConfig() *GenerationConfig
// 	GetGenerationFunction(field *GeneratedField) GenerationFunction
// }

// type StructGenerator struct {
// 	DefaultGenerationFunctions DefaultGenerationFunctionType
// 	GenerationConfig           *GenerationConfig
// }

// func NewStructGenerator(generationConfig *GenerationConfig) *StructGenerator {
// 	return &StructGenerator{
// 		GenerationConfig: generationConfig,
// 	}
// }

// func (sg *StructGenerator) GetGenerationFunction(field *GeneratedField) GenerationFunction {
// 	return GenerateStructFunc(sg.Field, sg.GenerationConfig)
// }

// func (sg *StructGenerator) GetGenerationConfig() *GenerationConfig {
// 	return sg.GenerationConfig
// }

// func (sg *StructGenerator) SetDefaultGenerationFunctions() {

// }

// type FieldGenerator struct {
// 	GeneratorFunction GenerationFunction
// }

// type ArrayGenerator struct {
// }

// type SliceGenerator struct {
// 	GenerationConfig
// }

// // switch gu.Field.Value.Kind() {
// // case reflect.Slice:
// // 	return GenerateSliceFunc(gu.Field)
// // case reflect.Struct:
// // 	return GenerateStructFunc(gu.Field)
// // case reflect.Ptr:
// // 	if gu.generationConfig.setNonRequiredFields {
// // 		return GeneratePointerValueFunc(gu.Field)
// // 	} else {
// // 		return GenerateNilValueFunc()
// // 	}
// // }
