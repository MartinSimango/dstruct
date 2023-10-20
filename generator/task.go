package generator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type TaskName string

type Task interface {
	GenerationFunction(taskProperties TaskProperties) GenerationFunction
	ExpectedParameterCount() int
	Name() string
}

var tasks map[string]Task

func init() {
	tasks = make(map[string]Task)
}

type TaskProperties struct {
	TaskName   string
	Parameters []string
	FieldName  string
}

func GetTask(task string) Task {
	return tasks[task]
}

func AddTask(task Task) error {
	name := task.Name()
	if tasks[name] != nil {
		return fmt.Errorf("Task with name %s already exists", name)
	}
	tasks[name] = task
	return nil
}

func getParameters(parameterCount int, tags reflect.StructTag) (parameters []string) {
	for i := 1; i <= parameterCount; i++ {
		parameters = append(parameters, tags.Get(fmt.Sprintf("gen_task_%d", i)))
	}
	return parameters
}

func CreateTaskProperties(fieldName string, tags reflect.StructTag) (*TaskProperties, error) {
	gen_task_tag := strings.TrimSpace(tags.Get("gen_task"))
	leftBraceIndex := strings.Index(gen_task_tag, "(")
	if leftBraceIndex == -1 {
		return nil, fmt.Errorf("error with field %s: task %s error: no ( found", fieldName, gen_task_tag)
	}

	if gen_task_tag[len(gen_task_tag)-1:] != ")" {
		return nil, fmt.Errorf("error with field %s: task %s error: last character of task must be )", fieldName, gen_task_tag)
	}
	taskName := gen_task_tag[:leftBraceIndex]
	parameterCount, err := strconv.Atoi(gen_task_tag[leftBraceIndex+1 : len(gen_task_tag)-1])
	if err != nil {
		return nil, fmt.Errorf("error getting task parameter count: %w", err)
	}

	return &TaskProperties{
		TaskName:   taskName,
		Parameters: getParameters(parameterCount, tags),
		FieldName:  fieldName,
	}, nil
}

func GetTagForTask(name TaskName, params ...any) reflect.StructTag {

	if tasks[string(name)] == nil {
		panic(fmt.Sprintf("Task '%s' is not registered", name))
	}

	tags := fmt.Sprintf(`gen_task:"%s(%d)"`, name, len(params))
	for i, p := range params {
		tags += fmt.Sprintf(` gen_task_%d:"%v"`, (i + 1), p)
	}

	return reflect.StructTag(tags)

}

func ValidateParamCount(task Task, taskProperties TaskProperties) {
	if len(taskProperties.Parameters) != task.ExpectedParameterCount() {
		panic(fmt.Sprintf("error with field %s. task '%s': task requires %d parameters but has %d", taskProperties.FieldName, task.Name(), task.ExpectedParameterCount(), len(taskProperties.Parameters)))
	}

}
