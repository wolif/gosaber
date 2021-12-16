package ref

import "reflect"

func (e *Entity) initStructFields() bool {
	if e.vStructFields != nil {
		return true
	}
	if !e.IsStruct() {
		return false
	}
	e.vStructFields = make(map[string]*reflect.StructField, e.GetType().NumField())
	for i := 0; i < e.GetType().NumField(); i++ {
		field := e.GetType().Field(i)
		e.vStructFields[field.Name] = &field
	}
	return true
}

func (e *Entity) StructFields() (map[string]*reflect.StructField, bool) {
	if !e.initStructFields() {
		return nil, false
	}
	return e.vStructFields, true
}

func (e *Entity) StructFieldsName() ([]string, bool) {
	if !e.initStructFields() {
		return nil, false
	}
	ret := make([]string, 0, len(e.vStructFields))
	for name, _ := range e.vStructFields {
		ret = append(ret, name)
	}
	return ret, true
}

func (e *Entity) StructFieldGet(name string) (field *reflect.StructField, ok bool) {
	if e.initStructFields() {
		field, ok = e.vStructFields[name]
		return
	}
	return nil, false
}
