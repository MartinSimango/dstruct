package dstruct

import (
	"fmt"
	"reflect"
	"time"

	"github.com/MartinSimango/dstruct/dreflect"
	"github.com/MartinSimango/dstruct/generator"
	"github.com/MartinSimango/dstruct/generator/config"
	"github.com/MartinSimango/dstruct/generator/core"
)

// CustomType allows types with the same type of 'Value' to have a specific generation function that
// is stored by the 'FunctionHolder' field. This is useful for types that are not supported by the
// default generation functions.
type CustomType struct {
	Value          any
	FunctionHolder core.FunctionHolder
}

// DStructGeneratedStruct implements GeneratedStruct.
type DStructGeneratedStruct[T any] struct {
	*DynamicStructModifierImpl
	fieldContexts GeneratedFieldContexts
	structConfig  core.GeneratedFieldConfig
	instance      T
	customTypes   map[reflect.Type]core.FunctionHolder
}

var _ GeneratedStruct = &DStructGeneratedStruct[int]{}

func NewGeneratedStruct[T any](val T) *DStructGeneratedStruct[T] {
	return NewGeneratedStructWithConfig(
		val,
		config.NewDstructConfig(),
		config.DefaultGenerationSettings(),
	)
}

func NewGeneratedStructWithConfig[T any](val T,
	cfg config.Config,
	settings config.GenerationSettings,
	customTypes ...CustomType,
) *DStructGeneratedStruct[T] {
	generatedStruct := &DStructGeneratedStruct[T]{
		DynamicStructModifierImpl: ExtendStruct(val).Build().(*DynamicStructModifierImpl),
		fieldContexts:             make(GeneratedFieldContexts),
		structConfig: core.NewGenerateFieldConfig(
			cfg,
			settings,
		),
		customTypes: make(map[reflect.Type]core.FunctionHolder),
	}

	for _, v := range customTypes {
		generatedStruct.addCustomType(v)
	}
	generatedStruct.addCustomTypes()
	generatedStruct.populateGeneratedFields(generatedStruct.root)
	return generatedStruct
}

// InstanceT returns the instance of the generated struct.
func (gs *DStructGeneratedStruct[T]) InstanceT() T {
	return gs.instance
}

// NewT returns a new instance of the generated struct
func (gs *DStructGeneratedStruct[T]) NewT() *T {
	gs.DynamicStructModifierImpl.New()
	return &gs.instance
}

// Generate implements GeneratedStruct.Generate.
func (gs *DStructGeneratedStruct[T]) Generate() {
	gs.generateFields()

	switch any(*new(T)).(type) {
	case nil:
		gs.instance = gs.DynamicStructModifierImpl.Instance().(T)
		return
	}

	gs.instance = toType[T](gs.DynamicStructModifierImpl)

	// TODO: should only be called if new struct fields are generated
	// an idea is to have generatedFields return a bool to indicate if new fields were generated
	gs.Update()
}

// SetFieldGenerationSettings implements GeneratedStruct.SetFieldGenerationSettings
func (gs *DStructGeneratedStruct[T]) SetFieldGenerationSettings(
	field string,
	settings config.GenerationSettings,
) error {
	if err := gs.assertFieldExists(field); err != nil {
		return err
	}

	if gs.fieldContexts[field] == nil {
		gs.propagateSettings(gs.fieldNodeMap[field], settings)
	} else {
		gs.fieldContexts[field].GeneratedField.SetGenerationSettings(settings)
	}

	return nil
}

// GetFieldGenerationSettings implements GeneratedStruct.GetFieldGenerationSettings
func (gs *DStructGeneratedStruct[T]) GetFieldGenerationSettings(
	field string,
) (config.GenerationSettings, error) {
	if err := gs.assertFieldExists(field); err != nil {
		return config.GenerationSettings{}, err
	}

	if gs.fieldContexts[field] == nil {
		return config.GenerationSettings{}, fmt.Errorf(
			"field %s does not have any generation settings",
			field,
		)
	}

	return gs.fieldContexts[field].GeneratedField.Config.GenerationSettings, nil
}

// SetGenerationSettings implements GeneratedStruct.SetGenerationSettings
func (gs *DStructGeneratedStruct[T]) SetGenerationSettings(
	settings config.GenerationSettings,
) {
	gs.structConfig.GenerationSettings = settings

	gs.propagateSettings(gs.root, gs.structConfig.GenerationSettings)
}

// GetGenerationSettings implements GeneratedStruct.GetGenerationSettings
func (gs *DStructGeneratedStruct[T]) GetGenerationSettings() config.GenerationSettings {
	return gs.structConfig.GenerationSettings
}

// SetFieldGenerationConfig implements GeneratedStruct.SetFieldGenerationConfig
func (gs *DStructGeneratedStruct[T]) SetFieldGenerationConfig(
	field string,
	cfg config.Config,
) error {
	if err := gs.assertFieldExists(field); err != nil {
		return err
	}

	if gs.fieldContexts[field] == nil {
		gs.propagateConfig(gs.fieldNodeMap[field], cfg)
	} else {
		// this will be fields that are either custom types or fields that have no children
		gs.fieldContexts[field].GeneratedField.SetConfig(cfg)
	}

	return nil
}

