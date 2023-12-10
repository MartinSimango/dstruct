package core

import (
	"reflect"

	"github.com/MartinSimango/dstruct/generator"
)

type coreGenerationFunction struct {
	_func func(...any) any
	args  []any
	kind  reflect.Kind
}

var _ generator.GenerationFunction = &coreGenerationFunction{}

func (cgf *coreGenerationFunction) Generate() any {
	return cgf._func(cgf.args...)
}

func (cgf *coreGenerationFunction) Kind() reflect.Kind {
	return cgf.kind
}
