package generator

import (
	"fmt"
	"reflect"
)

type Task interface {
	GetTags(params ...string) reflect.StructTag
	ParameterCount() int
	Name() string
	Instance(params ...string) TaskInstance
}

var tasks map[string]Task

func init() {
	tasks = make(map[string]Task)
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
