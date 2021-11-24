package ref

import "reflect"

func (e *entry) HasStructField(fieldName string) bool {
	if !e.IsStruct() {
		return false
	}
	_, found := e.refType.structField[fieldName]
	return found
}

func (e *entry) GetStructField(fieldName string) (structField *reflect.StructField, ok bool) {
	if e.HasStructField(fieldName) {
		return e.refType.structField[fieldName], true
	}
	return nil, false
}

func (e *entry) GetStructFields() (fields map[structFieldName]*reflect.StructField, ok bool) {
	if !e.IsStruct() {
		return nil, false
	}
	return e.refType.structField, true
}

func (e *entry) GetStructFieldValue(fieldName string) (value interface{}, ok bool) {
	if e.HasStructField(fieldName) {
		return e.value.FieldByName(fieldName).Interface(), true
	}
	return nil, false
}

func (e *entry) HasStructFieldTag(fieldName string, tagName string) bool {
	if !e.IsStruct() {
		return false
	}
	if !e.HasStructField(fieldName) {
		return false
	}

	if _, found := e.refType.structFieldTag[fieldName]; !found {
		e.refType.structFieldTag[fieldName] = make(map[structFieldTagName]string)
	}

	if _, found := e.refType.structFieldTag[fieldName][tagName]; found {
		return true
	}

	tag, ok := e.refType.structField[fieldName].Tag.Lookup(tagName)
	if ok {
		e.refType.structFieldTag[fieldName][tagName] = tag
		return true
	}

	return false
}

func (e *entry) GetStructFieldTag(fieldName, tagName string) (value string, ok bool) {
	if e.HasStructFieldTag(fieldName, tagName) {
		return e.refType.structFieldTag[fieldName][tagName], true
	}
	return "", false
}

func (e *entry) HasStructMethod(methodName string) bool {
	if !e.IsStruct() {
		return false
	}
	_, found := e.refType.structMethod[methodName]
	return found
}

func (e *entry) GetStructMethod(methodName string) (method *reflect.Method, ok bool) {
	if e.HasStructMethod(methodName) {
		return e.refType.structMethod[methodName], true
	}
	return nil, false
}

func (e *entry) CallStructMethod(methodName string, params ...interface{}) (retVal []interface{}, ok bool) {
	structMethod, ok := e.GetStructMethod(methodName)
	if !ok {
		return
	}

	funcParams := make([]reflect.Value, len(params)+1, len(params)+1)
	funcParams[0] = reflect.ValueOf(e.originData)
	for i, v := range params {
		funcParams[i+1] = reflect.ValueOf(v)
	}
	retVal = make([]interface{}, 0)
	retValOrigin := structMethod.Func.Call(funcParams)
	for _, ro := range retValOrigin {
		retVal = append(retVal, ro.Interface())
	}
	return retVal, true
}

func (e *entry) CallStructMethodSlice(methodName string, paramsSlice interface{}) (retVal []interface{}, ok bool) {
	if reflect.TypeOf(paramsSlice).Kind() != reflect.Slice {
		return
	}
	valParams := reflect.ValueOf(paramsSlice)
	numParams := valParams.Len()
	params := make([]interface{}, numParams, numParams)
	for i := 0; i < numParams; i++ {
		params[i] = valParams.Index(i).Interface()
	}

	return e.CallStructMethod(methodName, params...)
}
