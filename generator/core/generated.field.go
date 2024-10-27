package core

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/MartinSimango/dstruct/generator"
	"github.com/MartinSimango/dstruct/generator/config"
)

type GeneratedFieldConfig struct {
	GenerationFunctions DefaultGenerationFunctions
	GenerationSettings  config.GenerationSettings
	GenerationConfig    config.Config
}

func (gf *GeneratedFieldConfig) Copy(kind reflect.Kind) GeneratedFieldConfig {
	cfg := gf.GenerationConfig.Copy()
	return GeneratedFieldConfig{
		GenerationFunctions: gf.GenerationFunctions.Copy(cfg, kind),
		GenerationSettings:  gf.GenerationSettings,
		GenerationConfig:    cfg,
	}
}

func (gf *GeneratedFieldConfig) SetConfig(cfg config.Config) {
	gf.GenerationConfig = cfg
	for _, v := range gf.GenerationFunctions {
		v.SetConfig(cfg)
	}
}

func NewGenerateFieldConfig(
	cfg config.Config,
	settings config.GenerationSettings,
) GeneratedFieldConfig {
	return GeneratedFieldConfig{
		GenerationFunctions: NewDefaultGenerationFunctions(cfg),
		GenerationSettings:  settings,
		GenerationConfig:    cfg,
	}
}

type GeneratedField struct {
	Name                      string
	Value                     reflect.Value
	Tag                       reflect.StructTag
	Config                    GeneratedFieldConfig
	Parent                    *GeneratedField
	PointerValue              *reflect.Value
	customTypes               map[string]FunctionHolder
	goType                    string
	currentGenerationFunction generator.GenerationFunction
}

func NewGeneratedField(fqn string,
	value reflect.Value,
	tag reflect.StructTag,
	config GeneratedFieldConfig,
	customTypes map[string]FunctionHolder,
	goType string,
) *GeneratedField {
	generateField := &GeneratedField{
		Name:        fqn,
		Value:       value,
		Tag:         tag,
		Config:      config,
		customTypes: customTypes,
		goType:      goType,
	}
	// TODO: add custom type to GenerationFunctions
	if value.Kind() == reflect.Slice {
		config.GenerationFunctions[reflect.Slice] = NewSliceFunctionHolder(
			GenerateSliceFunc,
			generateField,
			config.GenerationConfig,
			generateField.Config.GenerationFunctions,
		)
	}
	if value.Kind() == reflect.Ptr {
		config.GenerationFunctions[reflect.Ptr] = NewPointerFunctionHolder(
			GeneratePointerValueFunc,
			generateField,
		)
	}

	// if field.IsCustomType() {
	// 	config.GenerationFunctions[field.customTypeFunctionHolder().GetFunction().Kind()] = field.customTypeFunctionHolder()
	// }

	return generateField
}

func (field *GeneratedField) IsCustomType() bool {
	return field.customTypes[field.goType] != nil
}

func (field *GeneratedField) customTypeFunctionHolder() FunctionHolder {
	return field.customTypes[field.goType]
}

func (field *GeneratedField) SetConfig(cfg config.Config) {
	field.Config.GenerationConfig.SetFrom(cfg)

	// kind := field.Value.Kind()
	// if field.IsCustomType() {
	// 	field.customTypeFunctionHolder().SetConfig(cfg)
	// } else if field.Config.GenerationFunctions[kind] != nil {
	// 	fmt.Println("Setting config for kind: ", kind)
	// 	field.Config.GenerationFunctions[kind].SetConfig(cfg)
	// } else {
	// 	field.Config.SetConfig(cfg)
	// }
}

func (field *GeneratedField) SetGenerationSettings(settings config.GenerationSettings) {
	field.Config.GenerationSettings = settings
}

func (field *GeneratedField) SetGenerationFunction(
	functionHolder FunctionHolder,
) {
	if field.IsCustomType() {
		field.customTypes[field.goType] = functionHolder
	} else if field.Config.GenerationFunctions[field.Value.Kind()] != nil {
		field.Tag = reflect.StructTag("") // remove tags to ensure the field is generated with the new function
		field.Config.GenerationFunctions[field.Value.Kind()] = functionHolder
	}
}

func (field *GeneratedField) SetGenerationFunctions(functions DefaultGenerationFunctions) {
	field.Config.GenerationFunctions = functions
}

func (field *GeneratedField) checkForRecursiveDefinition(fail bool) bool {
	var depth uint = 0
	var matchedField *GeneratedField
	for parent := field.Parent; parent != nil; parent = parent.Parent {
		if parent.Value.Type() == field.Value.Type() {
			if !field.Config.GenerationSettings.RecursiveDefinition.Allow || fail {
				panic(
					fmt.Sprintf(
						"github.com/MartinSimango/dstruct/generator: recursive definition found for field `%s` of type %s",
						parent.Name,
						parent.Value.Type(),
					),
				)
			}
			depth++
			if depth == 1 {
				matchedField = parent
			}
		}
		if depth == (field.Config.GenerationSettings.RecursiveDefinition.Depth + 1) {
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
	// check if the current field is a custom type with it's own generation function
	if customType := field.customTypes[field.goType]; customType != nil {
		field.Value.Set(reflect.ValueOf(customType.GetFunction().Generate()))
		return false
	}
	kind := field.Value.Kind()
	switch kind {
	case reflect.Struct:
		if field.checkForRecursiveDefinition(false) {
			return true
		}
		GenerateStructFunc(field).Generate()
	case reflect.Pointer:
		field.Value.Set(reflect.ValueOf(GeneratePointerValueFunc(field).Generate()))
	case reflect.Slice:
		// Don't allow recursion within slices
		if field.checkForRecursiveDefinition(true) {
			return true
		}
		field.Value.Set(
			reflect.ValueOf(field.Config.GenerationFunctions[kind].GetFunction().Generate()),
		)
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
			Name:        field.Name + "." + field.Value.Type().Field(j).Name,
			Value:       field.Value.Field(j),
			Tag:         field.Value.Type().Field(j).Tag,
			Config:      field.Config.Copy(field.Value.Field(j).Kind()),
			Parent:      field,
			customTypes: field.customTypes,
			goType:      field.Value.Field(j).Type().Name(),
		}
		structField.SetValue()
	}
}

func (field *GeneratedField) getGenerationFunction() generator.GenerationFunction {
	//  keep cache of the current generation function
	if field.currentGenerationFunction != nil {
		// TODO: implement cache
		return field.currentGenerationFunction
	}

	// check if field is a custom type with it's own generation function
	if field.customTypes[field.goType] != nil {
		return field.customTypes[field.goType].GetFunction()
	}

	kind := field.Value.Kind()
	tags := field.Tag

	switch kind {
	case reflect.Struct:
		return GenerateStructFunc(field)
	case reflect.Ptr:
		return GeneratePointerValueFunc(field)
	}

	if field.Config.GenerationSettings.ValueGenerationType == config.UseDefaults {
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
	// TODO : replace or remove this
	case "date-time":
		return GenerateDateTimeFunc()
	}

	enum, ok := tags.Lookup("enum")
	if ok {
		numEnums, _ := strconv.Atoi(enum)
		return GenerateFixedValueFunc(
			tags.Get(fmt.Sprintf("enum_%d", generateNum(0, numEnums-1)+1)),
		)
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

	// if we get no match, we default to the default generation function for the kind
	// kidds of type Slice will be handled here as their default generation function for a slice will be overwritten when the generated field is created.
	return field.Config.GenerationFunctions[kind].GetFunction()
}
