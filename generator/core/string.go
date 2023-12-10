package core

import (
	"reflect"

	"github.com/MartinSimango/dstruct/generator"
	"github.com/takahiromiyamoto/go-xeger"
)

func GenerateStringFromRegexFunc(regex string) generator.GenerationFunction {
	return &coreGenerationFunction{
		_func: func(parameters ...any) any {
			regex := parameters[0].(string)
			x, err := xeger.NewXeger(regex)
			if err != nil {
				panic(err)
			}
			return x.Generate()
		},
		args: []any{regex},
		kind: reflect.Slice,
	}
}
