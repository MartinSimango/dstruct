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

	type TestStructCopy struct {
		a int
		b *float32
	}
	var sliceInt []int
	var sliceInt8 []int8

	var tests = []ConvertTestData{
		{"ConvertIntToInt8", int(2), reflect.TypeOf(int8(0)), false, int8(2)},
		{"ConvertIntPointerToInt8Pointer", new(int), reflect.TypeOf(new(int8)), true, nil},
		{"ConvertIntToString", int(65), reflect.TypeOf("2"), false, "A"},
		{"ConvertFloatToString", float64(6), reflect.TypeOf("2"), true, nil},
		{"ConvertIntToIntPointer", int(2), reflect.TypeOf(new(int)), true, nil},
		{"ConvertIntSliceNilToInt8Slice", sliceInt, reflect.TypeOf(sliceInt8), false, sliceInt8},
		{"ConvertIntSliceToInt8Slice", []int{1, 2, 128}, reflect.TypeOf(sliceInt8), false, []int8{1, 2, -128}},
		{"ConvertIntArrayToInt8Slice", [1]int{1}, reflect.TypeOf([1]int8{}), false, [1]int8{1}},
		{"ConvertIntArrayToInt8DifferentLengths", [2]int{1}, reflect.TypeOf([1]int8{}), true, nil},
		{"ConvertPointerToStruct", &TestStruct{a: 2}, reflect.TypeOf(&TestStructCopy{}), false, &TestStructCopy{a: 2}},

		// TODO Mock dstruct.ExtendStruct
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

type TestGetPointerToInterfaceData struct {
	name        string
	input       any
	shouldPanic bool
	expected    any
}

func newValue[T any](val T) *T {
	v := new(T)
	*v = val
	return v
}

func TestGetPointerToInterface(t *testing.T) {
	assert := assert.New(t)

	var tests = []TestGetPointerToInterfaceData{
		{"GetPointerToBasicType", 2, false, new(int)},
		{"GetPointerToNil", nil, true, nil},
		{"GetPointerToStruct", struct{}{}, false, &struct{}{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				r := recover()
				assert.Equal(test.shouldPanic, r != nil, r)
			}()
			assert.EqualValues(test.expected, dreflect.GetPointerToInterface(test.input))
		})
	}
}

type TestGetUnderlyingPointerValueData struct {
	name        string
	input       any
	shouldPanic bool
	expected    any
}

func TestGetUnderlyingPointerValue(t *testing.T) {
	assert := assert.New(t)
	var pointer *int = new(int)
	*pointer = 2
	var pointerToPointer **int = &pointer
	var tests = []TestGetUnderlyingPointerValueData{
		{"GetUnderlyingPointerValueOfPointer", pointer, false, 2},
		{"GetUnderlyingPointerValueOfPointerToPointer", pointerToPointer, false, pointer},
		{"GetUnderlyingPointerValueOfNonPointer", 2, true, nil},
		{"GetUnderlyingPointerValueOfNilPointer", nil, true, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				r := recover()
				assert.Equal(test.shouldPanic, r != nil, r)
			}()
			assert.EqualValues(test.expected, dreflect.GetUnderlyingPointerValue(test.input))
		})
	}
}

type TestGetSliceTypeData struct {
	name        string
	input       reflect.Value
	shouldPanic bool
	expected    reflect.Type
}

func TestGetSliceType(t *testing.T) {
	var a []int
	dreflect.GetSliceType(reflect.ValueOf(a))

	assert := assert.New(t)

	var tests = []TestGetSliceTypeData{
		{"GetSliceTypeOfIntSlice", reflect.ValueOf([]int{}), false, reflect.TypeOf(2)},
		{"GetSliceTypeOfNonSliceType", reflect.ValueOf(2), true, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				r := recover()
				assert.Equal(test.shouldPanic, r != nil, r)
			}()
			assert.EqualValues(test.expected, dreflect.GetSliceType(test.input))
		})
	}
}