func (gs *DStructGeneratedStruct[T]) GetFieldGenerationConfig(field string) (config.Config, error) {
	if err := gs.assertFieldExists(field); err != nil {
		return nil, err
	}

	if gs.fieldContexts[field] == nil {
		return nil, fmt.Errorf(
			"field %s does not have a generation config",
			field,
		)
	}

	return gs.fieldContexts[field].GeneratedField.Config.GenerationConfig, nil
}

// SetGenerationConfig implements GeneratedStruct.SetGenerationConfig
func (gs *DStructGeneratedStruct[T]) SetGenerationConfig(cfg config.Config) {
	gs.structConfig.GenerationConfig = cfg

	gs.propagateConfig(gs.root, gs.structConfig.GenerationConfig)
}

// GetGenerationConfig implements GeneratedStruct.GetGenerationConfig
func (gs *DStructGeneratedStruct[T]) GetGenerationConfig() config.Config {
	return gs.structConfig.GenerationConfig
}

// SetFieldGenerationFunction implements GeneratedStruct.SetFieldGenerationFunction
func (gs *DStructGeneratedStruct[T]) SetFieldGenerationFunction(field string,
	functionHolder core.FunctionHolder,
) error {
	if err := gs.assertFieldExists(field); err != nil {
		return err
	}

	if gs.fieldContexts[field] == nil {
		return fmt.Errorf("field %s does not have a generation function", field)
	}

	gs.fieldContexts[field].GeneratedField.SetGenerationFunction(functionHolder)

	return nil
}

// // GetFieldGenerationFunction implements GeneratedStruct.GetFieldGenerationFunction
// func (gs *DStructGeneratedStruct[T]) GetFieldGenerationFunction(
// 	field string,
// ) (core.FunctionHolder, error) {
// 	if err := gs.assertFieldExists(field); err != nil {
// 		return nil, err
// 	}
//
// 	if gs.fieldContexts[field] == nil {
// 		return nil, fmt.Errorf("field %s does not have a generation function", field)
// 	}
//
// 	return gs.fieldContexts[field].GeneratedgField. , nil
// }
//
// GetFieldGenerationFunction_ implements GeneratedStruct.GetFieldGenerationFunction_

//	func (gs *DStructGeneratedStruct[T]) GetFieldGenerationFunction_(field string) core.FunctionHolder {
//		if err := gs.assertFieldExists(field); err != nil {
//			panic(err)
//		}
//
//		if gs.fieldContexts[field] == nil {
//			panic(fmt.Errorf("field %s does not have a generation function", field))
//		}
//
//		return gs.fieldContexts[field].GeneratedField.GenerationFunction
//	}
//
// SetFieldGenerationFunctions implements GeneratedStruct.SetFieldGenerationFunctions
func (gs *DStructGeneratedStruct[T]) SetFieldGenerationFunctions(
	field string,
	functions core.DefaultGenerationFunctions,
) error {
	if err := gs.assertFieldExists(field); err != nil {
		return err
	}

	if gs.fieldContexts[field] == nil {
		gs.propagateGenerationFunctions(gs.fieldNodeMap[field], functions)
	} else {
		gs.fieldContexts[field].GeneratedField.SetGenerationFunctions(functions)
	}

	return nil
}

// SetGenerationFunctions implements GeneratedStruct.SetGenerationFunctions
func (gs *DStructGeneratedStruct[T]) SetGenerationFunctions(
	functions core.DefaultGenerationFunctions,
) {
	gs.structConfig.GenerationFunctions = functions
	gs.propagateGenerationFunctions(gs.root, functions)
}

// SetFieldFromTask implements GeneratedStruct.SetFieldFromTask
func (gs *DStructGeneratedStruct[T]) SetFieldFromTask(
	field string,
	task generator.Task,
	params ...any,
) error {
	taskProperties, err := generator.CreateTaskProperties(
		field,
		generator.GetTagForTask(generator.TaskName(task.Name()), params...),
	)
	if err != nil {
		return err
	}

	gs.SetFieldGenerationFunction(
		field,
		core.NewFunctionHolderNoArgs(task.GenerationFunction(*taskProperties)),
	)
	return nil
}

// GetFieldGenerationConfig_ implements GeneratedStruct.
func (gs *DStructGeneratedStruct[T]) GetFieldGenerationConfig_(field string) config.Config {
	if err := gs.assertFieldExists(field); err != nil {
		panic(err)
	}
	if gs.fieldContexts[field] == nil {
		panic(fmt.Errorf("field %s does not have a generation config", field))
	}
	return gs.fieldContexts[field].GeneratedField.Config.GenerationConfig
}

