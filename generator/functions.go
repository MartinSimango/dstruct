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

var _ GenerationFunction = generateNumberFunc[int]{}

func (gn generateNumberFunc[N]) Copy(newConfig *GenerationConfig) GenerationFunction {
	return GenerateNumberFunc(*gn.args[0].(*N), *gn.args[1].(*N), newConfig)
}
func assign[T any](configMin, configMax *T, min, max T) (*T, *T) {
	*configMin = min
	*configMax = max
	return configMin, configMax
}

func GenerateNumberFunc[n number](min, max n, generationConfig *GenerationConfig) GenerationFunction {

	f := generateNumberFunc[n]{
		GenerationConfig:        generationConfig,
		basicGenerationFunction: generateNumber,
	}

	paramKind := reflect.ValueOf(new(n)).Elem().Kind()
	var param_1, param_2 any
	switch paramKind {
	case reflect.Int:
		param_1, param_2 = assign(&f.intMin, &f.intMax, int(min), int(max))
	case reflect.Int32:
		param_1, param_2 = assign(&f.int32Min, &f.int32Max, int32(min), int32(max))
	case reflect.Int64:
		param_1, param_2 = assign(&f.int64Min, &f.int64Max, int64(min), int64(max))
	case reflect.Float32:
		param_1, param_2 = assign(&f.float32Min, &f.float32Max, float32(min), float32(max))
	case reflect.Float64:
		param_1, param_2 = assign(&f.float64Min, &f.float64Max, float64(min), float64(max))
	default:
		panic(fmt.Sprintf("Invalid number type: %s", paramKind))
	}
	f.args = []any{param_1, param_2}
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

	f := generateStructFunc{
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
