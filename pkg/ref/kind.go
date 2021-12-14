package ref

import "reflect"

func (e *Entry)  IsKind(kind reflect.Kind) bool {
	return e.kind == kind
}

func (e *Entry) IsKindIn(kinds ...reflect.Kind) bool {
	for _, k := range kinds {
		if k == e.kind {
			return true
		}
	}
	return false
}

func (e *Entry) IsOriKind(kind reflect.Kind) bool {
	return e.oriKind == kind
}

func (e *Entry) IsOriKindIn(kinds ...reflect.Kind) bool {
	for _, k := range kinds {
		if k == e.oriKind {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------------------------------------------------
func (e *Entry) IsInt() bool {
	return e.IsKindIn(reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64)
}

func (e *Entry) IsUint() bool {
	return e.IsKindIn(reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64)
}

func (e *Entry) IsFloat() bool {
	return e.IsKindIn(reflect.Float32, reflect.Float64)
}

func (e *Entry) IsNumber() bool {
	return e.IsInt() || e.IsUint() || e.IsFloat()
}

func (e *Entry) IsString() bool {
	return e.IsKindIn(reflect.String)
}

func (e *Entry) IsStringOrNumber() bool {
	return e.IsString() || e.IsNumber()
}

func (e *Entry) IsComplex() bool {
	return e.IsKindIn(reflect.Complex64, reflect.Complex128)
}

func (e *Entry) IsBool() bool {
	return e.IsKind(reflect.Bool)
}

// ---------------------------------------------------------------------------------------------------------------------
func (e *Entry) IsSlice() bool {
	return e.IsKind(reflect.Slice)
}

func (e *Entry) IsMap() bool {
	return e.IsKind(reflect.Map)
}

func (e *Entry) IsStruct() bool {
	return e.IsKind(reflect.Struct)
}

func (e *Entry) IsFunc() bool {
	return e.IsKind(reflect.Func)
}

func (e *Entry) IsArray() bool {
	return e.IsKind(reflect.Array)
}

func (e *Entry) IsChan() bool {
	return e.IsKind(reflect.Chan)
}

// ---------------------------------------------------------------------------------------------------------------------
func (e *Entry) IsInterface() bool {
	return e.IsOriKind(reflect.Interface)
}

func (e *Entry) IsPtr() bool {
	return e.IsOriKind(reflect.Ptr)
}

func (e *Entry) IsUintPtr() bool {
	return e.IsOriKind(reflect.Uintptr)
}

func (e *Entry) IsUnsafePointer() bool {
	return e.IsOriKind(reflect.UnsafePointer)
}
