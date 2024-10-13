package core

import (
	"github.com/MartinSimango/dstruct/generator"
)

type FixedFunctionHolderFunc[T any] func(value T) generator.GenerationFunction

type FixedFunctionHolder[T any] struct {
	value T
	BaseFunctionHolder
}

var _ FunctionHolder = &FixedFunctionHolder[int]{}

func NewFixedFunctionHolder[T any](f FixedFunctionHolderFunc[T], value T) *FixedFunctionHolder[T] {
	return &FixedFunctionHolder[T]{
		BaseFunctionHolder: BaseFunctionHolder{
			fun: f,
		},
		value: value,
	}

}

// Override
func (c *FixedFunctionHolder[T]) GetFunction() generator.GenerationFunction {
	if c.generationFunction != nil {
		return c.generationFunction
	}
	c.generationFunction = c.fun.(FixedFunctionHolderFunc[T])(c.value)

	return c.generationFunction
}

func (c *FixedFunctionHolder[T]) Copy() FunctionHolder {
	return &FixedFunctionHolder[T]{
		BaseFunctionHolder: c.BaseFunctionHolder.Copy(),
		value:              c.value,
	}
}
