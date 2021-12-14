package ref

import "reflect"

func PeelOffPtr(i interface{}) interface{} {
	for reflect.TypeOf(i).Kind() == reflect.Ptr {
		i = reflect.ValueOf(i).Elem().Interface()
	}
	return i
}

func Is(o1, o2 interface{}) bool {
	return o1 == o2
}

func IsIn(o interface{}, scope ...interface{}) bool {
	for _, s := range scope {
		if o == s {
			return true
		}
	}
	return false
}

func IsInSlice(o interface{}, slice interface{}) bool {
	if reflect.TypeOf(slice).Kind() != reflect.Slice {
		return false
	}

	sliceVal := reflect.ValueOf(slice)
	sliceLen := sliceVal.Len()
	if sliceLen == 0 {
		return false
	}

	for i := 0; i < sliceLen; i++ {
		if o == sliceVal.Index(i).Interface() {
			return true
		}
	}
	return false
}
