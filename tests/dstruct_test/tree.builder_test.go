package dstruct_test

import (
	"reflect"
	"testing"

	"github.com/MartinSimango/dstruct"
	"github.com/MartinSimango/dstruct/dreflect"
	"github.com/stretchr/testify/assert"
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
	}

	type testStructEmbedded struct {
		age  int `json:"Age"`
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
	var tests = []TestExtendData{
		{"ExtendInt", 2, true, nil},
		{"ExtendString", "hello", true, nil},
		{"ExtendReflectValue", reflect.ValueOf(2), true, reflect.ValueOf(2)},
		{"ExtendNil", nil, true, nil},
		{"ExtendBool", true, true, false},
		{"ExtendStructWithAnyNotSet", TestExtendData{}, true, TestExtendData{}},
		{"ExtendStruct", testStruct1{Age: 20}, false, testStruct1{Age: 20}},
		{"ExtendStructWithEmbeddedField", testStructEmbedded{}, false, testStructEmbedded{}},
		{"ExtendStructWithUnexportedEmbeddedField", testStructUnexportedEmbedded{}, true, testStructUnexportedEmbedded{}},
	}

	assert := assert.New(t)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				r := recover()
				assert.Equal(test.shouldPanic, r != nil, r)
			}()
			dynamicStruct := dstruct.ExtendStruct(test.input).Build().Instance()

			dynamicStructConverted := dreflect.Convert(reflect.ValueOf(dynamicStruct), reflect.TypeOf(test.expected)).Interface()
			assert.EqualValues(test.expected, dynamicStructConverted)
		})
	}

}

type AddFieldData struct {
	name string
}

func TestAddField(t *testing.T) {

	b := dstruct.NewBuilder().
		AddField("Age", 20, `json:"Age"`)
	assert := assert.New(t)

	output := b.Build().Instance()

	assert.True(reflect.DeepEqual(output, struct {
		Age int `json:"Age"`
	}{20}))
}
