package generator

import (
	"fmt"
	"reflect"
	"time"
)

type (
	generateNumberFunc[N number]  GenerationFunctionImpl
	generateStringFromRegexFunc   GenerationFunctionImpl
	generateFixedValueFunc[T any] GenerationFunctionImpl
)

var _ GenerationFunction = &generateNumberFunc[int]{}

func (gn generateNumberFunc[N]) Copy(newConfig *GenerationConfig) GenerationFunction {
	return GenerateNumberFunc(*gn.args[0].(*N), *gn.args[1].(*N)).SetGenerationConfig(newConfig)
}
func assign[T any](configMin, configMax *T, min, max T) (*T, *T) {
	*configMin = min
	*configMax = max
	return configMin, configMax
}

func (gf *generateNumberFunc[n]) GetGenerationConfig() *GenerationConfig {
	return gf.GenerationConfig
}

func (gf *generateNumberFunc[n]) SetGenerationConfig(generationConfig *GenerationConfig) GenerationFunction {
	min, max := *gf.args[0].(*n), *gf.args[1].(*n)
	paramKind := reflect.ValueOf(min).Kind()
	var param_1, param_2 any
	switch paramKind {
	case reflect.Int:
		param_1, param_2 = assign(&generationConfig.intMin, &generationConfig.intMax, int(min), int(max))
	case reflect.Int32:
		param_1, param_2 = assign(&generationConfig.int32Min, &generationConfig.int32Max, int32(min), int32(max))
	case reflect.Int64:
		param_1, param_2 = assign(&generationConfig.int64Min, &generationConfig.int64Max, int64(min), int64(max))
	case reflect.Float32:
		param_1, param_2 = assign(&generationConfig.float32Min, &generationConfig.float32Max, float32(min), float32(max))
	case reflect.Float64:
		param_1, param_2 = assign(&generationConfig.float64Min, &generationConfig.float64Max, float64(min), float64(max))
	default:
		panic(fmt.Sprintf("Invalid number type: %s", paramKind))
	}
	gf.args = []any{param_1, param_2}
	return gf
}

func GenerateNumberFunc[n number](min, max n) GenerationFunction {
	f := &generateNumberFunc[n]{
		basicGenerationFunction: generateNumber,
	}

	f.args = []any{&min, &max}

	return f
}

func GenerateStringFromRegexFunc(regex string) GenerationFunction {
	f := generateStringFromRegexFunc{
		basicGenerationFunction: generateStringFromRegex,
	}
	f.args = []any{regex}
	return f
}

func GenerateFixedValueFunc[T any](n T) GenerationFunction {
	f := generateFixedValueFunc[T]{}
	f._func = func(p ...any) any {
		return n
	}
	return f
}

func GenerateBoolFunc() basicGenerationFunction {
	f := generateBool
	return f
}

func GenerateNilValueFunc() basicGenerationFunction {
	f := generateNilValue
	return f
}

func GenerateSliceFunc(field *GeneratedField) GenerationFunction {
	f := generateSlice
	f.args = []any{field}
	return f
}

type generateStructFunc struct {
	GenerationFunctionImpl
	field *GeneratedField
}

func GenerateStructFunc(field *GeneratedField) GenerationFunction {

	f := &generateStructFunc{
		GenerationFunctionImpl: GenerationFunctionImpl{
			GenerationConfig:        field.Generator.GenerationConfig,
			basicGenerationFunction: generateStruct,
		},
		field: field,
	}
	f.args = []any{field}
	return f
}

func (g generateStructFunc) Copy(newConfig *GenerationConfig) GenerationFunction {
	return GenerateStructFunc(g.field)
}

func GeneratePointerValueFunc(field *GeneratedField) GenerationFunction {
	f := generatePointerValue
	f.args = []any{field}
	return f
}

func GenerateDateTimeFunc() GenerationFunction {
	f := generateDateTime
	return f
}
func GenerateDateTimeBetweenDatesFunc(startDate, endDate time.Time) GenerationFunction {
	// TODO implement
	f := generateDateTime
	return f
}
