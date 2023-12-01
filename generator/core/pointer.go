package core

import (
	"reflect"

	"github.com/MartinSimango/dstruct/generator"
)

func GeneratePointerValueFunc(field *GeneratedField) generator.GenerationFunction {

	return &coreGenerationFunction{
		_func: func(parameters ...any) any {
			field := parameters[0].(*GeneratedField)
			if !field.GenerationValueConfig.SetNonRequiredFields {
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
		args: []any{field},
	}

}
