package core

import (
	"reflect"

	"github.com/MartinSimango/dstruct/generator"
)

func GenerateFixedValueFunc[T any](n T) generator.GenerationFunction {
	return &coreGenerationFunction{
		_func: func(p ...any) any {
			return n
		},
		kind: reflect.ValueOf(n).Kind(),
	}
}
