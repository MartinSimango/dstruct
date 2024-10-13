package core

import (
	"github.com/MartinSimango/dstruct/generator"
	"github.com/MartinSimango/dstruct/generator/config"
)

type FunctionHolderFuncNoArgs func() generator.GenerationFunction

type FunctionHolderWithNoArgs struct {
	BaseFunctionHolder
}

var _ FunctionHolder = &FunctionHolderWithNoArgs{}

func NewFunctionHolderNoArgs(generationFunction generator.GenerationFunction) *FunctionHolderWithNoArgs {
	return &FunctionHolderWithNoArgs{
		BaseFunctionHolder: BaseFunctionHolder{
			generationFunction: generationFunction,
		},
	}
}

func (c *FunctionHolderWithNoArgs) GetFunction() generator.GenerationFunction {
	return c.generationFunction
}

func (c *FunctionHolderWithNoArgs) SetConfig(config config.Config) {}

func (c *FunctionHolderWithNoArgs) Copy() FunctionHolder {
	return &FunctionHolderWithNoArgs{
		BaseFunctionHolder: c.BaseFunctionHolder.Copy(),
	}
}
