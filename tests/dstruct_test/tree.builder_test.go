package dstruct_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MartinSimango/dstruct"
	"github.com/MartinSimango/dstruct/dreflect"
)

type TestExtendData struct {
	name        string
	input       any
	shouldPanic bool
	expected    any
}

func TestExtend(t *testing.T) {
	type Embedded struct {
		Name string
	}

	type testStruct1 struct {
		Age  int `json:"Age"`
		Name string
		Id   *string
	}

	type testStructEmbedded struct {
		age  int
		name string
		Embedded
	}

	type unexportedStruct struct {
		Name string
	}

	type testStructUnexportedEmbedded struct {
		Age  int `json:"Age"`
		Name string
		unexportedStruct
	}
	tests := []TestExtendData{
		{"ExtendInt", 2, true, nil},
		{"ExtendString", "hello", true, nil},
		{"ExtendReflectValue", reflect.ValueOf(2), true, reflect.ValueOf(2)},
		{"ExtendNil", nil, true, nil},
		{"ExtendBool", true, true, false},
		{"ExtendStructWithAnyNotSet", TestExtendData{}, true, TestExtendData{}},
		{"ExtendStruct", testStruct1{Age: 20}, false, testStruct1{Age: 20}},
		{"ExtendPointerToStruct", &testStruct1{Age: 20}, false, &testStruct1{Age: 20}},
		{"ExtendStructWithEmbeddedField", testStructEmbedded{}, false, testStructEmbedded{}},
		{
			"ExtendStructWithUnexportedEmbeddedField",
			testStructUnexportedEmbedded{},
			true,
			testStructUnexportedEmbedded{},
		},
	}

	assert := assert.New(t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				r := recover()
				assert.Equal(test.shouldPanic, r != nil, r)
			}()
			dynamicStruct := dstruct.ExtendStruct(test.input).Build().Instance()

			dynamicStructConverted := dreflect.Convert(reflect.ValueOf(dynamicStruct), reflect.TypeOf(test.expected)).
				Interface()
			assert.EqualValues(test.expected, dynamicStructConverted)
		})
	}
}

type TestAddFieldData struct {
	name string
}

func TestAddField(t *testing.T) {
	assert := assert.New(t)

	b := dstruct.NewBuilder().
		AddField("Age", 20, `json:"Age"`)

	expected := struct {
		Age int `json:"Age"`
	}{20}

	instance := b.Build().Instance()

	assert.EqualValues(
		expected,
		dreflect.Convert(reflect.ValueOf(instance), reflect.TypeOf(expected)).Interface(),
	)
}

func TestAddEmbeddedField(t *testing.T) {
	assert := assert.New(t)

	type Embedded struct {
		Age    int
		Height float32
	}

	b := dstruct.NewBuilder().
		AddEmbeddedField(Embedded{Age: 20}, `json:"Embedded"`)

	expected := struct {
		Embedded
	}{Embedded: Embedded{Age: 20}}

	instance := b.Build().Instance()

	assert.EqualValues(
		expected,
		dreflect.Convert(reflect.ValueOf(instance), reflect.TypeOf(expected)).Interface(),
	)
}

func TestGetField(t *testing.T) {
	assert := assert.New(t)

	type Person struct {
		Age  int
		Name string
	}

	b := dstruct.NewBuilder().
		AddEmbeddedField(Person{Age: 20, Name: "John"}, ``)

	expected := struct {
		Age  int
		Name string
		Cool int
	}{20, "John", 2}

	instance := b.GetField("Person").AddField("Cool", 2, "").Build()

	fmt.Println(b.Build().GetFields())

	assert.EqualValues(
		expected,
		dreflect.Convert(reflect.ValueOf(instance.Instance()), reflect.TypeOf(expected)).
			Interface(),
	)
	// Original builder must also be altered and have new field
	assert.EqualValues(
		expected,
		dreflect.Convert(reflect.ValueOf(b.GetField("Person").Build().Instance()), reflect.TypeOf(expected)).
			Interface(),
	)
}

// GetNewBuilderFromField returns a new builder instance where the subfield of the struct "field" is the root of the struct.

func TestNewBuilderFromField(t *testing.T) {
	assert := assert.New(t)

	type Person struct {
		Age  int
		Name string
	}

	b := dstruct.NewBuilder().
		AddEmbeddedField(Person{Age: 20, Name: "John"}, ``)

	expected := struct {
		Age  int
		Name string
		Cool int
	}{20, "John", 2}

	expectedOld := struct {
		Age  int
		Name string
	}{20, "John"}

	instance := b.NewBuilderFromField("Person").AddField("Cool", 2, "").Build()

	assert.EqualValues(
		expected,
		dreflect.Convert(reflect.ValueOf(instance.Instance()), reflect.TypeOf(expected)).
			Interface(),
	)
	// Original builder must NOT be altered and have new field
	assert.EqualValues(
		expectedOld,
		dreflect.Convert(reflect.ValueOf(b.GetField("Person").Build().Instance()), reflect.TypeOf(expectedOld)).
			Interface(),
	)
}

func TestRemoveField(t *testing.T) {
	assert := assert.New(t)

	type Person struct {
		Age  int
		Name **struct{ Name string }
	}
	c := &struct{ Name string }{"Cool"}
	b := dstruct.ExtendStruct(Person{Age: 20, Name: &c})

	b.RemoveField("Name.Name")

	d := &struct{}{}
	expected := struct {
		Age  int `json:"Age"`
		Name **struct{}
	}{20, &d}

	instance := b.Build().Instance()

	assert.EqualValues(
		expected,
		dreflect.Convert(reflect.ValueOf(instance), reflect.TypeOf(expected)).Interface(),
	)
}
