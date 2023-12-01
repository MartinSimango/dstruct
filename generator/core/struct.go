package core

import (
	"github.com/MartinSimango/dstruct/generator"
)

func GenerateStructFunc(field *GeneratedField) generator.GenerationFunction {

	return &coreGenerationFunction{
		_func: func(parameters ...any) any {
			field := parameters[0].(*GeneratedField)
			field.setStructValues()
			return field.Value.Interface()
		},
		args: []any{field},
	}

}
