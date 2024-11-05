package generator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type TaskProperties struct {
	TaskName   string
	parameters []string
}

// CreateTaskPropertiesFromTag creates a TaskProperties struct from a reflect.StructTag
// structure of gen task tag is as follows: gen_task:"task_name(parameter_count)" gen_task_1:"parameter_1" gen_task_2:"parameter_2"
func CreateTaskPropertiesFromTag(tag reflect.StructTag) (*TaskProperties, error) {
	gen_task_tag := strings.TrimSpace(tag.Get("gen_task"))
	leftBraceIndex := strings.Index(gen_task_tag, "(")
	if leftBraceIndex == -1 {
		return nil, fmt.Errorf(
			"error creating task properties for task %s: no ( found",
			gen_task_tag,
		)
	}

	if gen_task_tag[len(gen_task_tag)-1:] != ")" {
		return nil, fmt.Errorf(
			"error creating task properties for task %s: last character of task must be )",
			gen_task_tag,
		)
	}
	taskName := gen_task_tag[:leftBraceIndex]
	parameterCount, err := strconv.Atoi(gen_task_tag[leftBraceIndex+1 : len(gen_task_tag)-1])
	if err != nil {
		return nil, fmt.Errorf("error getting task parameter count: %w", err)
	}

	return &TaskProperties{
		TaskName:   taskName,
		parameters: getParameters(parameterCount, tag),
	}, nil
}

func CreateTaskPropertiesFromParams(taskName string, params ...any) *TaskProperties {
	tp := &TaskProperties{
		TaskName: taskName,
	}
	tp.SetParameters(params...)
	return tp
}

func (tp *TaskProperties) SetParameters(params ...any) {
	p := make([]string, len(params))
	for i, param := range params {
		p[i] = fmt.Sprintf("%v", param)
	}
	tp.parameters = p
}

func (tp *TaskProperties) Parameters() []string {
	return tp.parameters
}

func getParameters(parameterCount int, tags reflect.StructTag) (parameters []string) {
	for i := 1; i <= parameterCount; i++ {
		parameters = append(parameters, tags.Get(fmt.Sprintf("gen_task_%d", i)))
	}
	return parameters
}
