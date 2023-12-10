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

type GenerationFields map[string]*core.GenerationUnit

type GeneratedStruct interface {
	DynamicStructModifier
	// Generate generates fields for the struct
	Generate()

	// GenerateAndUpdate Generates fields and updates the root tree for the underlying struct. Allowing
	// new generated fields to be accessed and modified by Set and Get methods.
	GenerateAndUpdate()

	SetFieldGenerationValueConfig(field string, config config.GenerationValueConfig)

	// GetFieldValueGenerationConfig gets the generation config for field within the struct.
	GetFieldValueGenerationConfig(field string) config.GenerationValueConfig

	GetValueGenerationConfig() config.GenerationValueConfig

	// SetFieldGenerationConfig sets the generation config for field within the struct. It returns
	// an error if the field does not exist or if the field cannot be generated.
	// Fields that can be generated are struct fields of the most basic type i.e a struct fields
	// that are structs cannot be generated, however it's fields can be.
	//
	// Fields types that cannot be generated: structs, func, chan, any (will default to a nil value being generated).
	//
	// Note: Pointers to structs can be generated.
	SetFieldConfig(field string, generationConfig config.Config) error

	GetFieldConfig(field string) config.Config

	SetConfig(config config.Config)

	GetConfig() config.Config

	SetGenerationValueConfig(config config.GenerationValueConfig)

	SetFieldFromTask(field string, task generator.Task, params ...any) error
}

type GeneratedStructImpl[T any] struct {
	*DynamicStructModifierImpl
	generatedFields      GenerationFields
	generatedFieldConfig core.GeneratedFieldConfig
	config               config.Config
	instance             T
	excludedTypes        map[reflect.Type]bool
}

// var _ GeneratedStruct = &GeneratedStructImpl[int]{}

func NewGeneratedStruct[T any](val T) *GeneratedStructImpl[T] {
	return NewGeneratedStructWithConfig(val, config.NewConfig())
}

func NewGeneratedStructWithConfig[T any](val T,
	cfg config.Config) *GeneratedStructImpl[T] {
	generatedStruct := &GeneratedStructImpl[T]{
		DynamicStructModifierImpl: ExtendStruct(val).Build().(*DynamicStructModifierImpl),
		generatedFields:           make(GenerationFields),
		config:                    cfg,
		generatedFieldConfig:      core.NewGenerateFieldConfig(cfg, config.DefaultGenerationValueConfig()),
		excludedTypes:             make(map[reflect.Type]bool),
	}
	generatedStruct.ExcludeType(time.Time{}, core.DefaultDateFunctionHolder(cfg.Date()))
	generatedStruct.populateGeneratedFields(generatedStruct.root)
	return generatedStruct
}

func (gs *GeneratedStructImpl[T]) ExcludeType(val any, function core.FunctionHolder) {
	gs.excludedTypes[reflect.TypeOf(val)] = true
}

func (gs *GeneratedStructImpl[T]) populateGeneratedFields(node *Node[structField]) {

	for _, field := range node.children {
		if field.HasChildren() {
			if gs.excludedTypes[field.data.goType] {
				field.data.value.Set(reflect.ValueOf(time.Now()))
				continue
			}
			gs.populateGeneratedFields(field)
			continue
		}
		fieldKind := field.data.value.Kind()

		generatedField := core.NewGeneratedField(field.data.fqn,
			field.data.value,
			field.data.tag,
			gs.generatedFieldConfig.Copy(fieldKind),
			gs.config.Copy(), // TODO account for nil
		)

		gs.generatedFields[field.data.fqn] = core.NewGenerationUnit(generatedField)

	}

}

func (gs *GeneratedStructImpl[T]) Generate() {
	gs.generateFields()

	v := any(*new(T))
	switch v.(type) {
	case nil:
		gs.instance = gs.DynamicStructModifierImpl.Instance().(T)
		return
	}
	gs.instance = ToType[T](gs.DynamicStructModifierImpl)
}

func (gs *GeneratedStructImpl[T]) GenerateAndUpdate() {
	gs.Generate()
	gs.Update()
}

func (gs *GeneratedStructImpl[T]) changeChildrenConfig(node *Node[structField], cfg config.Config) {
	for _, v := range node.children {
		if v.HasChildren() {
			gs.changeChildrenConfig(v, cfg)
			continue
		}
		gs.generatedFields[v.data.fqn].GenerationFunctions[v.data.typ.Kind()].SetConfig(cfg)
	}

}

