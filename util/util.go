package util

import (
	"fmt"
	"reflect"
)

func InArray(elem any, slice any) (bool, error) {
	val := reflect.ValueOf(slice)

	if val.Kind() != reflect.Slice {
		return false, fmt.Errorf("inArray() given a non-slice type")
	}
	for i := 0; i < val.Len(); i++ {
		// 比较切片中的元素和要查找的元素
		if val.Index(i).Interface() == elem {
			return true, nil
		}
	}
	return false, nil
}
