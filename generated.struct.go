package dstruct

import (
	"github.com/MartinSimango/dstruct/generator"
	"github.com/MartinSimango/dstruct/generator/config"
	"github.com/MartinSimango/dstruct/generator/core"
)

type GeneratedFieldContexts map[string]*core.GeneratedFieldContext

type GeneratedStruct interface {
	DynamicStructModifier
	// Generate generates fields for the struct. If new fields are generated, the root tree for the underlying
	// struct is updated. This allows new generated fields to be accessed and modified by Set and Get methods
	Generate()

	// SetFieldGenerationSettings sets the generation value config for field within the struct.
	// If the field does not exist or if the field has no generation settings an error is returned.
	SetFieldGenerationSettings(
		field string,
		settings config.GenerationSettings,
	) error

	// GetFieldGenerationSettings gets the generation config for field within the struct.
	// If the field does not exist or if the field has no generation settings an error is returned.
	GetFieldGenerationSettings(field string) (config.GenerationSettings, error)

	// GetFieldGenerationSettings_ is like GetFieldGenerationSettings but panics if an error occurs.
	GetFieldGenerationSettings_(field string) config.GenerationSettings

	// SetGenerationSettings sets the generation settings for the struct and propagates the settings to all fields

	SetGenerationSettings(settings config.GenerationSettings)

	// GetGenerationSettings gets the generation settings for the struct.
	GetGenerationSettings() config.GenerationSettings

	// SetFieldGenerationConfig sets the generation config for field within the struct. It returns
	// an error if the field does not exist or if the field cannot be generated.
	// Fields that can be generated are struct fields of the most basic type i.e a struct fields
	// that are structs cannot be generated, however it's fields can be.
	//
	// Fields types that cannot be generated: structs, func, chan, any (will default to a nil value being generated).
	//
	// Note: Pointers to structs can be generated.
	SetFieldGenerationConfig(field string, generationConfig config.Config) error

	// GetFieldGenerationConfig gets the generation config for field within the struct.
	GetFieldGenerationConfig(field string) (config.Config, error)

	// GetFieldGenerationConfig_ is like GetFieldGenerationConfig but panics if an error occurs.
	GetFieldGenerationConfig_(field string) config.Config

	// SetGenerationConfig sets the generation config for the struct and propagates the settings to all fields
	SetGenerationConfig(config config.Config)

	// GetGenerationConfig gets the generation config for the struct.
	GetGenerationConfig() config.Config

	// SetFieldGenerationFunction sets the generation function for field within the struct. It returns an error if the field does not exist or if the field cannot be generated.
	SetFieldGenerationFunction(field string, functionHolder core.FunctionHolder) error

	// // GetFieldGenerationConfig gets the generation function for field within the struct.
	// GetFieldGenerationFunction(field string) (core.FunctionHolder, error)
	//
	// // GetFieldGenerationFunction_ is like GetFieldGenerationFunction but panics if an error occurs.
	// GetFieldGenerationFunction_(field string) core.FunctionHolder
	//
	// SetFieldDefaultFunctions sets the default generation functions for field within the struct. It returns an error if the field does not exist or if the field cannot be generated.
	SetFieldGenerationFunctions(field string, functions core.DefaultGenerationFunctions) error

	// SetGenerationFunctions sets the generation functions for the struct and propagates the settings to all fields
	SetGenerationFunctions(functions core.DefaultGenerationFunctions)

	// SetFieldFromTask sets the field value from the task. The task is used to generate the value for the field.
	SetFieldFromTask(field string, task generator.Task, params ...any) error
}
