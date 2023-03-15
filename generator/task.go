package generator

import (
	"fmt"
	"strconv"
	"strings"
)

type TaskName string

const (
	GenInt32 TaskName = "GenInt32"
)

type Task struct {
	Name       TaskName
	Parameters string
	FieldName  string
}

type GenInt32Params struct {
	min int32
	max int32
}

func (t *Task) GenInt32Params() GenInt32Params {
	params := strings.Split(t.Parameters, ",")

	if len(params) != 2 {
		panic(fmt.Sprintf("error with field %s: task %s: task requires 2 parameters but has %d", t.FieldName, t.Name, len(params)))
	}
	param_1, err := strconv.Atoi(params[0])
	if err != nil {
		panic(fmt.Sprintf("error with field %s: task %s error: %s", t.FieldName, t.Name, err))
	}

	param_2, err := strconv.Atoi(params[1])
	if err != nil {
		panic(fmt.Sprintf("error with field %s: task %s error: %s", t.FieldName, t.Name, err))
	}

	if param_1 > param_2 {
		err = fmt.Errorf("min must be less or equal to the max value min = %d max = %d", param_1, param_2)
		panic(fmt.Sprintf("error with field %s: task %s error: %s", t.FieldName, t.Name, err))

	}

	return GenInt32Params{
		min: int32(param_1),
		max: int32(param_2),
	}
}

func getTask(task string, fieldName string) Task {
	task = strings.TrimSpace(task)
	leftBraceIndex := strings.Index(task, "(")
	if leftBraceIndex == -1 {
		panic(fmt.Sprintf("error with field %s: task %s error: no ( found", fieldName, task))
	}

	if task[len(task)-1:] != ")" {
		panic(fmt.Sprintf("error with field %s: task %s error: last character of task must be )", fieldName, task))
	}
	taskName := task[:leftBraceIndex]
	parameters := task[leftBraceIndex+1 : len(task)-1]
	return Task{
		Name:       TaskName(taskName),
		Parameters: parameters,
		FieldName:  fieldName,
	}
}

func (t Task) getFunction() GenerationFunction {
	switch t.Name {
	case GenInt32:
		params := t.GenInt32Params()
		return GenerateNumberFunc(params.min, params.max, NewGenerationConfig())
	}
	panic(fmt.Sprintf("Invalid task name '%s' for field %s ", t.Name, t.FieldName))
}

func GetTagsForGenInt32Task(min int32, max int32) string {
	return fmt.Sprintf(`gen_task:"%s(%d,%d)"`, GenInt32, min, max)
}
