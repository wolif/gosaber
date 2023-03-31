package errors

import (
	"fmt"
)

type Error struct {
	kind *kind
	code *int64
	err  error
}

func (e Error) Kind() *kind {
	return e.kind
}

func (e Error) IsA(kind *kind, strict ...bool) bool {
	return e.kind.Is(kind, strict...)
}

func (e Error) Code() int64 {
	return *(e.code)
}

func (e Error) Error() string {
	return e.err.Error()
}

var DefKind = NewBaseKind("default error Kind", -1, "")

func New(kindCodeErr ...interface{}) *Error {
	e := &Error{}
	for _, param := range kindCodeErr {
		switch p := param.(type) {
		case int:
			c := int64(p)
			e.code = &c
		case int64:
			e.code = &p
		case string:
			e.err = fmt.Errorf(p)
		case error:
			e.err = fmt.Errorf(p.Error())
		case *kind:
			e.kind = p
		}
	}
	if e.kind == nil {
		e.kind = DefKind
	}
	if e.code == nil {
		e.code = e.kind.code
	}
	if e.err == nil {
		e.err = fmt.Errorf(*(e.kind.err))
	}
	return e
}
