package core

import (
	"reflect"

	"github.com/MartinSimango/dstruct/generator"
)

func GeneratePointerValueFunc(field *GeneratedField) generator.GenerationFunction {
	return &coreGenerationFunction{
		_func: func(parameters ...any) any {
			field := parameters[0].(*GeneratedField)
			if !field.Config.GenerationSettings.SetNonRequiredFields {
				return nil
			}
			// t := field.Value.Interface()
			// field.Value.Set(reflect.New(field.Value.Type().Elem()))
			fieldValueCopy := reflect.New(field.Value.Type()).Elem()
			fieldValueCopy.Set(reflect.New(field.Value.Type().Elem()))
			fieldPointerValue := *field
			fieldPointerValue.Value = fieldValueCopy.Elem()
			fieldPointerValue.PointerValue = &fieldValueCopy
			fieldPointerValue.SetValue()

			if fieldValueCopy.Elem().CanSet() {
				fieldValueCopy.Elem().Set(fieldPointerValue.Value)
			}
			return fieldValueCopy.Interface()

			// field.Value.Set(reflect.New(field.Value.Type().Elem()))
			// fieldPointerValue := *field
			// fieldPointerValue.Value = field.Value.Elem()
			// fieldPointerValue.PointerValue = &field.Value
			// fieldPointerValue.SetValue()

			// if field.Value.Elem().CanSet() {
			// 	field.Value.Elem().Set(fieldPointerValue.Value)
			// }
			//
			//
			// return field.Value.Interface()
		},
		kind: reflect.Ptr,
		args: []any{field},
	}
}
