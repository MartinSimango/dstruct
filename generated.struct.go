package dstruct

import (
	"encoding/json"
	"fmt"

	"github.com/MartinSimango/dstruct/generator"
)

type GenerationFields map[string]*generator.GenerationUnit

type GeneratedStruct interface {
	Generate()
	SetGenerationConfig(generationConfig *generator.GenerationConfig)
	SetFieldGenerationConfig(field string, generationConfig *generator.GenerationConfig)
}

type GeneratedStructImpl struct {
	*DynamicStructModifierImpl
	generatedFields  GenerationFields
	GenerationConfig *generator.GenerationConfig
}

var _ GeneratedStruct = &GeneratedStructImpl{}

func NewGeneratedStruct(val any) *GeneratedStructImpl {
	return NewGeneratedStructWithConfig(val, generator.NewGenerationConfig())
}

func NewGeneratedStructWithConfig(val any, config *generator.GenerationConfig) *GeneratedStructImpl {
	generatedStruct := &GeneratedStructImpl{
		DynamicStructModifierImpl: ExtendStruct(val).Build().(*DynamicStructModifierImpl),
		generatedFields:           make(GenerationFields),
		GenerationConfig:          config,
	}
	generatedStruct.populateGeneratedFields()
	return generatedStruct
}

func (gs *GeneratedStructImpl) populateGeneratedFields() {

	for name, field := range gs.fieldMap {
		if field.HasChildren() {
			continue
		}

		gs.generatedFields[name] = generator.NewGenerationUnit(gs.GenerationConfig, getGeneratorField(field.data))
	}
}

func (gs *GeneratedStructImpl) Generate() {
	gs.generateFields()
}

func (gs *GeneratedStructImpl) SetGenerationConfig(generationConfig *generator.GenerationConfig) {
	gs.GenerationConfig = generationConfig
}

func (gs *GeneratedStructImpl) SetFieldGenerationConfig(field string, generationConfig *generator.GenerationConfig) {
	if gs.fieldMap[field] == nil {

	}
	if gs.generatedFields[field] == nil {

	}
	gs.generatedFields[field].GenerationConfig = generationConfig
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
		Value: field.Value,
		Tag:   field.Tag,
	}
}
