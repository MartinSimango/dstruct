package core

import (
	"github.com/MartinSimango/dstruct/generator"
	"github.com/MartinSimango/dstruct/generator/config"
)

type PointerFunctionHolderFunc func(*GeneratedField) generator.GenerationFunction

type PointerFunctionHolder struct {
	BaseFunctionHolder
	field *GeneratedField
}

var _ FunctionHolder = &PointerFunctionHolder{}

func NewPointerFunctionHolder(
	f PointerFunctionHolderFunc,
	field *GeneratedField,
) FunctionHolder {
	return &PointerFunctionHolder{
		BaseFunctionHolder: BaseFunctionHolder{
			config:        nil,
			fun:           f,
			resetFunction: nil,
		},
		field: field,
	}
}

func (c *PointerFunctionHolder) GetFunction() generator.GenerationFunction {
	if c.generationFunction != nil {
		return c.generationFunction
	}
	c.generationFunction = c.fun.(PointerFunctionHolderFunc)(c.field)
	return c.generationFunction
}

func (c *PointerFunctionHolder) Copy(cfg config.Config) FunctionHolder {
	return &SliceFunctionHolder{
		BaseFunctionHolder: c.BaseFunctionHolder.Copy(cfg),
	}
}
