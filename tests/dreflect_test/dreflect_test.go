package dreflect_test

import (
	"reflect"
	"testing"

	"github.com/MartinSimango/dstruct"
	"github.com/MartinSimango/dstruct/dreflect"
	"github.com/stretchr/testify/assert"
)

type ConvertTestData struct {
	name        string
	input       any
	outputType  reflect.Type
	shouldPanic bool
	expected    any
}

func TestConvert(t *testing.T) {
	type TestStruct struct {
		a int
		b *float32
	}

	var tests = []ConvertTestData{
		{"ConvertIntToInt", int(2), reflect.TypeOf(int(0)), false, int(2)},
		{"ConvertIntToInt*", int(2), reflect.TypeOf(new(int)), true, int(2)},
		{"ConvertDStructToGoStruct", dstruct.ExtendStruct(TestStruct{a: 2}).Build().Instance(),
			reflect.TypeOf(TestStruct{}), false, TestStruct{a: 2}},
		{"ConvertGoStructToDStruct", TestStruct{a: 2},
			reflect.TypeOf(TestStruct{}), false, dstruct.ExtendStruct(TestStruct{a: 2}).Build().Instance()},
	}

	assert := assert.New(t)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				r := recover()
				assert.Equal(test.shouldPanic, r != nil, r)
			}()
			assert.EqualValues(test.expected, dreflect.Convert(reflect.ValueOf(test.input), test.outputType).Interface())
		})
	}

}

func TestGetPointerToInterface(t *testing.T) {
	var b int
	dreflect.GetPointerToInterface(b)
	// return reflect.New(reflect.ValueOf(str).Type()).Interface()
}

func TestGetUnderlyingPointerValue(t *testing.T) {
	var b *int = new(int)

	dreflect.GetUnderlyingPointerValue(b)
}

func TestGetSliceType(t *testing.T) {
	var a []int
	dreflect.GetSliceType(reflect.ValueOf(a))
}

func TestGetPointerToSliceType(t *testing.T) {
	dreflect.GetPointerToSliceType(reflect.TypeOf(2))
}

// Convert e
