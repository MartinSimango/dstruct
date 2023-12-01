package core

import (
	"github.com/MartinSimango/dstruct/generator"
	"github.com/MartinSimango/dstruct/generator/config"
)

type NumberFunctionHolderFunc func(config.NumberConfig) generator.GenerationFunction

type NumberFunctionHolder struct {
	BaseFunctionHolder
}

var _ FunctionHolder = &NumberFunctionHolder{}

func NewNumberFunctionHolder(f NumberFunctionHolderFunc, cfg config.NumberConfig) *NumberFunctionHolder {
	return &NumberFunctionHolder{
		BaseFunctionHolder: BaseFunctionHolder{
			config: config.NewConfigBuilder().WithNumberConfig(cfg).Build(),
			fun:    f,
		},
	}
}

// Override
func (c *NumberFunctionHolder) GetFunction() generator.GenerationFunction {
	if c.generationFunction != nil {
		return c.generationFunction
	}
	c.generationFunction = c.fun.(NumberFunctionHolderFunc)(c.config.Number())

	return c.generationFunction
}

func (c *NumberFunctionHolder) Copy() FunctionHolder {
	return &NumberFunctionHolder{
		BaseFunctionHolder: c.BaseFunctionHolder.Copy(),
	}
}
