package dstruct

import (
	"fmt"

	"github.com/MartinSimango/dstruct/generator"
)

type GenerationFields map[string]*generator.GenerationUnit

type GeneratedStruct interface {
	DynamicStructModifier
	// Generate generates fields for the struct
	Generate()

	// GenerateAndUpdate Generates fields and updates the root tree for the underlying struct. Allowing
	// new generated fields to be accessed and modified by Set and Get methods.
	GenerateAndUpdate()

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
	generatedFields GenerationFields
	generator       *generator.Generator
}

var _ GeneratedStruct = &GeneratedStructImpl{}

func NewGeneratedStruct(val any) *GeneratedStructImpl {
	return NewGeneratedStructWithConfig(val, generator.NewGenerationFunctionDefaults(generator.NewGenerationConfig()))
}

func NewGeneratedStructWithConfig(val any,
	gen *generator.Generator) *GeneratedStructImpl {
	generatedStruct := &GeneratedStructImpl{
		DynamicStructModifierImpl: ExtendStruct(val).Build().(*DynamicStructModifierImpl),
		generatedFields:           make(GenerationFields),
		generator:                 gen,
	}
	generatedStruct.populateGeneratedFields()
	return generatedStruct
}

func (gs *GeneratedStructImpl) populateGeneratedFields() {

	for name, field := range gs.fieldMap {
		if field.HasChildren() || field.data.isFQNAlias {
			continue
		}

		gs.generatedFields[name] = generator.NewGenerationUnit(
			generator.NewGeneratedField(field.data.fqn,
				field.data.value,
				field.data.tag,
				gs.generator.Copy(),
			))
	}
}

func (gs *GeneratedStructImpl) Generate() {
	gs.generateFields()
}

func (gs *GeneratedStructImpl) GenerateAndUpdate() {
	gs.Generate()
	gs.Update()
	// for k, v := range gs.GetFields() {
	// 	fmt.Println(k, v.anonymous, v.isFQNAlias)
	// }
}

func (gs *GeneratedStructImpl) SetFieldGenerationConfig(field string, generationConfig *generator.GenerationConfig) error {
	if gs.fieldMap[field] == nil {
		return fmt.Errorf("field %s does not exist within the struct", field)
	}
	if gs.generatedFields[field] == nil {
		return fmt.Errorf("cannot set config for field %s", field)
	}
	gs.generatedFields[field].Field.Generator.GenerationConfig = generationConfig
	return nil
}

func (gs *GeneratedStructImpl) GetFieldGenerationConfig(field string) *generator.GenerationConfig {
	return gs.generatedFields[field].Field.Generator.GenerationConfig
}

func (gs *GeneratedStructImpl) GetFieldGenerator(field string) *generator.Generator {
	return gs.generatedFields[field].Field.Generator
}
func (gs *GeneratedStructImpl) SetFieldDefaultGenerationFunction(field string,
	generationFunction generator.GenerationFunction) {
	kind := gs.fieldMap[field].data.GetType().Kind()
	generator := gs.generatedFields[field].Field.Generator
	generator.DefaultGenerationFunctions[kind] = generationFunction
	generationFunction.SetGenerationConfig(generator.GenerationConfig)
	gs.generatedFields[field].UpdateCurrentFunction = true
}

func (gs *GeneratedStructImpl) SetFieldGenerator(field string,
	generator *generator.Generator) {
	gs.generatedFields[field].Field.Generator = generator
	gs.generatedFields[field].UpdateCurrentFunction = true

}

func (gs *GeneratedStructImpl) generateFields() {
	for k, genFunc := range gs.generatedFields {
		if err := gs.Set(k, genFunc.Generate()); err != nil {
			fmt.Println(err)
		}
	}
}

func getGeneratorField(field *field, defaultGenerationFunction *generator.Generator) *generator.GeneratedField {
	return &generator.GeneratedField{
		Name:      field.fqn,
		Value:     field.value,
		Tag:       field.tag,
		Generator: defaultGenerationFunction,
	}
}
