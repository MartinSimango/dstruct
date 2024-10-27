package core

import (
	"github.com/MartinSimango/dstruct/generator"
	"github.com/MartinSimango/dstruct/generator/config"
)

type SliceFunctionHolderFunc func(*GeneratedField, config.Config, DefaultGenerationFunctions) generator.GenerationFunction

type SliceFunctionHolder struct {
	BaseFunctionHolder
	field               *GeneratedField
	generationFunctions DefaultGenerationFunctions
}

var _ FunctionHolder = &SliceFunctionHolder{}

func NewSliceFunctionHolder(
	f SliceFunctionHolderFunc,
	field *GeneratedField,
	cfg config.Config,
	generationFunctions DefaultGenerationFunctions,
) *SliceFunctionHolder {
	return &SliceFunctionHolder{
		BaseFunctionHolder: BaseFunctionHolder{
			config: cfg,
			fun:    f,
			resetFunction: func(cfg config.Config) generator.GenerationFunction {
				return f(field, cfg, generationFunctions)
			},
		},
		field:               field,
		generationFunctions: generationFunctions,
	}
}

func (c *SliceFunctionHolder) GetFunction() generator.GenerationFunction {
	if c.generationFunction != nil {
		return c.generationFunction
	}
	c.generationFunction = c.fun.(SliceFunctionHolderFunc)(c.field, c.config, c.generationFunctions)

	return c.generationFunction
}

func (c *SliceFunctionHolder) Copy(cfg config.Config) FunctionHolder {
	return &SliceFunctionHolder{
		BaseFunctionHolder: c.BaseFunctionHolder.Copy(cfg),
		// TODO address this
		// generationFunctions: c.generationFunctions.Copy(reflect.Slice),
	}
}
