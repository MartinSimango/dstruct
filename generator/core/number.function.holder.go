package core

import (
	"github.com/MartinSimango/dstruct/generator"
	"github.com/MartinSimango/dstruct/generator/config"
)

// TODO: rename to NumberFunctionHolderFunc?
type NumberFunctionHolderFunc func(config.NumberRangeConfig) generator.GenerationFunction

type NumberFunctionHolder struct {
	BaseFunctionHolder
}

var _ FunctionHolder = &NumberFunctionHolder{}

func NewNumberFunctionHolder(
	f NumberFunctionHolderFunc,
	cfg config.NumberRangeConfig,
) *NumberFunctionHolder {
	return &NumberFunctionHolder{
		BaseFunctionHolder: BaseFunctionHolder{
			config:             config.NewDstructConfigBuilder().WithNumberConfig(cfg).Build(),
			fun:                f,
			generationFunction: f(cfg),
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
