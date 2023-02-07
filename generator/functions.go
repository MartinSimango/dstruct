package generator

import (
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"time"

	"github.com/takahiromiyamoto/go-xeger"
)

const ISO8601 string = "2018-03-20T09:12:28Z"

func init() {
	rand.Seed(time.Now().Unix())
}

type number interface {
	int8 | int | int32 | int64 | float32 | float64
}

func generateNum[n number](min, max n) n {
	return min + (n(rand.Float64() * float64(max+1-min)))
}

var (
	generateStringFromRegex GenerationFunc = GenerationFunc{
		_func: func(parameters ...any) any {
			regex := parameters[0].(string)
			x, err := xeger.NewXeger(regex)
			if err != nil {
				panic(err)
			}
			return x.Generate()
		},
	}

	generateNumber GenerationFunc = GenerationFunc{
		_func: func(parameters ...any) any {
			min := parameters[0]
			max := parameters[1]
			paramKind := reflect.ValueOf(min).Elem().Kind()

			switch paramKind {
			case reflect.Int:
				return generateNum(*min.(*int), *max.(*int))
			case reflect.Int32:
				return generateNum(*min.(*int32), *max.(*int32))
			case reflect.Int64:
				return generateNum(*min.(*int64), *max.(*int64))
			case reflect.Float32:
				return generateNum(*min.(*float32), *max.(*float32))
			case reflect.Float64:
				return generateNum(*min.(*float64), *max.(*float64))
			default:
				panic(fmt.Sprintf("Invalid number type: %s", paramKind))

			}
		},
	}

	generateBool GenerationFunc = GenerationFunc{
		_func: func(parameters ...any) any {
			return generateNum(0, 1) == 0
		},
	}

	generateObject GenerationFunc = GenerationFunc{
		_func: func(parameters ...any) any {
			return nil
		},
	}

	generateNilValue GenerationFunc = GenerationFunc{
		_func: func(parameters ...any) any {
			return nil
		},
	}

	generateStruct GenerationFunc = GenerationFunc{

		_func: func(parameters ...any) any {
			generationConfig := parameters[0].(*GenerationConfig)
			field := parameters[1].(Field)
			setStructValues(field, generationConfig)
			return field.Value
		},
	}

	generateSlice GenerationFunc = GenerationFunc{
		_func: func(parameters ...any) any {
			generationConfig := parameters[0].(*GenerationConfig)

			field := parameters[1].(Field)
			sliceType := reflect.TypeOf(field.Value.Interface()).Elem()
			min := generationConfig.sliceMinLength
			max := generationConfig.sliceMaxLength

			len := min + (int(rand.Float64() * float64(max+1-min)))
			sliceOfElementType := reflect.SliceOf(sliceType)
			slice := reflect.MakeSlice(sliceOfElementType, 0, 1024)

			switch sliceType.Kind() {
			case reflect.Struct:
				sliceElement := reflect.New(sliceType)
				for i := 0; i < len; i++ {
					newField := Field{
						Name:  field.Name,
						Value: reflect.ValueOf(sliceElement.Interface()).Elem(),
						Tag:   field.Tag,
					}
					setStructValues(newField, generationConfig)

					slice = reflect.Append(slice, sliceElement.Elem())

				}
			}

			return slice.Interface()

		},
	}

	generatePointerValue GenerationFunc = GenerationFunc{

		_func: func(parameters ...any) any {
			generationConfig := parameters[0].(*GenerationConfig)
			field := parameters[1].(Field)
			ptr := reflect.New(field.Value.Type().Elem())
			field.Value = ptr.Elem()
			setValue(field, generationConfig)
			return ptr.Interface()

		},
	}

	generateDateTime GenerationFunc = GenerationFunc{
		_func: func(parameters ...any) any {
			return time.Now().UTC().Format(time.RFC3339)
		},
	}
)

func GenerateStringFromRegexFunc(regex string) GenerationFunction {
	f := generateStringFromRegex
	f.args = []any{regex}
	return f
}

func GenerateNumberFunc[n number](min, max *n) GenerationFunction {
	f := generateNumber
	f.args = []any{min, max}
	return f
}

func GenerateFixedValueFunc[T any](n T) GenerationFunction {
	var f GenerationFunc
	f._func = func(p ...any) any {
		return n
	}
	return f
}

func GenerateBoolFunc() GenerationFunction {
	f := generateBool
	return f
}

func GenerateNilValueFunc() GenerationFunction {
	f := generateNilValue
	return f
}

func GenerateSliceFunc(generationConfig *GenerationConfig, field Field) GenerationFunction {
	f := generateSlice
	f.args = []any{generationConfig, field}
	return f
}

func GenerateStructFunc(generationConfig *GenerationConfig, field Field) GenerationFunction {
	f := generateStruct
	f.args = []any{generationConfig, field}
	return f
}

func GeneratePointerValueFunc(generationConfig *GenerationConfig, field Field) GenerationFunction {
	f := generatePointerValue
	f.args = []any{generationConfig, field}
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

func generationFunctionFromTags(field Field,
	generationConfig *GenerationConfig) GenerationFunction {
	kind := field.Value.Kind()
	tags := field.Tag

	if generationConfig.valueGenerationType == UseDefaults {
		example, ok := tags.Lookup("example")
		if !ok {
			example, ok = tags.Lookup("default")
		}
		if ok {
			switch kind {
			case reflect.Int:
				v, _ := strconv.Atoi(example)
				return GenerateFixedValueFunc(v)
			case reflect.Int32:
				v, _ := strconv.Atoi(example)
				return GenerateFixedValueFunc(int32(v))
			case reflect.Int64:
				v, _ := strconv.Atoi(example)
				return GenerateFixedValueFunc(int64(v))
			case reflect.Float64:
				v, _ := strconv.ParseFloat(example, 64)
				return GenerateFixedValueFunc(float64(v))
			case reflect.String:
				return GenerateFixedValueFunc(example)
			case reflect.Bool:
				v, _ := strconv.ParseBool(example)
				return GenerateFixedValueFunc(v)
			default:
				fmt.Println("Unsupported types for defaults: ", kind, example)
			}

		}
	}

	pattern := tags.Get("pattern")
	if pattern != "" {
		return GenerateStringFromRegexFunc(pattern)
	}

	format := tags.Get("format")

	switch format {
	case "date-time":
		return GenerateDateTimeFunc()
	}

	enum, ok := tags.Lookup("enum")
	if ok {
		numEnums, _ := strconv.Atoi(enum)
		return GenerateFixedValueFunc(tags.Get(fmt.Sprintf("enum_%d", generateNum(0, numEnums-1)+1)))
	}

	gen_task, ok := tags.Lookup("gen_task")
	if ok {
		return getTask(gen_task, field.Name).getFunction()
	}
	return generationConfig.DefaultGenerationFunctions[kind]
}

type Field struct {
	Name  string
	Value reflect.Value
	Tag   reflect.StructTag
}
