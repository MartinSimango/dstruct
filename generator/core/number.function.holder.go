package core

import (
	"github.com/MartinSimango/dstruct/generator"
	"github.com/MartinSimango/dstruct/generator/config"
)

type NumberFunctionHolderFunc func(config.NumberRangeConfig) generator.GenerationFunction

type NumberFunctionHolder struct {
	BaseFunctionHolder
}

var _ FunctionHolder = &NumberFunctionHolder{}

func NewNumberFunctionHolder(
	f NumberFunctionHolderFunc,
	cfg config.NumberRangeConfig,
) *NumberFunctionHolder {
	nfh := &NumberFunctionHolder{
		BaseFunctionHolder: BaseFunctionHolder{
			config: config.NewDstructConfigBuilder().WithNumberConfig(cfg).Build(),
			fun:    f,
			resetFunction: func(cfg config.Config) generator.GenerationFunction {
				// fmt.Println("Resetting number function holder")
				return f(cfg.Number())
			},
			generationFunction: f(cfg),
		},
	}

	return nfh
}

func (c *NumberFunctionHolder) Copy(cfg config.Config) FunctionHolder {
	nf := &NumberFunctionHolder{
		BaseFunctionHolder: c.BaseFunctionHolder.Copy(cfg),
	}
	return nf
}
