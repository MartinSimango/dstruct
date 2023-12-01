package dstruct

import (
	"fmt"
	"reflect"

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

	// GetFieldValueGenerationConfig gets the generation config for field within the struct.
	GetFieldValueGenerationConfig(field string) config.GenerationValueConfig

	GetFieldGenerationConfig(field string) config.Config

	GetValueGenerationConfig() config.GenerationValueConfig

	// SetFieldGenerationConfig sets the generation config for field within the struct. It returns
	// an error if the field does not exist or if the field cannot be generated.
	// Fields that can be generated are struct fields of the most basic type i.e a struct fields
	// that are structs cannot be generated, however it's fields can be.
	//
	// Fields types that cannot be generated: structs, func, chan, any (will default to a nil value being generated).
	//
	// Note: Pointers to structs can be generated.
	SetFieldGenerationConfig(field string, generationConfig config.Config) error

	SetFieldGenerationValueConfig(field string, config config.GenerationValueConfig)

	SetGenerationConfig(config config.Config)

	SetGenerationValueConfig(config config.GenerationValueConfig)
}

type GeneratedStructImpl struct {
	*DynamicStructModifierImpl
	generatedFields      GenerationFields
	generatedFieldConfig core.GeneratedFieldConfig
	config               config.Config
}

var _ GeneratedStruct = &GeneratedStructImpl{}

func NewGeneratedStruct(val any) *GeneratedStructImpl {
	return NewGeneratedStructWithConfig(val, config.NewConfig())
}

func NewGeneratedStructWithConfig(val any,
	cfg config.Config) *GeneratedStructImpl {
	generatedStruct := &GeneratedStructImpl{
		DynamicStructModifierImpl: ExtendStruct(val).Build().(*DynamicStructModifierImpl),
		generatedFields:           make(GenerationFields),
		config:                    cfg,
		generatedFieldConfig:      core.NewGenerateFieldConfig(cfg, config.DefaultGenerationValueConfig()),
	}

	generatedStruct.populateGeneratedFields()
	return generatedStruct
}

func (gs *GeneratedStructImpl) populateGeneratedFields() {

	for name, field := range gs.fieldNodeMap {
		if field.HasChildren() {
			continue
		}
		fieldKind := field.data.value.Kind()

		gs.generatedFields[name] = core.NewGenerationUnit(
			core.NewGeneratedField(field.data.fqn,
				field.data.value,
				field.data.tag,
				gs.generatedFieldConfig.Copy(fieldKind),
			))
		if fieldKind == reflect.Slice {
			f := gs.generatedFields[name]
			f.GenerationFunctions[fieldKind] =
				core.NewSliceFunctionHolder(core.GenerateSliceFunc, f.GeneratedField, gs.config)
		}
	}
}

func (gs *GeneratedStructImpl) Generate() {
	gs.generateFields()
}

func (gs *GeneratedStructImpl) GenerateAndUpdate() {
	gs.Generate()
	gs.Update()
}

func (gs *GeneratedStructImpl) changeChildrenConfig(node *Node[structField], cfg config.Config) {
	for k, v := range node.children {
		if v.HasChildren() {
			gs.changeChildrenConfig(v, cfg)
			continue
		}
		gs.generatedFields[k].GenerationFunctions[v.data.typ.Kind()].SetConfig(cfg)
	}

}

func (gs *GeneratedStructImpl) SetFieldGenerationConfig(field string, cfg config.Config) error {
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

func (gs *GeneratedStructImpl) SetFieldGenerationFunction(field string,
	functionHolder core.FunctionHolder) {
	kind := gs.fieldNodeMap[field].data.GetType().Kind()
	generator := gs.generatedFields[field].GeneratedFieldConfig
	generator.GenerationFunctions[kind] = functionHolder
	gs.generatedFields[field].UpdateCurrentFunction = true
}

func (gs *GeneratedStructImpl) SetFieldGenerationFunctions(field string,
	defaultGenerationFunctions core.DefaultGenerationFunctions) {
	gs.generatedFields[field].GenerationFunctions = defaultGenerationFunctions
	gs.generatedFields[field].UpdateCurrentFunction = true
}

func (gs *GeneratedStructImpl) SetGenerationConfig(config config.Config) {
	for name, field := range gs.fieldNodeMap {
		if field.HasChildren() {
			continue
		}
		gs.generatedFields[name].SetConfig(config.Copy())
	}
}

func (gs *GeneratedStructImpl) SetFieldGenerationValueConfig(field string, config config.GenerationValueConfig) {
	gs.generatedFields[field].GenerationValueConfig = config
}

func (gs *GeneratedStructImpl) SetGenerationValueConfig(config config.GenerationValueConfig) {
	for name, field := range gs.fieldNodeMap {
		if field.HasChildren() {
			continue
		}
		gs.generatedFields[name].GenerationValueConfig = config
	}
}

// GetFieldValueGenerationConfig implements GeneratedStruct.
func (gs *GeneratedStructImpl) GetFieldValueGenerationConfig(field string) config.GenerationValueConfig {
	return gs.generatedFields[field].GenerationValueConfig
}

// GetFieldValueGenerationConfig implements GeneratedStruct.
func (gs *GeneratedStructImpl) GetFieldGenerationConfig(field string) config.Config {
	k := gs.generatedFields[field].GeneratedField.Value.Kind()
	return gs.generatedFields[field].GeneratedField.GenerationFunctions[k].GetConfig()
}

// GetValueGenerationConfig implements GeneratedStruct.
func (gs *GeneratedStructImpl) GetValueGenerationConfig() config.GenerationValueConfig {
	return gs.generatedFieldConfig.GenerationValueConfig
}

func (gs *GeneratedStructImpl) generateFields() {
	for k, genFunc := range gs.generatedFields {
		if err := gs.Set(k, genFunc.Generate()); err != nil {
			fmt.Println(err)
		}
	}
}
