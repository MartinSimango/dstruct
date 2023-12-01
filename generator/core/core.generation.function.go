package core

import (
	"github.com/MartinSimango/dstruct/generator"
)

type coreGenerationFunction struct {
	_func func(...any) any
	args  []any
}

var _ generator.GenerationFunction = &coreGenerationFunction{}

func (cgf *coreGenerationFunction) Generate() any {
	return cgf._func(cgf.args...)
}
