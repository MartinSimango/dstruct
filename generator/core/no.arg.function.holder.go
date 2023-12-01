package core

import (
	"github.com/MartinSimango/dstruct/generator"
)

type FunctionHolderFuncNoArgs func() generator.GenerationFunction

type FunctionHolderWithNoArgs struct {
	BaseFunctionHolder
	generationFunction generator.GenerationFunction
}

var _ FunctionHolder = &FunctionHolderWithNoArgs{}

func NewFunctionHolderNoArgs(generationFunction generator.GenerationFunction) *FunctionHolderWithNoArgs {
	return &FunctionHolderWithNoArgs{
		generationFunction: generationFunction,
	}
}

func (c *FunctionHolderWithNoArgs) Copy() FunctionHolder {
	return &FunctionHolderWithNoArgs{
		BaseFunctionHolder: c.BaseFunctionHolder.Copy(),
	}
}
