package ref

import "reflect"

type Entry struct {
	oriData interface{}
	oriVal  reflect.Value
	oriTyp  reflect.Type
	oriKind reflect.Kind

	data interface{}
	val  reflect.Value
	typ  reflect.Type
	kind reflect.Kind

	vMap           map[interface{}]interface{}
	vStructFields  map[string]*reflect.StructField
	vStructValues  map[string]interface{}
	vStructMethods map[string]*reflect.Method
}

func New(o interface{}) *Entry {
	e := new(Entry)
	e.oriData = o
	e.oriVal = reflect.ValueOf(e.oriData)
	e.oriTyp = e.oriVal.Type()
	e.oriKind = e.oriTyp.Kind()

	e.data = PeelOffPtr(e.oriData)
	e.val = reflect.ValueOf(e.data)
	e.typ = e.val.Type()
	e.kind = e.oriTyp.Kind()

	return e
}
