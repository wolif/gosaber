package ref

import "reflect"

func (e *entry) GetOriginType() reflect.Type {
	return e.refType.originTyp
}

func (e *entry) GetOriginKind() reflect.Kind {
	return e.refType.originKind
}

func (e *entry) GetOriginVal() reflect.Value {
	return e.originVal
}

func (e *entry) IsOriginZero() bool {
	return e.originVal.IsZero()
}

func (e *entry) GetType() reflect.Type {
	return e.refType.typ
}

func (e *entry) GetKind() reflect.Kind {
	return e.refType.kind
}

func (e *entry) GetValue() reflect.Value {
	return e.value
}

func (e *entry) IsZero() bool {
	return e.value.IsZero()
}

func (e *entry) IsNil() bool {
	return e.value.IsNil()
}

func (e *entry) IsValid() bool {
	return e.value.IsValid()
}
