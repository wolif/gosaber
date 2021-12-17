package ref

import (
	"fmt"
	"reflect"
)

func (e *Entity) initStructMethods() bool {
	if e.vStructMethods != nil {
		return true
	}
	if !e.IsStruct() {
		return false
	}
	e.vStructMethods = make(map[string]*reflect.Method, e.GetType().NumMethod())
	for i := 0; i < e.GetType().NumMethod(); i++ {
		method := e.GetType().Method(i)
		e.vStructMethods[method.Name] = &method
	}
	return true
}

func (e *Entity) StructMethods() (methods map[string]*reflect.Method, ok bool) {
	if e.initStructMethods() {
		return e.vStructMethods, true
	}
	return nil, false
}

func (e *Entity) StructMethodGet(name string) (method *reflect.Method, ok bool) {
	if e.initStructMethods() {
		method, ok = e.vStructMethods[name]
		return
	}
	return nil, false
}

func (e *Entity) StructMethodCall(methodName string, args ...interface{}) (result []interface{}, err error) {
	method, ok := e.StructMethodGet(methodName)
	if !ok {
		return nil, fmt.Errorf("method with name %s not found", methodName)
	}
	if len(args)+1 != method.Type.NumIn() {
		err = fmt.Errorf("arguments num error, %d args expected, %d given", method.Type.NumIn()-1, len(args))
		return
	}
	params := []reflect.Value{e.GetOriValue()}
	for i := 1; i < method.Type.NumIn(); i++ {
		if method.Type.In(i) != reflect.TypeOf(args[i-1]) {
			err = fmt.Errorf("arg %d type error, %T need, %T given", i-1, method.Type.In(i), args[i-1])
			return
		}
		params = append(params, reflect.ValueOf(args[i-1]))
	}
	oriRes := method.Func.Call(params)
	result = make([]interface{}, len(oriRes))
	for i, r := range oriRes {
		result[i] = r.Interface()
	}
	return
}
