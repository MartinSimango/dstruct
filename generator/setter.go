package generator

import (
	"reflect"
)

func setValue(val reflect.Value, tags reflect.StructTag, generationConfig *GenerationConfig) {
	switch val.Kind() {
	case reflect.Struct:
		setStructValues(val, generationConfig)
	case reflect.Slice:
		panic("Unhandled setValue case")
	case reflect.Pointer:
		panic("Unhandled setValue case")
	case reflect.Interface:
		val.Set(reflect.Zero(val.Type()))
	default:
		val.Set(reflect.ValueOf(generationFunctionFromTags(val.Kind(), tags, generationConfig).Generate()))
	}
}

func setStructValues(config reflect.Value, generationConfig *GenerationConfig) {
	for j := 0; j < config.NumField(); j++ {
		setValue(config.Field(j), config.Type().Field(j).Tag, generationConfig)
	}

}
