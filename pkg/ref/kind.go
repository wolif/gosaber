package ref

import "reflect"

func (e *Entity)  IsKind(kinds ...reflect.Kind) bool {
	for _, k := range kinds {
		if k == e.kind {
			return true
		}
	}
	return false
}

func (e *Entity) IsOriKind(kinds ...reflect.Kind) bool {
	for _, k := range kinds {
		if k == e.oriKind {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------------------------------------------------
func (e *Entity) IsInt() bool {
	return e.IsKind(reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64)
}

func (e *Entity) IsUint() bool {
	return e.IsKind(reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64)
}

func (e *Entity) IsFloat() bool {
	return e.IsKind(reflect.Float32, reflect.Float64)
}

func (e *Entity) IsNumber() bool {
	return e.IsInt() || e.IsUint() || e.IsFloat()
}

func (e *Entity) IsString() bool {
	return e.IsKind(reflect.String)
}

func (e *Entity) IsStringOrNumber() bool {
	return e.IsString() || e.IsNumber()
}

func (e *Entity) IsComplex() bool {
	return e.IsKind(reflect.Complex64, reflect.Complex128)
}

func (e *Entity) IsBool() bool {
	return e.IsKind(reflect.Bool)
}

// ---------------------------------------------------------------------------------------------------------------------
func (e *Entity) IsSlice() bool {
	return e.IsKind(reflect.Slice)
}

func (e *Entity) IsMap() bool {
	return e.IsKind(reflect.Map)
}

func (e *Entity) IsStruct() bool {
	return e.IsKind(reflect.Struct)
}

func (e *Entity) IsFunc() bool {
	return e.IsKind(reflect.Func)
}

func (e *Entity) IsArray() bool {
	return e.IsKind(reflect.Array)
}

func (e *Entity) IsChan() bool {
	return e.IsKind(reflect.Chan)
}

// ---------------------------------------------------------------------------------------------------------------------
func (e *Entity) IsInterface() bool {
	return e.IsOriKind(reflect.Interface)
}

func (e *Entity) IsPtr() bool {
	return e.IsOriKind(reflect.Ptr)
}

func (e *Entity) IsUintPtr() bool {
	return e.IsOriKind(reflect.Uintptr)
}

func (e *Entity) IsUnsafePointer() bool {
	return e.IsOriKind(reflect.UnsafePointer)
}
