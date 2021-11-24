package ref

import "reflect"

var refTypeCache = map[reflect.Type]*refType{}

type structFieldName = string
type structFieldTagName = string
type structMethodName = string

type refType struct {
	// 原始数据相关
	originTyp  reflect.Type
	originKind reflect.Kind
	// 去指针后数据相关
	typ  reflect.Type
	kind reflect.Kind
	// 结构体相关
	structMethod   map[structMethodName]*reflect.Method
	structField    map[structFieldName]*reflect.StructField
	structFieldTag map[structFieldName]map[structFieldTagName]string
}

func newRefType(i interface{}) *refType {
	t, found := refTypeCache[reflect.TypeOf(i)]
	if found {
		return t
	}

	t = new(refType)
	t.originTyp = reflect.TypeOf(i)
	t.originKind = t.originTyp.Kind()
	t.typ = reflect.TypeOf(PeelOffPtr(i))
	t.kind = t.typ.Kind()

	if t.kind == reflect.Struct {
		t.structMethod = make(map[structMethodName]*reflect.Method)
		t.structField = make(map[structFieldName]*reflect.StructField)
		t.structFieldTag = make(map[structFieldName]map[structFieldTagName]string)
		for i := 0; i < t.typ.NumField(); i++ {
			field := t.typ.Field(i)
			t.structField[field.Name] = &field
		}
		for i := 0; i < t.originTyp.NumMethod(); i++ {
			method := t.originTyp.Method(i)
			t.structMethod[method.Name] = &method
		}
	}

	refTypeCache[t.originTyp] = t
	return t
}
