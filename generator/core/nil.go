package core

import (
	"github.com/MartinSimango/dstruct/generator"
)

func GenerateNilValueFunc() generator.GenerationFunction {
	return &coreGenerationFunction{
		_func: func(parameters ...any) any {
			return nil
		},
	}

}
