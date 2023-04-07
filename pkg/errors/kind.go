package errors

import (
	"fmt"
	"sync"
)

type kind struct {
	super *kind
	name  string
	code  *int64
	err   *string
}

func (k *kind) Super() *kind {
	return k.super
}

func (k *kind) Name() string {
	return k.name
}

func (k *kind) Code() int64 {
	return *(k.code)
}

func (k *kind) Err() string {
	return *(k.err)
}

func (k *kind) IsAncestor() bool {
	return k.super == nil
}


var kindsRegistered = map[string]struct{}{}

func kindRegister(name string) string {
	if name == "" {
		panic("kind name can't empty")
	}
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	if _, ok := kindsRegistered[name]; ok {
		panic(fmt.Sprintf("kind name [%s] has already setted", name))
	}
	kindsRegistered[name] = struct{}{}
	return name
}

func NewBaseKind(name string, code int64, err string) *kind {
	return &kind{
		super: nil,
		name:  kindRegister(name),
		code:  &code,
		err:   &err,
	}
}

func (k *kind) Extend(name string, codeErr ...interface{}) *kind {
	ret := &kind{
		super: k,
		name:  kindRegister(name),
	}
	for _, param := range codeErr {
		switch p := param.(type) {
		case string:
			ret.err = &p
		case int:
			c := int64(p)
			ret.code = &c
		case int64:
			ret.code = &p
		}
	}
	if ret.code == nil {
		ret.code = k.code
	}
	if ret.err == nil {
		ret.err = k.err
	}
	return ret
}

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
