package ref

import "reflect"

/*
*功能和 reflectutils 功能类似, 部分方法名称有差异
*不同点在于: 反射不同的实体时,如果反射出的reflect.Type相同,就会使用同一份类型信息
*这样可以节省程序的资源消耗
 */
type entry struct {
	refType *refType
	// 原始数据相关
	originData interface{}
	originVal  reflect.Value
	// 去指针后数据相关
	data  interface{}
	value reflect.Value
	// map相关
	mapKeys    []interface{}
	fieldValue map[string]interface{}
}

func New(i interface{}) *entry {
	rv := new(entry)
	rv.refType = newRefType(i)
	rv.originData = i
	rv.originVal = reflect.ValueOf(i)
	rv.data = PeelOffPtr(i)
	rv.value = reflect.ValueOf(rv.data)

	if rv.IsMap() {
		rv.mapKeys = make([]interface{}, 0)
		for _, k := range rv.value.MapKeys() {
			rv.mapKeys = append(rv.mapKeys, k.Interface())
		}
	} else if rv.IsStruct() {
		rv.fieldValue = make(map[string]interface{})
		for fieldName, _ := range rv.refType.structField {
			rv.fieldValue[fieldName] = rv.value.FieldByName(fieldName).Interface()
		}
	}

	return rv
}
