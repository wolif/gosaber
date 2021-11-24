package ref

import "reflect"

func PeelOffPtr(i interface{}) interface{} {
	for reflect.TypeOf(i).Kind() == reflect.Ptr {
		i = reflect.ValueOf(i).Elem().Interface()
	}
	return i
}

func IsInScope(i interface{}, scope ...interface{}) bool {
	for _, s := range scope {
		if i == s {
			return true
		}
	}
	return false
}

func IsInSlice(i interface{}, slice interface{}) bool {
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		return false
	}

	sliceVal := reflect.ValueOf(slice)
	sliceLen := sliceVal.Len()
	if sliceLen == 0 {
		return false
	}

	scope := make([]interface{}, sliceLen, sliceLen)
	for i := 0; i < sliceLen; i++ {
		scope[i] = sliceVal.Index(i).Interface()
	}

	return IsInScope(i, scope...)
}