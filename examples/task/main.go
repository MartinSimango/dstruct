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
	GenInt32 generator.TaskName = "GenInt32"
)

type genInt32Params struct {
	min int32
	max int32
}

type GenInt32Task struct{}

// Ensure GenInt32Task implements Task.
var _ generator.Task = &GenInt32Task{}

// Tags implements Task.
func (g *GenInt32Task) Name() string {
	return string(GenInt32)
}

// GenerationFunction implements Task.
func (g *GenInt32Task) GenerationFunction(
	taskProperties generator.TaskProperties,
) generator.GenerationFunction {
	generator.ValidateParamCount(g, taskProperties)
	params := g.getInt32Params(taskProperties.FieldName, taskProperties.Parameters)
	numberConfig := config.NewNumberRangeConfig()
	numberConfig.Int().SetRange(int(params.min), int(params.max))
	return core.GenerateNumberFunc[int](numberConfig)
}

func (g *GenInt32Task) ExpectedParameterCount() int {
	return 2
}

func (g *GenInt32Task) FunctionHolder(taskProperties generator.TaskProperties) core.FunctionHolder {
	return core.NewFunctionHolderNoArgs(g.GenerationFunction(taskProperties))
}

func (g *GenInt32Task) getInt32Params(fieldName string, params []string) genInt32Params {
	param_1, err := strconv.Atoi(params[0])
	if err != nil {
		panic(fmt.Sprintf("error with field %s: task %s error: %s", fieldName, GenInt32, err))
	}

	param_2, err := strconv.Atoi(params[1])
	if err != nil {
		panic(fmt.Sprintf("error with field %s: task %s error: %s", fieldName, GenInt32, err))
	}

	if param_1 > param_2 {
		err = fmt.Errorf(
			"min must be less or equal to the max value min = %d max = %d",
			param_1,
			param_2,
		)
		panic(fmt.Sprintf("error with field %s: task %s error: %s", fieldName, GenInt32, err))

	}

	return genInt32Params{
		min: int32(param_1),
		max: int32(param_2),
	}
}

type M struct {
	Name string
}

type Person struct {
	// Love ABC
	Age   *int
	Time  time.Time
	Other *M

	// Parray []Person

	// Person *Person
}

type P struct {
	P     Person
	Value int
}

type Test struct {
	// A      *int
	// S      string
	// C      int
	Person P `json:"person"`

	// Cpoint *int
	T time.Time

	// Pa     []Pa
	// Parray []Person

	// L
}

func main() {
	e := dstruct.ExtendStruct(Test{}).Build().Instance()
	createTime := time.Now()
	// for i := 0; i < 1; i++ {
	// 	estruct.AddField(fmt.Sprintf("Test_%d", i), Test{}, "")
	// }
	fmt.Println(time.Since(createTime))
	// b := estruct.Build()
	//
	st := config.GenerationSettings{SetNonRequiredFields: true}
	gestruct := dstruct.NewGeneratedStructWithConfig(e,
		config.NewDstructConfig(),
		st,
	)

	c := gestruct.GetGenerationConfig()
	c.SetIntRange(20, 50)
	st.SetNonRequiredFields = false

	// gestruct.SetFieldGenerationSettings("Person", st)
	// gestruct.SetFieldGenerationConfig("Person", c)

	gTime := time.Now()
	gestruct.Generate()
	// gestruct.Update()
	// gestruct.Set("Person.P", Person{Age: new(int), Time: time.Now()})
	err := gestruct.Set(
		"Person.P",
		Person{Age: new(int), Time: time.Now(), Other: &M{Name: "Martin"}},
	)
	if err != nil {
		// panic(err)
	}

	fmt.Println("Time to generate: ", time.Since(gTime))
	fmt.Printf("%+v\n", gestruct.Get_("Person.P.Other.Name"))
	// for field := range gestruct.GetFields() {
	// fmt.Println("Field: ", field)
	// fmt.Printf(
	// 	"'%s': dstruct: %+v  goType: %+v\n",
	// 	field,
	// 	value.GetDstructType(),
	// 	value.GetGoType(),
	// )
	// }

	// generatedStuct := dstruct.NewGeneratedStructWithConfig(
	// 	Test{Cpoint: new(int)},
	// 	config.NewDstructConfig().SetSliceLength(3, 3),
	// 	config.DefaultGenerationSettings(),
	// )
	// gt := &GenInt32Task{}
	// generator.AddTask(gt)
	// if err := generatedStuct.Set("Person.P", &Person{}); err != nil {
	// 	panic(err)
	// }
	// generatedStuct.SetFieldGenerationConfig(
	// 	"Person.Value",
	// 	config.NewDstructConfig().SetIntRange(800, 1000),
	// )
	//
	// // generatedStuct.SetFieldFromTask("C", gt, 300, 400)
	// // generatedStuct.Update()
	// generatedStuct.Generate()
	//
	// fmt.Printf("%+v\n", generatedStuct.Get_("Person.P"))
	//
	fmt.Println("Testing task")
}
