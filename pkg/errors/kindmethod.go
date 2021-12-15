package errors

import "fmt"

func (k *kind) Is(kind *kind) bool {
	if k.name == kind.name {
		return true
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
