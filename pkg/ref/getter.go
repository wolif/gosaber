package ref

import "reflect"

func (e *Entry) GetOriData() interface{} {
	return e.oriData
}

func (e *Entry) GetOriType() reflect.Type {
	return e.oriTyp
}

func (e *Entry) GetOriValue() reflect.Value {
	return e.oriVal
}

func (e *Entry) GetOriKind() reflect.Kind {
	return e.oriKind
}

func (e *Entry) GetData() interface{} {
	return e.data
}

func (e *Entry) GetType() reflect.Type {
	return e.typ
}

func (e *Entry) GetValue() reflect.Value {
	return e.val
}

func (e *Entry) GetKind() reflect.Kind {
	return e.kind
}
