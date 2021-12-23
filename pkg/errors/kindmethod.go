package errors

import "fmt"

func (k *kind) Is(kind *kind, strict ...bool) bool {
	if k.name == kind.name {
		return true
	}
	if len(strict) > 0 && strict[0] {
		return false
	}
	if k.name != kind.name && k.superKind == nil {
		return false
	}
	return k.superKind.Is(kind)
}

func (k *kind) Error(codeAndErr ...interface{}) *Error {
	return NewError(k, codeAndErr...)
}

func (k *kind) Errorf(format string, a ...interface{}) *Error {
	return k.Error(fmt.Sprintf(format, a...))
}

func (k *kind) ErrorCodeF(code int64, format string, a ...interface{}) *Error {
	return NewError(k, code, fmt.Sprintf(format, a...))
}
