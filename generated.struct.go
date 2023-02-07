package dstruct

import (
	"encoding/json"
	"fmt"

	"github.com/MartinSimango/dstruct/generator"
)

type GenerationFields map[string]*generator.GenerationUnit

type GeneratedStruct interface {
	DynamicStructModifier
	// Generate generates fields for the struct
	Generate()

	// GetFieldGenerationConfig gets the generation config for field within the struct.
	GetFieldGenerationConfig(field string) *generator.GenerationConfig

	// SetFieldGenerationConfig sets the generation config for field within the struct. It returns
	// an error if the field does not exist or if the field cannot be generated.
	// Fields that can be generated are struct fields of the most basic type i.e a struct fields
	// that are structs cannot be generated, however it's fields can be.
	//
	// Fields types that cannot be generated: structs, func, chan, any (will default to a nil value being generated).
	//
	// Note: Pointers to structs can be generated.
	SetFieldGenerationConfig(field string, generationConfig *generator.GenerationConfig) error
}

type GeneratedStructImpl struct {
	*DynamicStructModifierImpl
	generatedFields  GenerationFields
	generationConfig *generator.GenerationConfig
}

var _ GeneratedStruct = &GeneratedStructImpl{}

func NewGeneratedStruct(val any) *GeneratedStructImpl {
	return NewGeneratedStructWithConfig(val, generator.NewGenerationConfig())
}

func NewGeneratedStructWithConfig(val any, config *generator.GenerationConfig) *GeneratedStructImpl {
	generatedStruct := &GeneratedStructImpl{
		DynamicStructModifierImpl: ExtendStruct(val).Build().(*DynamicStructModifierImpl),
		generatedFields:           make(GenerationFields),
		generationConfig:          config,
	}
	generatedStruct.populateGeneratedFields()
	return generatedStruct
}

func (gs *GeneratedStructImpl) populateGeneratedFields() {

	for name, field := range gs.fieldMap {
		if field.HasChildren() {
			continue
		}

		gs.generatedFields[name] = generator.NewGenerationUnit(gs.generationConfig, getGeneratorField(field.data))
	}
}

func (gs *GeneratedStructImpl) Generate() {
	gs.generateFields()
}

func (gs *GeneratedStructImpl) SetFieldGenerationConfig(field string, generationConfig *generator.GenerationConfig) error {
	if gs.fieldMap[field] == nil {
		return fmt.Errorf("field %s does not exist within the struct", field)
	}
	if gs.generatedFields[field] == nil {
		return fmt.Errorf("cannot set config for field %s", field)
	}
	gs.generatedFields[field].GenerationConfig = generationConfig
	return nil
}

func (gs *GeneratedStructImpl) GetFieldGenerationConfig(field string) *generator.GenerationConfig {
	return gs.generatedFields[field].GenerationConfig
}

func (gs *GeneratedStructImpl) generateFields() {
	for k, genFunc := range gs.generatedFields {
		if err := gs.Set(k, genFunc.Generate()); err != nil {
			fmt.Println(err)
		}
	}
}

func Print(strct any) string {
	val, _ := json.MarshalIndent(strct, "", "\t")
	return string(val)
}

func getGeneratorField(field *Field) generator.Field {
	return generator.Field{
		Name:  field.fqn,
		Value: field.Value,
		Tag:   field.Tag,
	}
}
