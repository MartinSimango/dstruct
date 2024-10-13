package generator

import "reflect"

type GenerationFunction interface {
	Generate() any
	Kind() reflect.Kind
}