// GetFieldGenerationSettings_ implements GeneratedStruct.
func (gs *DStructGeneratedStruct[T]) GetFieldGenerationSettings_(
	field string,
) config.GenerationSettings {
	if err := gs.assertFieldExists(field); err != nil {
		panic(err)
	}
	if gs.fieldContexts[field] == nil {
		panic(fmt.Errorf("field %s does not have any generation settings", field))
	}
	return gs.fieldContexts[field].GeneratedField.Config.GenerationSettings
}

func toType[T any](gs DynamicStructModifier) T {
	return dreflect.ConvertToType[T](gs.Instance())
}

func toPointerType[T any](gs DynamicStructModifier) *T {
	return dreflect.ConvertToType[*T](gs.New())
}

func (gs *DStructGeneratedStruct[T]) generateFields() {
	for k, genFunc := range gs.fieldContexts {
		if err := gs.Set(k, genFunc.Generate()); err != nil {
			fmt.Println(err)
		}
	}
}

func (gs *DStructGeneratedStruct[T]) addCustomTypes() {
	gs.addCustomType(
		CustomType{
			time.Time{},
			core.DefaultDateFunctionHolder(gs.structConfig.GenerationConfig.Date()),
		},
	)
}

func (gs *DStructGeneratedStruct[T]) addCustomType(customType CustomType) {
	// TODO: restrict some types from being added such as nil, ints etc
	gs.customTypes[reflect.TypeOf(customType.Value)] = customType.FunctionHolder
	// the function holder kind is the find that the function retrusn which could either be an existing kind or a new kind i.ie time.Time would be a new kind
	// TODO:idea: GenerationsFunction key should be type of CustomKind and not reflect.Kind
	gs.structConfig.GenerationFunctions[customType.FunctionHolder.Kind()] = customType.FunctionHolder
}

func (gs *DStructGeneratedStruct[T]) createGeneratedField(
	field *Node[structField],
	kind reflect.Kind,
) *core.GeneratedField {
	return core.NewGeneratedField(field.data.fullyQualifiedName,
		field.data.value,
		field.data.tag,
		gs.structConfig.Copy(kind),
		gs.customTypes,
		field.data.goType)
}

func (gs *DStructGeneratedStruct[T]) populateGeneratedFields(node *Node[structField]) {
	for _, field := range node.children {
		if customType := gs.customTypes[field.data.goType]; customType != nil {
			gs.fieldContexts[field.data.fullyQualifiedName] = core.NewGeneratedFieldContext(
				gs.createGeneratedField(field, customType.Kind()),
			)
		} else if field.HasChildren() {
			gs.populateGeneratedFields(field)
		} else {
			gs.fieldContexts[field.data.fullyQualifiedName] = core.NewGeneratedFieldContext(
				gs.createGeneratedField(field, field.data.value.Kind()),
			)
		}
	}
}

func (gs *DStructGeneratedStruct[T]) propagateConfig(
	node *Node[structField],
	cfg config.Config,
) {
	for _, field := range node.children {
		// Don't propagate changes to children nodes if the field is a custom type
		if field.HasChildren() &&
			!gs.fieldContexts[field.data.fullyQualifiedName].GeneratedField.IsCustomType() {
			gs.propagateConfig(field, cfg)
		} else {
			gs.fieldContexts[field.data.fullyQualifiedName].GeneratedField.SetConfig(cfg)
		}
	}
}

func (gs *DStructGeneratedStruct[T]) propagateSettings(
	node *Node[structField],
	settings config.GenerationSettings,
) {
	for _, field := range node.children {
		if field.HasChildren() &&
			!gs.fieldContexts[field.data.fullyQualifiedName].GeneratedField.IsCustomType() {
			gs.propagateSettings(field, settings)
		} else {
			gs.fieldContexts[field.data.fullyQualifiedName].GeneratedField.SetGenerationSettings(settings)
		}
	}
}

func (gs *DStructGeneratedStruct[T]) propagateGenerationFunctions(
	node *Node[structField],
	functions core.DefaultGenerationFunctions,
) {
	for _, field := range node.children {
		if field.HasChildren() &&
			!gs.fieldContexts[field.data.fullyQualifiedName].GeneratedField.IsCustomType() {
			gs.propagateGenerationFunctions(field, functions)
		} else {
			gs.fieldContexts[field.data.fullyQualifiedName].GeneratedField.SetGenerationFunctions(functions)
		}
	}
}

func (gs *DStructGeneratedStruct[T]) propagateGenerationFunction(
	node *Node[structField],
	functionHolder core.FunctionHolder,
) {
	for _, field := range node.children {
		if field.HasChildren() &&
			!gs.fieldContexts[field.data.fullyQualifiedName].GeneratedField.IsCustomType() {
			gs.propagateGenerationFunction(field, functionHolder)
		} else {
			gs.fieldContexts[field.data.fullyQualifiedName].GeneratedField.SetGenerationFunction(functionHolder)
		}
	}
}

func (gs *DStructGeneratedStruct[T]) assertFieldExists(field string) error {
	if gs.fieldNodeMap[field] == nil {
		return fmt.Errorf("field %s does not exist within the struct", field)
	}
	return nil
}