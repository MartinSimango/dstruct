package generator

type TaskInstance interface {
	GenerationFunction() GenerationFunction
	SetParameters(params ...string)
}
