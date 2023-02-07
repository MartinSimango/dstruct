package generator

import (
	"reflect"
)

func setValue(field Field, generationConfig *GenerationConfig) {
	switch field.Value.Kind() {
	case reflect.Struct:
		setStructValues(field, generationConfig)
	case reflect.Slice:
		panic("Unhandled setValue case slice ")
	case reflect.Pointer:
		panic("Unhandled setValue case ptr")
	case reflect.Interface:
		field.Value.Set(reflect.Zero(field.Value.Type()))
	default:
		field.Value.Set(reflect.ValueOf(generationFunctionFromTags(field, generationConfig).Generate()))
	}
}

func setStructValues(field Field, generationConfig *GenerationConfig) {
	for j := 0; j < field.Value.NumField(); j++ {
		structField := Field{
			Name:  field.Value.Type().Field(j).Name,
			Value: field.Value.Field(j),
			Tag:   field.Value.Type().Field(j).Tag,
		}
		if field.Name != "" {
			field.Name = field.Name + "." + structField.Name
		}
		setValue(structField, generationConfig)
	}

}
