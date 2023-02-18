package dreflect

import (
	"fmt"
	"reflect"
	"unsafe"
)

func GetPointerToInterface(str any) any {
	return reflect.New(reflect.ValueOf(str).Type()).Interface()
}

func GetUnderlyingPointerValue(ptr any) any {
	return reflect.ValueOf(ptr).Elem().Interface()
}

func GetSliceType(value reflect.Value) reflect.Type {
	return reflect.TypeOf(value.Interface()).Elem()
}

func GetPointerToSliceType(sliceType reflect.Type) any {
	return reflect.New(sliceType).Elem().Interface()
}

// Convert extends the reflect.Convert function an proceeds to convert subtypes
func Convert(value reflect.Value, t reflect.Type) reflect.Value {
	defer func() {
		if recover() != nil {
			panic(fmt.Sprintf("dreflect.Convert: value of type %s cannot be converted to type %s", value.Type(), t))
		}
	}()
	dst := reflect.New(t).Elem()
	return convert(value, dst)
}

func convert(src reflect.Value, dst reflect.Value) reflect.Value {

	if src.Kind() != dst.Kind() {
		if src.Kind() < 1 || src.Kind() > 14 {
			panic(fmt.Sprintf("dreflect.Convert: value of type %s cannot be converted to type %s", src.Type(), dst.Type()))
		}
	}
	switch src.Kind() {
	case reflect.Struct:
		// var structFields []reflect.StructField
		newStruct := reflect.New(dst.Type()).Elem()
		pointerToStruct := reflect.New(src.Type())
		pointerToStruct.Elem().Set(src)

		for i := 0; i < newStruct.NumField(); i++ {
			f := pointerToStruct.Elem().Field(i)
			fieldValue := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
			newField := convert(fieldValue, dst.Field(i))

			reflect.NewAt(newField.Type(), unsafe.Pointer(newStruct.Field(i).UnsafeAddr())).
				Elem().
				Set(newField)

		}
		return newStruct
	case reflect.Slice:
		newSliceType, dstSliceType := getSliceArrayType(src, dst)
		newSlice := reflect.MakeSlice(reflect.SliceOf(newSliceType), src.Len(), src.Cap())

		for i := 0; i < src.Len(); i++ {
			newSlice.Index(i).Set(convert(src.Index(i), reflect.New(dstSliceType).Elem()))
		}
		return newSlice
	case reflect.Array:
		if src.Len() != dst.Len() {
			panic(fmt.Sprintf("dreflect.Convert: value of type %s cannot be converted to type %s", src.Type(), dst.Type()))
		}
		newArrayType, dstSliceType := getSliceArrayType(src, dst)
		newArray := reflect.New(reflect.ArrayOf(src.Len(), newArrayType)).Elem()
		for i := 0; i < src.Len(); i++ {
			newArray.Index(i).Set(convert(src.Index(i), reflect.New(dstSliceType).Elem()))
		}
		return newArray
	case reflect.Pointer:

		dstPointerType := dst.Type().Elem()
		if src.IsNil() {
			return reflect.Zero(dstPointerType)
		}

		srcPointerValue := reflect.ValueOf(GetUnderlyingPointerValue(src.Interface()))

		if src.Type().Elem().Kind() != dst.Type().Elem().Kind() && src.Elem().Kind() >= 1 && src.Elem().Kind() <= 14 {
			panic(fmt.Sprintf("dreflect.Convert: value of type %s cannot be converted to type %s", src.Type(), dst.Type()))
		}

		p := convert(srcPointerValue, reflect.New(dstPointerType).Elem())
		ps := reflect.New(dstPointerType)
		ps.Elem().Set(p)
		return ps

	}
	return src.Convert(dst.Type())
}

func getSliceArrayType(src reflect.Value, dst reflect.Value) (reflect.Type, reflect.Type) {
	srcSliceType := GetSliceType(src)
	dstSliceType := GetSliceType(dst)
	return convert(reflect.New(srcSliceType).Elem(), reflect.New(dstSliceType).Elem()).Type(), dstSliceType
}
