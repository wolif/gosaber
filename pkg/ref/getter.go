package ref

import "reflect"

func (e *Entity) GetOriData() interface{} {
	return e.oriData
}

func (e *Entity) GetOriType() reflect.Type {
	return e.oriTyp
}

func (e *Entity) GetOriValue() reflect.Value {
	return e.oriVal
}

func (e *Entity) GetOriKind() reflect.Kind {
	return e.oriKind
}

func (e *Entity) GetData() interface{} {
	return e.data
}

func (e *Entity) GetType() reflect.Type {
	return e.typ
}

func (e *Entity) GetValue() reflect.Value {
	return e.val
}

func (e *Entity) GetKind() reflect.Kind {
	return e.kind
}
