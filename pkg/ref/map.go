package ref

import "reflect"

func (e *entry) HasMapKey(key interface{}) bool {
	if !e.IsMap() {
		return false
	}
	return IsInScope(key, e.mapKeys...)
}

func (e *entry) GetMapValue(key interface{}) (value interface{}, ok bool) {
	if e.HasMapKey(key) {
		return e.value.MapIndex(reflect.ValueOf(key)).Interface(), true
	}
	return nil, false
}