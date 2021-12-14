package errors

import "fmt"

type Error struct {
	kind *kind
	code int64
	err  error
}

func (e *Error) Kind() *kind {
	return e.kind
}

func (e *Error) IsKind(kind *kind) bool {
	return e.kind.Is(kind)
}

func (e *Error) Code() int64 {
	return e.code
}

func (e *Error) Error() string {
	return e.err.Error()
}

func NewError(kind *kind, codeAndErr ...interface{}) *Error {
	e := &Error{kind: kind}
	for _, param := range codeAndErr {
		switch p := param.(type) {
		case int:
			e.code = int64(p)
		case int64:
			e.code = p
		case string:
			e.err = fmt.Errorf(p)
		case error:
			e.err = p
		}
	}
	if e.code == 0 {
		e.code = e.kind.codeDef
	}
	if e.err == nil {
		e.err = fmt.Errorf(e.kind.errMsgDef)
	}
	return e
}
