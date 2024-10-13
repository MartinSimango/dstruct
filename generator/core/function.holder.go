package core

import (
	"reflect"

	"github.com/MartinSimango/dstruct/generator"
	"github.com/MartinSimango/dstruct/generator/config"
)

type FunctionHolder interface {
	GetFunction() generator.GenerationFunction
	SetFunction(generationFunction generator.GenerationFunction)
	GetConfig() config.Config
	SetConfig(cfg config.Config)
	Copy() FunctionHolder
	Kind() reflect.Kind
}

type ResetFunction func(config config.Config) generator.GenerationFunction

type BaseFunctionHolder struct {
	config             config.Config
	generationFunction generator.GenerationFunction
	fun                any
	resetFunction      ResetFunction
}

func (c *BaseFunctionHolder) SetFunction(generationFunction generator.GenerationFunction) {
	c.generationFunction = generationFunction
}

func (c *BaseFunctionHolder) GetFunction() generator.GenerationFunction {
	return c.generationFunction
}

func (c *BaseFunctionHolder) GetConfig() config.Config {
	return c.config
}

func (c *BaseFunctionHolder) SetConfig(cfg config.Config) {
	c.config = cfg
	c.generationFunction = c.resetFunction(cfg)
}

func (c *BaseFunctionHolder) Copy() BaseFunctionHolder {
	// TODO: should this panic?
	var configCopy config.Config
	if c.config != nil {
		configCopy = c.config.Copy()
	}

	return BaseFunctionHolder{
		fun:                c.fun,
		config:             configCopy,
		generationFunction: c.resetFunction(configCopy),
	}
}

func (c *BaseFunctionHolder) Kind() reflect.Kind {
	return c.generationFunction.Kind()
}
