package core

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/MartinSimango/dstruct/generator"
	"github.com/MartinSimango/dstruct/generator/config"
)

type GeneratedFieldConfig struct {
	GenerationFunctions   DefaultGenerationFunctions
	GenerationValueConfig config.GenerationValueConfig
}

func (gf *GeneratedFieldConfig) Copy(kind reflect.Kind) (gfc GeneratedFieldConfig) {
	return GeneratedFieldConfig{
		GenerationFunctions:   gf.GenerationFunctions.Copy(kind),
		GenerationValueConfig: gf.GenerationValueConfig,
	}
}

func (gf *GeneratedFieldConfig) SetConfig(cfg config.Config) {
	for _, v := range gf.GenerationFunctions {
		v.SetConfig(cfg)
	}
}

func NewGenerateFieldConfig(cfg config.Config, gvc config.GenerationValueConfig) GeneratedFieldConfig {
	return GeneratedFieldConfig{
		GenerationFunctions:   NewDefaultGenerationFunctions(cfg),
		GenerationValueConfig: gvc,
	}
}

type GeneratedField struct {
	Name  string
	Value reflect.Value
	Tag   reflect.StructTag
	GeneratedFieldConfig
	Parent       *GeneratedField
	PointerValue *reflect.Value
}

func NewGeneratedField(fqn string,
	value reflect.Value,
	tag reflect.StructTag,
	generatedFieldConfig GeneratedFieldConfig) *GeneratedField {
	g := &GeneratedField{
		Name:                 fqn,
		Value:                value,
		Tag:                  tag,
		GeneratedFieldConfig: generatedFieldConfig,
	}
	return g
}

func (field *GeneratedField) checkForRecursiveDefinition(fail bool) bool {
	var depth uint = 0
	var matchedField *GeneratedField
	for parent := field.Parent; parent != nil; parent = parent.Parent {
		if parent.Value.Type() == field.Value.Type() {
			if !field.GenerationValueConfig.RecursiveDefinition.Allow || fail {
				panic(fmt.Sprintf("github.com/MartinSimango/dstruct/generator: recursive definition found for field `%s` of type %s", parent.Name, parent.Value.Type()))
			}
			depth++
			if depth == 1 {
				matchedField = parent
			}
		}
		if depth == (field.GenerationValueConfig.RecursiveDefinition.Depth + 1) {
			// fmt.Println(":DF ", field.Name, matchedField.Name, field.Value.Type(), matchedField.Value.Type(), matchedField.Value, depth)
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

func (field *GeneratedField) SetValue() bool {
	kind := field.Value.Kind()
	switch kind {
	case reflect.Struct:
		if field.checkForRecursiveDefinition(false) {
			return true
		}
		GenerateStructFunc(field).Generate()
	case reflect.Pointer:
		GeneratePointerValueFunc(field).Generate()
	case reflect.Slice:
		// Don't allow recursion within slices
		if field.checkForRecursiveDefinition(true) {
			return true
		}
		field.Value.Set(reflect.ValueOf(field.GenerationFunctions[kind].GetFunction().Generate()))
	case reflect.Interface:
		field.Value.Set(reflect.Zero(field.Value.Type()))
	default:
		field.Value.Set(reflect.ValueOf(field.getGenerationFunction().Generate()))
	}
	return false
}

func (field *GeneratedField) setStructValues() {
	for j := 0; j < field.Value.NumField(); j++ {
		structField := &GeneratedField{
			Name:                 field.Name + "." + field.Value.Type().Field(j).Name,
			Value:                field.Value.Field(j),
			Tag:                  field.Value.Type().Field(j).Tag,
			GeneratedFieldConfig: field.GeneratedFieldConfig.Copy(field.Value.Field(j).Kind()),
			Parent:               field,
		}
		structField.SetValue()
	}
}

func (field *GeneratedField) getGenerationFunction() generator.GenerationFunction {
	kind := field.Value.Kind()
	tags := field.Tag

	// TODO see if needed
	switch kind {
	case reflect.Slice:
		return GenerateSliceFunc(field, nil)
	case reflect.Struct:
		return GenerateStructFunc(field)
	case reflect.Ptr:
		return GeneratePointerValueFunc(field)
	}
	if field.GenerationValueConfig.ValueGenerationType == config.UseDefaults {
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

	_, ok = tags.Lookup("gen_task")
	if ok {
		taskProperties, err := generator.CreateTaskProperties(field.Name, tags)
		if err != nil {
			panic(fmt.Sprintf("Error decoding gen_task: %s", err.Error()))
		}
		task := generator.GetTask(taskProperties.TaskName)
		if task == nil {
			panic(fmt.Sprintf("Unregistered task %s", taskProperties.TaskName))
		}
		return task.GenerationFunction(*taskProperties)
	}
	return field.GenerationFunctions[kind].GetFunction()
}
