package generator

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/takahiromiyamoto/go-xeger"
)

type basicGenerationFunction struct {
	_func func(...any) any
	args  []any
}

const ISO8601 string = "2018-03-20T09:12:28Z"

type number interface {
	int8 | int | int32 | int64 | float32 | float64
}

func generateNum[n number](min, max n) n {
	return min + (n(rand.Float64() * float64(max+1-min)))
}

var (
	generateStringFromRegex basicGenerationFunction = basicGenerationFunction{
		_func: func(parameters ...any) any {
			regex := parameters[0].(string)
			x, err := xeger.NewXeger(regex)
			if err != nil {
				panic(err)
			}
			return x.Generate()
		},
	}

	generateNumber basicGenerationFunction = basicGenerationFunction{
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

	generateBool basicGenerationFunction = basicGenerationFunction{
		_func: func(parameters ...any) any {
			return generateNum(0, 1) == 0
		},
	}

	generateObject basicGenerationFunction = basicGenerationFunction{
		_func: func(parameters ...any) any {
			return nil
		},
	}

	generateNilValue basicGenerationFunction = basicGenerationFunction{
		_func: func(parameters ...any) any {
			return nil
		},
	}

	generateStruct basicGenerationFunction = basicGenerationFunction{

		_func: func(parameters ...any) any {
			field := parameters[0].(*GeneratedField)
			field.setStructValues()
			return field.Value
		},
	}

	generateSlice basicGenerationFunction = basicGenerationFunction{
		_func: func(parameters ...any) any {

			field := parameters[0].(*GeneratedField)
			generationConfig := field.Generator.GenerationConfig
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
					newField := GeneratedField{
						Name:      field.Name,
						Value:     reflect.ValueOf(sliceElement.Interface()).Elem(),
						Tag:       field.Tag,
						Generator: field.Generator.Copy(),
					}
					newField.setStructValues()

					slice = reflect.Append(slice, sliceElement.Elem())

				}
			}

			return slice.Interface()

		},
	}

	generatePointerValue basicGenerationFunction = basicGenerationFunction{

		_func: func(parameters ...any) any {
			field := parameters[0].(*GeneratedField)
			ptr := reflect.New(field.Value.Type().Elem())
			field.Value = ptr.Elem()
			field.SetValue()
			return ptr.Interface()

		},
	}

	generateDateTime basicGenerationFunction = basicGenerationFunction{
		_func: func(parameters ...any) any {
			return time.Now().UTC().Format(time.RFC3339)
		},
	}
)
