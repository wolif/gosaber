package errors1

import "fmt"

type Error struct {
	kind *Kind
	typ  *Type
	code int64
	err  error
}

func (e *Error) Kind() *Kind {
	return e.kind
}

func (e *Error) IsKind(kind *Kind) bool {
	return e.kind.Name == kind.Name
}

func (e *Error) Type() *Type {
	return e.typ
}

func (e *Error) IsType(typ *Type) bool {
	return e.typ.Name == typ.Name
}

func (e *Error) Code() int64 {
	return e.code
}

func (e *Error) Error() string {
	return e.err.Error()
}

func NewError(typ *Type, codeAndErr ...interface{}) *Error {
	e := &Error{
		kind: typ.Kind,
		typ:  typ,
		code: typ.DefCode,
		err:  fmt.Errorf(typ.DefErr),
	}
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
	return e
}
