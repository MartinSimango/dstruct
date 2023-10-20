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
	// TODO consider making args exportable for more customizable generation functions
	args []any
}

const ISO8601 string = "2018-03-20T09:12:28Z"

type number interface {
	int8 | int | int32 | int64 | float32 | float64
}

func generateNum[n number](min, max n) n {
	return min + (n(rand.Float64() * float64(max+1-min)))
}

func init() {
	generateSlice = basicGenerationFunction{
		_func: func(parameters ...any) any {

			field := parameters[0].(*GeneratedField)
			generationConfig := field.Generator.GenerationConfig
			sliceType := reflect.TypeOf(field.Value.Interface()).Elem()
			min := generationConfig.sliceMinLength
			max := generationConfig.sliceMaxLength

			len := generateNum(min, max)
			sliceOfElementType := reflect.SliceOf(sliceType)
			slice := reflect.MakeSlice(sliceOfElementType, 0, 1024)
			sliceElement := reflect.New(sliceType)

			for i := 0; i < len; i++ {
				newField := &GeneratedField{
					Name:      fmt.Sprintf("%s#%d", field.Name, i),
					Value:     reflect.ValueOf(sliceElement.Interface()).Elem(),
					Tag:       field.Tag,
					Generator: field.Generator.Copy(),
					Parent:    field,
				}

				newField.SetValue()

				slice = reflect.Append(slice, sliceElement.Elem())
			}
			return slice.Interface()

		},
	}

	generatePointerValue = basicGenerationFunction{

		_func: func(parameters ...any) any {
			field := parameters[0].(*GeneratedField)
			if !field.Generator.GenerationConfig.setNonRequiredFields {
				return nil
			}

			field.Value.Set(reflect.New(field.Value.Type().Elem()))
			fieldPointerValue := *field
			fieldPointerValue.Value = field.Value.Elem()
			fieldPointerValue.PointerValue = &field.Value
			fieldPointerValue.SetValue()

			if field.Value.Elem().CanSet() {
				field.Value.Elem().Set(fieldPointerValue.Value)
			}

			return field.Value.Interface()

		},
	}

	generateStruct = basicGenerationFunction{

		_func: func(parameters ...any) any {
			field := parameters[0].(*GeneratedField)
			field.setStructValues()
			return field.Value.Interface()
		},
	}
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

	generateStruct basicGenerationFunction

	generateSlice basicGenerationFunction

	generatePointerValue basicGenerationFunction

	generateDateTime basicGenerationFunction = basicGenerationFunction{
		_func: func(parameters ...any) any {
			return time.Now().UTC().Format(time.RFC3339)
		},
	}
)
