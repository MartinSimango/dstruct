package dreflect

import (
	"fmt"
	"reflect"
	"unsafe"
)

func GetPointerOfValueType(str any) any {
	ptr := reflect.New(reflect.ValueOf(str).Type())
	ptr.Elem().Set(reflect.ValueOf(str))
	return ptr.Interface()
}

func GetUnderlyingPointerValue(ptr any) any {
	return reflect.ValueOf(ptr).Elem().Interface()
}

func GetSliceType(slice any) reflect.Type {
	return reflect.TypeOf(slice).Elem()
}

func ConvertToType[T any](val any) T {
	return Convert(reflect.ValueOf(val), reflect.TypeOf(*new(T))).Interface().(T)
}

// Convert extends the reflect.Convert function an proceeds to convert subtypes
func Convert(value reflect.Value, t reflect.Type) reflect.Value {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Sprintf("dreflect.Convert: value of type %v cannot be converted to type %v", value.Type(), t))
		}
	}()
	dst := reflect.New(t).Elem()
	if value.Type().ConvertibleTo(t) {
		return value.Convert(t)
	}
	return convert(value, dst)
}

func convertibleTo(src, dst reflect.Type) bool {

	return !src.ConvertibleTo(dst) &&
		src.Kind() != reflect.Struct &&
		src.Kind() != reflect.Slice &&
		src.Kind() != reflect.Array &&
		src.Kind() != reflect.Pointer
}

func convert(src reflect.Value, dst reflect.Value) reflect.Value {

	if convertibleTo(src.Type(), dst.Type()) {
		panic(fmt.Sprintf("dreflect.Convert: value of type %s cannot be converted to type %s", src.Type(), dst.Type()))
	}

	switch src.Kind() {
	case reflect.Struct:
		newStruct := reflect.New(dst.Type()).Elem()
		srcStruct := reflect.ValueOf(GetPointerOfValueType(src.Interface())).Elem()
		dNum := newStruct.NumField()
		sNum := srcStruct.NumField()

		if dNum != sNum {
			panic(fmt.Sprintf("dreflect.Convert: Number of struct fields differ %s has %d subfield(s) and %s has %d subfield(s)", src.Type(), sNum, dst.Type(), dNum))

		}
		for i := 0; i < newStruct.NumField(); i++ {
			f := srcStruct.Field(i)
			if srcStruct.Type().Field(i).Name != newStruct.Type().Field(i).Name {
				panic(fmt.Sprintf("dreflect.Convert: value of type %s cannot be converted to type %s", src.Type(), dst.Type()))
			}
			fieldValue := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
			newField := convert(fieldValue, dst.Field(i))

			reflect.NewAt(newField.Type(), unsafe.Pointer(newStruct.Field(i).UnsafeAddr())).
				Elem().
				Set(newField)
		}
		return newStruct
	case reflect.Slice:

		if src.IsNil() {
			return reflect.Zero(dst.Type())
		}

		dstSliceType := dst.Type().Elem()

		newSliceType := getSliceArrayType(src, dstSliceType)
		newSlice := reflect.MakeSlice(reflect.SliceOf(newSliceType), src.Len(), src.Cap())

		for i := 0; i < src.Len(); i++ {
			newSlice.Index(i).Set(convert(src.Index(i), reflect.New(dstSliceType).Elem()))
		}
		return newSlice
	case reflect.Array:
		if src.Len() != dst.Len() {
			panic(fmt.Sprintf("dreflect.Convert: value of type %s cannot be converted to type %s", src.Type(), dst.Type()))
		}
		dstSliceType := GetSliceType(dst.Interface())

		newArrayType := getSliceArrayType(src, dstSliceType)
		newArray := reflect.New(reflect.ArrayOf(src.Len(), newArrayType)).Elem()
		for i := 0; i < src.Len(); i++ {
			newArray.Index(i).Set(convert(src.Index(i), reflect.New(dstSliceType).Elem()))
		}
		return newArray
	case reflect.Pointer:

		if src.IsNil() {
			return reflect.Zero(dst.Type())
		}
		dstPointerValueType := dst.Type().Elem()
		srcPointerValue := reflect.ValueOf(GetUnderlyingPointerValue(src.Interface()))

		if src.Type().Elem().Kind() != dst.Type().Elem().Kind() && src.Elem().Kind() >= 1 && src.Elem().Kind() <= 14 {
			panic(fmt.Sprintf("dreflect.Convert: value of type %s cannot be converted to type %s", src.Type(), dst.Type()))
		}

		retPointer := reflect.New(dst.Type())
		//ensure that new pointer uses same memory address as src pointer
		reflect.NewAt(src.Type(), unsafe.Pointer(retPointer.Elem().UnsafeAddr())).Elem().
			Set(src)

		// now copy over the values from the src
		retPointer.Elem().Elem().Set(convert(srcPointerValue, reflect.New(dstPointerValueType).Elem()))

		return retPointer.Elem()

	}
	return src.Convert(dst.Type())
}

func getSliceArrayType(src reflect.Value, dstSliceType reflect.Type) reflect.Type {
	srcSliceType := GetSliceType(src.Interface())
	return convert(reflect.New(srcSliceType).Elem(), reflect.New(dstSliceType).Elem()).Type()
}
