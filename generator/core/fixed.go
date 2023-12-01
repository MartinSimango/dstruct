package core

import (
	"github.com/MartinSimango/dstruct/generator"
)

func GenerateFixedValueFunc[T any](n T) generator.GenerationFunction {
	return &coreGenerationFunction{
		_func: func(p ...any) any {
			return n
		},
	}
}
