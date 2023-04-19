package generator

import (
	"fmt"
	"reflect"
	"strconv"
)

type ParentType map[string]int

type GeneratedField struct {
	Name         string
	Value        reflect.Value
	Tag          reflect.StructTag
	Generator    *Generator
	Parent       *GeneratedField
	PointerValue *reflect.Value
}

func NewGeneratedField(fqn string,
	value reflect.Value,
	tag reflect.StructTag,
	generator *Generator) *GeneratedField {
	g := &GeneratedField{
		Name:      fqn,
		Value:     value,
		Tag:       tag,
		Generator: generator,
	}
	return g
}

func (field *GeneratedField) checkForRecursiveDefinition() bool {
	var depth uint = 0
	var matchedField *GeneratedField
	for parent := field.Parent; parent != nil; parent = parent.Parent {
		if parent.Value.Type() == field.Value.Type() {
			if !field.Generator.GenerationConfig.recursiveDefinition.Allow {
				panic(fmt.Sprintf("github.com/MartinSimango/dstruct/generator: recursive definition found for field `%s` of type %s", parent.Name, parent.Value.Type()))
			}
			depth++
			if depth == 1 {
				matchedField = parent
			}
		}
		if depth == (field.Generator.GenerationConfig.recursiveDefinition.Count + 1) {

			if matchedField.PointerValue != nil {
				matchedField.PointerValue.SetZero()
			} else {
				matchedField.Value.SetZero()
			}
			return true
		}
	}
	return false

}

func (field *GeneratedField) SetValue() {

	switch field.Value.Kind() {
	case reflect.Struct:
		if field.checkForRecursiveDefinition() {
			return
		}
		GenerateStructFunc(field).Generate()
	case reflect.Pointer:
		GeneratePointerValueFunc(field).Generate()
	case reflect.Slice:
		if field.checkForRecursiveDefinition() {
			return
		}
		field.Value.Set(reflect.ValueOf(GenerateSliceFunc(field).Generate()))
	case reflect.Interface:
		field.Value.Set(reflect.Zero(field.Value.Type()))
	default:
		field.Value.Set(reflect.ValueOf(field.getGenerationFunction().Generate()))
	}
}

func (field *GeneratedField) setStructValues() {
	for j := 0; j < field.Value.NumField(); j++ {
		structField := &GeneratedField{
			Name:      field.Name + "." + field.Value.Type().Field(j).Name,
			Value:     field.Value.Field(j),
			Tag:       field.Value.Type().Field(j).Tag,
			Generator: field.Generator.Copy(),
			Parent:    field,
		}
		structField.SetValue()
	}
}

func (field *GeneratedField) getGenerationFunction() GenerationFunction {
	kind := field.Value.Kind()
	tags := field.Tag

	switch kind {
	case reflect.Slice:
		return GenerateSliceFunc(field)
	case reflect.Struct:
		return GenerateStructFunc(field)
	case reflect.Ptr:
		return GeneratePointerValueFunc(field)
	}

	if field.Generator.GenerationConfig.valueGenerationType == UseDefaults {
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
	return field.Generator.DefaultGenerationFunctions[kind]
}
