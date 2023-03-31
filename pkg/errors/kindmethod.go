package errors

import "fmt"

func (k *kind) Is(kind *kind, strict ...bool) bool {
	if k.name == kind.name {
		return true
	}
	if len(strict) > 0 && strict[0] {
		return false
	}
	if k.name != kind.name && k.super == nil {
		return false
	}
	return k.super.Is(kind)
}

func (k *kind) New(codeAndErr ...interface{}) *Error {
	return New(append(codeAndErr, k)...)
}

func (k *kind) NewF(format string, a ...interface{}) *Error {
	return k.New(fmt.Sprintf(format, a...))
}
