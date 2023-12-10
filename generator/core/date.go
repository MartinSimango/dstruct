package core

import (
	"reflect"
	"time"

	"github.com/MartinSimango/dstruct/generator"
	"github.com/MartinSimango/dstruct/generator/config"
)

const ISO8601 string = "2018-03-20T09:12:28Z"

func GenerateDateTimeFunc() generator.GenerationFunction {

	// TODO have a proper implementation
	return &coreGenerationFunction{
		_func: func(parameters ...any) any {
			return time.Now().UTC().Format(time.RFC3339)
		},
		kind: reflect.String,
	}

}

func GenerateDateTimeBetweenDatesFunc(dc config.DateConfig) generator.GenerationFunction {

	// TODO have a proper implementation
	return &coreGenerationFunction{
		_func: func(parameters ...any) any {
			return time.Now().UTC().Format(time.RFC3339)
		},
		kind: reflect.String,
	}

}

type DateFunctionHolderFunc func(config.DateConfig) generator.GenerationFunction

type DateFunctionHolder struct {
	BaseFunctionHolder
}

func NewDateFunctionHolder(f DateFunctionHolderFunc, cfg config.DateConfig) *DateFunctionHolder {
	return &DateFunctionHolder{
		BaseFunctionHolder: BaseFunctionHolder{
			config:             config.NewConfigBuilder().WithDateConfig(cfg).Build(),
			fun:                f,
			generationFunction: f(cfg),
		},
	}
}

func DefaultDateFunctionHolder(cfg config.DateConfig) *DateFunctionHolder {
	return NewDateFunctionHolder(GenerateDateTimeBetweenDatesFunc, cfg)
}

func (c *DateFunctionHolder) Copy() FunctionHolder {
	return &DateFunctionHolder{
		BaseFunctionHolder: c.BaseFunctionHolder.Copy(),
	}
}
