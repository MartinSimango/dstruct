package core

import (
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
		kind: NewKind(time.Time{}),
	}
}

func GenerateDateTimeBetweenDatesFunc(dc config.DateRangeConfig) generator.GenerationFunction {
	// TODO have a proper implementation
	return &coreGenerationFunction{
		_func: func(parameters ...any) any {
			timeDiffInSeconds := dc.GetEndDate().Sub(dc.GetStartDate()).Seconds()
			return dc.GetStartDate().
				Add(time.Second * time.Duration(generateNum(0, timeDiffInSeconds)))
		},
		kind: NewKind(time.Time{}),
	}
}

type DateFunctionHolderFunc func(config.DateRangeConfig) generator.GenerationFunction

type DateFunctionHolder struct {
	BaseFunctionHolder
}

func NewDateFunctionHolder(
	f DateFunctionHolderFunc,
	cfg config.DateRangeConfig,
) *DateFunctionHolder {
	return &DateFunctionHolder{
		BaseFunctionHolder: BaseFunctionHolder{
			config:             config.NewDstructConfigBuilder().WithDateRangeConfig(cfg).Build(),
			fun:                f,
			generationFunction: f(cfg),
		},
	}
}

func DefaultDateFunctionHolder(cfg config.DateRangeConfig) *DateFunctionHolder {
	return NewDateFunctionHolder(GenerateDateTimeBetweenDatesFunc, cfg)
}

func (c *DateFunctionHolder) Copy(cfg config.Config) FunctionHolder {
	return &DateFunctionHolder{
		BaseFunctionHolder: c.BaseFunctionHolder.Copy(cfg),
	}
}
