package core

import (
	"reflect"

	"github.com/MartinSimango/dstruct/generator"
)

// GenerateBoolFunc returns a GenerationFunction that generates a boolean value.
func GenerateBoolFunc() generator.GenerationFunction {

	return &coreGenerationFunction{
		_func: func(parameters ...any) any {
			return generateNum(0, 1) == 0
		},
		kind: reflect.Bool,
	}
}
