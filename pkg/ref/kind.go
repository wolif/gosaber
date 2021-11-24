package ref

import "reflect"

func (e *entry) IsInt() bool {
	return IsInScope(e.refType.kind, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64)
}
func (e *entry) IsUint() bool {
	return IsInScope(e.refType.kind, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64)
}
func (e *entry) IsFloat() bool {
	return IsInScope(e.refType.kind, reflect.Float32, reflect.Float64)
}
func (e *entry) IsNumber() bool {
	return e.IsInt() || e.IsUint() || e.IsFloat()
}
func (e *entry) IsString() bool {
	return e.refType.kind == reflect.String
}
func (e *entry) IsStringOrNumber() bool {
	return e.IsString() || e.IsNumber()
}

func (e *entry) IsPtr() bool {
	return e.refType.originKind == reflect.Ptr
}
func (e *entry) IsSlice() bool {
	return e.refType.kind == reflect.Slice
}
func (e *entry) IsMap() bool {
	return e.refType.kind == reflect.Map
}
func (e *entry) IsStruct() bool {
	return e.refType.kind == reflect.Struct
}
func (e *entry) IsFunc() bool {
	return e.refType.kind == reflect.Func
}
