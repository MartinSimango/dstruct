package generator

import (
	"fmt"
	"reflect"
)

type BaseTask struct {
	taskName       string
	parameterCount int
}

func NewBaseTask(taskName string, parameterCount int) *BaseTask {
	return &BaseTask{
		taskName:       taskName,
		parameterCount: parameterCount,
	}
}

func (b *BaseTask) ParameterCount() int {
	return b.parameterCount
}

func (b *BaseTask) GetTags(params ...string) reflect.StructTag {
	b.ValidateParamCount(params...)
	tags := fmt.Sprintf(`gen_task:"%s(%d)"`, b.taskName, len(params))
	for i, p := range params {
		tags += fmt.Sprintf(` gen_task_%d:"%v"`, (i + 1), p)
	}
	return reflect.StructTag(tags)
}

func (b *BaseTask) Name() string {
	return b.taskName
}

func (b *BaseTask) ValidateParamCount(params ...string) {
	if len(params) != b.parameterCount && b.parameterCount != -1 {
		panic(
			fmt.Sprintf(
				"error with task '%s': task requires %d parameters but has %d",
				b.taskName,
				b.parameterCount,
				len(params),
			),
		)
	}
}