func (gs *GeneratedStructImpl[T]) SetFieldGenerationConfig(field string, cfg config.Config) error {
	if gs.fieldNodeMap[field] == nil {
		return fmt.Errorf("field %s does not exist within the struct", field)
	}
	// TODO TEST for structs
	if gs.fieldNodeMap[field].HasChildren() {
		gs.changeChildrenConfig(gs.fieldNodeMap[field], cfg)
	}

	if gs.generatedFields[field] == nil {
		return fmt.Errorf("cannot set config for field %s", field)
	}
	kind := gs.generatedFields[field].Value.Kind()

	gs.generatedFields[field].GeneratedField.GenerationFunctions[kind].SetConfig(cfg)
	return nil
}

func (gs *GeneratedStructImpl[T]) SetFieldGenerationFunction(field string,
	functionHolder core.FunctionHolder) {

	// kind := gs.fieldNodeMap[field].data.GetType().Kind()
	// _ = functionHolder.(*core.FunctionHolderWithNoArgs)
	gs.generatedFields[field].GeneratedField.GenerationFunctions[functionHolder.Kind()] = functionHolder
}

func (gs *GeneratedStructImpl[T]) SetFieldGenerationFunctions(field string,
	defaultGenerationFunctions core.DefaultGenerationFunctions) {
	gs.generatedFields[field].GenerationFunctions = defaultGenerationFunctions
}

func (gs *GeneratedStructImpl[T]) SetGenerationConfig(config config.Config) {
	for name, field := range gs.fieldNodeMap {
		if field.HasChildren() {
			continue
		}
		gs.generatedFields[name].SetConfig(config.Copy())
	}
}

func (gs *GeneratedStructImpl[T]) SetFieldGenerationValueConfig(field string, config config.GenerationValueConfig) {
	gs.generatedFields[field].GenerationValueConfig = config
}

func (gs *GeneratedStructImpl[T]) SetGenerationValueConfig(config config.GenerationValueConfig) {
	for name, field := range gs.fieldNodeMap {
		if field.HasChildren() {
			continue
		}
		gs.generatedFields[name].GenerationValueConfig = config
	}
}

// GetFieldValueGenerationConfig implements GeneratedStruct.
func (gs *GeneratedStructImpl[T]) GetFieldValueGenerationConfig(field string) config.GenerationValueConfig {
	return gs.generatedFields[field].GenerationValueConfig
}

// GetFieldValueGenerationConfig implements GeneratedStruct.
func (gs *GeneratedStructImpl[T]) GetFieldGenerationConfig(field string) config.Config {
	k := gs.generatedFields[field].GeneratedField.Value.Kind()
	return gs.generatedFields[field].GeneratedField.GenerationFunctions[k].GetConfig()
}

// GetValueGenerationConfig implements GeneratedStruct.
func (gs *GeneratedStructImpl[T]) GetValueGenerationConfig() config.GenerationValueConfig {
	return gs.generatedFieldConfig.GenerationValueConfig
}

func (gs *GeneratedStructImpl[T]) SetFieldFromTask(field string, task generator.Task, params ...any) error {
	taskProperties, err := generator.CreateTaskProperties(field, generator.GetTagForTask(generator.TaskName(task.Name()), params...))
	if err != nil {
		return err
	}

	gs.SetFieldGenerationFunction(field, core.NewFunctionHolderNoArgs(task.GenerationFunction(*taskProperties)))
	return nil

}

func ToType[T any](gs DynamicStructModifier) T {
	return dreflect.ConvertToType[T](gs.Instance())
}

func ToPointerType[T any](gs DynamicStructModifier) *T {
	return dreflect.ConvertToType[*T](gs.New())
}

func (gs *GeneratedStructImpl[T]) generateFields() {
	for k, genFunc := range gs.generatedFields {
		if err := gs.Set(k, genFunc.Generate()); err != nil {
			fmt.Println("E", err)
		}
	}

}

func (gs *GeneratedStructImpl[T]) Instance() T {
	return gs.instance
}

func (gs *GeneratedStructImpl[T]) New() *T {
	gs.DynamicStructModifierImpl.New()
	return &gs.instance
}
