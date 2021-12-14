package errors

type Type struct {
	Kind    *Kind
	Name    string
	DefCode int64
	DefErr  string
}

func (t *Type) IsKind(kind *Kind) bool {
	return t.Kind == kind
}

func (t *Type) NewError(codeAndErr ...interface{}) *Error {
	return NewError(t, codeAndErr...)
}
