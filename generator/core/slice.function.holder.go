package core

import (
	"github.com/MartinSimango/dstruct/generator"
	"github.com/MartinSimango/dstruct/generator/config"
)

type SliceFunctionHolderFunc func(*GeneratedField, config.Config) generator.GenerationFunction

type SliceFunctionHolder struct {
	BaseFunctionHolder
	field *GeneratedField
}

var _ FunctionHolder = &SliceFunctionHolder{}

func NewSliceFunctionHolder(f SliceFunctionHolderFunc, field *GeneratedField, cfg config.Config) *SliceFunctionHolder {
	return &SliceFunctionHolder{
		BaseFunctionHolder: BaseFunctionHolder{
			config: cfg,
			fun:    f,
		},
		field: field,
	}
}

func (c *SliceFunctionHolder) GetFunction() generator.GenerationFunction {
	if c.generationFunction != nil {
		return c.generationFunction
	}
	c.generationFunction = c.fun.(SliceFunctionHolderFunc)(c.field, c.config)

	return c.generationFunction
}

func (c *SliceFunctionHolder) Copy() FunctionHolder {
	return &NumberFunctionHolder{
		BaseFunctionHolder: c.BaseFunctionHolder.Copy(),
	}
}
