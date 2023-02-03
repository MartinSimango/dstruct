package dstruct

import (
	"fmt"
	"reflect"
)

// MergeStructs merges two structs
func MergeStructs(strcts ...interface{}) (a any, err error) {
	if len(strcts) < 2 {
		return nil, fmt.Errorf("failed to merge: number of structs to merge must be 2 or more")
	}

	for i := 0; i < len(strcts); i++ {
		if reflect.ValueOf(strcts[i]).Kind() != reflect.Struct {
			return nil, fmt.Errorf("failed to merged structs: %d interface is not a struct ", (i + 1))
		}
	}

	left := ExtendStruct(strcts[0])
	var right Builder
	for i := 1; i < len(strcts); i++ {
		right = ExtendStruct(strcts[i])
		mergedStruct, err := mergeStructs(left, right, reflect.Struct)
		if err != nil {
			return nil, err
		}
		left = ExtendStruct(mergedStruct)

	}

	return left.Build().Instance(), nil
}

// TODO clean this function up
func mergeStructs(left, right Builder, parentKind reflect.Kind) (any, error) {

	// struct to be returned
	newStruct := ExtendStruct(left.Build().Instance())

	for name, field := range right.(*treeBuilderImpl).root.children {

		elementName := field.data.Name
		cV := left.(*treeBuilderImpl).root.GetNode(name)

		if cV == nil {
			newStruct.AddField(elementName, field.data.Value.Interface(), string(field.data.Tag))
			continue
		}
		if err := validateTypes(field.data.Value, cV.data.Value, field.data.fqn); err != nil {
			return nil, err
		}

		if field.data.Value.Kind() == reflect.Slice {
			vSliceType := getSliceType(field.data.Value)
			cVSliceType := getSliceType(cV.data.Value)
			if err := validateSliceTypes(vSliceType, cVSliceType, field.data.Value, cV.data.Value, field.data.fqn); err != nil {
				return nil, err
			}
			newStruct.RemoveField(field.data.fqn)
			if cVSliceType.Kind() == reflect.Struct {
				newSliceTypeStruct, err := mergeStructs(left.GetField(name),
					right.GetField(name), reflect.Slice)

				if err != nil {
					return nil, err
				}
				newStruct.AddField(field.data.Name, newSliceTypeStruct, "")
			} else {
				newStruct.AddField(field.data.Name, field.data.Value.Interface(), "")
			}

		} else if field.data.Value.Kind() == reflect.Struct {
			updatedSchema, err := mergeStructs(left.GetField(name), right.GetField(name), reflect.Struct)
			if err != nil {
				return nil, err
			}
			newStruct.RemoveField(field.data.GetFieldName())
			newStruct.AddField(field.data.Name, updatedSchema, string(field.data.Tag))
		}

	}

	if parentKind == reflect.Slice {
		sliceOfElementType := reflect.SliceOf(reflect.ValueOf(newStruct.Build().Instance()).Elem().Type())
		return reflect.MakeSlice(sliceOfElementType, 0, 1024).Interface(), nil
	}
	return newStruct.Build().Instance(), nil
}

func shouldTypeMatch(kind reflect.Kind) bool {
	if kind == reflect.Array || kind == reflect.Struct || kind == reflect.Slice {
		return false
	}
	return true
}

func validateTypes(v, cV reflect.Value, fullFieldName string) error {
	currentElementType := reflect.TypeOf(cV.Interface())
	newElementType := reflect.TypeOf(v.Interface())
	if shouldTypeMatch(v.Kind()) || shouldTypeMatch(cV.Kind()) {
		if currentElementType != newElementType {
			return fmt.Errorf("mismatching types for field '%s': %s and %s", fullFieldName, currentElementType, newElementType)
		}
	} else {
		if v.Kind() != cV.Kind() {
			return fmt.Errorf("mismatching types for field '%s': %s and %s", fullFieldName, currentElementType, newElementType)
		}
	}
	return nil
}

func validateSliceTypes(vSliceType, cVSliceType reflect.Type, v, cV reflect.Value, fullFieldName string) error {
	currentElementType := reflect.TypeOf(reflect.New(cVSliceType).Interface())
	newElementType := reflect.TypeOf(reflect.New(vSliceType).Interface())

	if shouldTypeMatch(vSliceType.Kind()) || shouldTypeMatch(cVSliceType.Kind()) {
		if currentElementType != newElementType {
			return fmt.Errorf("mismatching types for field '%s': %s and %s", fullFieldName, reflect.TypeOf(v.Interface()), reflect.TypeOf(cV.Interface()))
		}
	} else {
		if v.Kind() != cV.Kind() {
			return fmt.Errorf("mismatching types for field '%s': %s and %s", fullFieldName, reflect.TypeOf(v.Interface()), reflect.TypeOf(cV.Interface()))
		}
	}
	return nil
}
