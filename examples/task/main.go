package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/MartinSimango/dstruct"
	"github.com/MartinSimango/dstruct/generator"
	"github.com/MartinSimango/dstruct/generator/config"
	"github.com/MartinSimango/dstruct/generator/core"
)

const (
	GenInt32 string = "GenInt32"
)

type genInt32Params struct {
	min int32
	max int32
}

type GenInt32Task struct {
	generator.BaseTask
	numberConfig config.NumberRangeConfig
}

type GenInt32TaskInstance struct {
	GenInt32Task
	genFunc generator.GenerationFunction
}

var (
	_ generator.TaskInstance = &GenInt32TaskInstance{}
	// Ensure GenInt32Task implements Task.
	_ generator.Task = &GenInt32Task{}
)

func NewGenInt32Task() *GenInt32Task {
	gt := &GenInt32Task{
		numberConfig: config.NewNumberRangeConfig(),
		BaseTask:     *generator.NewBaseTask(string(GenInt32), 2),
	}
	return gt
}

func (g GenInt32Task) Instance(params ...string) generator.TaskInstance {
	gt := &GenInt32TaskInstance{
		GenInt32Task: g,
	}
	gt.numberConfig = g.numberConfig.Copy()
	gt.SetParameters(params...)
	gt.genFunc = core.GenerateNumberFunc[int32](gt.numberConfig)
	return gt
}

func (g *GenInt32TaskInstance) GenerationFunction() generator.GenerationFunction {
	return g.genFunc
}

func (g *GenInt32TaskInstance) SetParameters(params ...string) {
	g.ValidateParamCount(params...)
	min, err := strconv.Atoi(params[0])
	if err != nil {
		panic(fmt.Sprintf("Error parsing min value: %s", err))
	}
	max, err := strconv.Atoi(params[1])
	if err != nil {
		panic(fmt.Sprintf("Error parsing max value: %s", err))
	}

	if min > max {
		panic(fmt.Sprintf("min value %d is greater than max value %d", min, max))
	}
	g.numberConfig.Int32().SetRange(int32(min), int32(max))
}

type M struct {
	Name string
}

type Person struct {
	Age   *int
	Time  time.Time
	Other *M
}

type P struct {
	P     Person
	Value int
}

type Test struct {
	B      int32 `gen_task:"GenInt32(2)" gen_task_1:"10" gen_task_2:"20"`
	C      int32
	Person P
	Cpoint *int
	T      time.Time
}

func main() {
	generatedStuct := dstruct.NewGeneratedStructWithConfig(
		Test{Cpoint: new(int)},
		config.NewDstructConfig().SetSliceLength(3, 3),
		config.DefaultGenerationSettings(),
	)
	gt := NewGenInt32Task()
	generator.AddTask(gt)
	if err := generatedStuct.Set("Person.P", Person{}); err != nil {
		panic(err)
	}
	generatedStuct.SetFieldGenerationConfig(
		"Person.Value",
		config.NewDstructConfig().SetIntRange(800, 1000),
	)
	gti := gt.Instance("20", "30").(*GenInt32TaskInstance)
	generatedStuct.SetFieldFromTaskInstance("C", gti)
	generatedStuct.Generate()

	fmt.Printf("%+v\n", generatedStuct)
	gti.SetParameters("100", "200")

	generatedStuct.Generate()

	fmt.Printf("%+v\n", generatedStuct)
}
