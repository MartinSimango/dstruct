package generator

// var _ GenerationFunction = &generateNumberFunc[int]{}

// func (gn generateNumberFunc[N]) Copy(newConfig *FunctionGenerationConfig) GenerationFunction {
// 	return GenerateNumberFunc(*gn.args[0].(*N), *gn.args[1].(*N)).SetGenerationConfig(newConfig)
// }
// func assign[T any](configMin, configMax *T, min, max T) (*T, *T) {
// 	*configMin = min
// 	*configMax = max
// 	return configMin, configMax
// }

// func (gf *generateNumberFunc[n]) GetGenerationConfig() *FunctionGenerationConfig {
// 	return gf.FunctionGenerationConfig
// }

// func (gf *generateNumberFunc[n]) SetGenerationConfig(generationConfig *FunctionGenerationConfig) GenerationFunction {
// 	min, max := *gf.args[0].(*n), *gf.args[1].(*n)
// 	paramKind := reflect.ValueOf(min).Kind()
// 	var param_1, param_2 any
// 	switch paramKind {
// 	case reflect.Int:
// 		param_1, param_2 = assign(&generationConfig.intMin, &generationConfig.intMax, int(min), int(max))
// 	case reflect.Int32:
// 		param_1, param_2 = assign(&generationConfig.int32Min, &generationConfig.int32Max, int32(min), int32(max))
// 	case reflect.Int64:
// 		param_1, param_2 = assign(&generationConfig.int64Min, &generationConfig.int64Max, int64(min), int64(max))
// 	case reflect.Float32:
// 		param_1, param_2 = assign(&generationConfig.float32Min, &generationConfig.float32Max, float32(min), float32(max))
// 	case reflect.Float64:
// 		param_1, param_2 = assign(&generationConfig.float64Min, &generationConfig.float64Max, float64(min), float64(max))
// 	default:
// 		panic(fmt.Sprintf("Invalid number type: %s", paramKind))
// 	}
// 	gf.args = []any{param_1, param_2}
// 	return gf
// }

// func (g generateStructFunc) Copy(newConfig *FunctionGenerationConfig) GenerationFunction {
// 	return GenerateStructFunc(g.field)
// }
