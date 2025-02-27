package core

import (
	"fmt"
	"reflect"

	"github.com/MartinSimango/dstruct/dreflect"
	"github.com/MartinSimango/dstruct/generator"
	"github.com/MartinSimango/dstruct/generator/config"
)

func GenerateSliceFunc(
	field *GeneratedField,
	config config.Config,
	generationFunctions DefaultGenerationFunctions,
) generator.GenerationFunction {
	return &coreGenerationFunction{
		_func: func(parameters ...any) any {
			field := parameters[0].(*GeneratedField)
			sliceConfig := field.Config.GenerationFunctions[reflect.Slice].GetConfig().Slice()
			sliceType := reflect.TypeOf(field.Value.Interface()).Elem()
			min := sliceConfig.MinLength()
			max := sliceConfig.MaxLength()

			len := generateNum(min, max)
			sliceOfElementType := reflect.SliceOf(sliceType)
			slice := reflect.MakeSlice(sliceOfElementType, 0, 1024)
			sliceElement := reflect.New(sliceType)

			for i := 0; i < len; i++ {
				elemValue := reflect.ValueOf(sliceElement.Interface()).Elem()
				fieldConfig := NewGenerateFieldConfig(
					field.Config.GenerationFunctions[reflect.Slice].GetConfig(),
					field.Config.GenerationSettings,
				)
				fieldConfig.GenerationFunctions = generationFunctions
				newField := &GeneratedField{
					Name:        fmt.Sprintf("%s#%d", field.Name, i),
					Value:       elemValue,
					Tag:         field.Tag,
					Config:      fieldConfig,
					Parent:      field,
					customTypes: field.customTypes,
					typeHash:    dreflect.GetTypeHash(elemValue.Interface()),
				}

				newField.SetValue()

				slice = reflect.Append(slice, sliceElement.Elem())
			}
			return slice.Interface()
		},
		args: []any{field},
		kind: reflect.Slice,
	}
}
