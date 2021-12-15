package errors1

type Kind struct {
	Name    string
	DefCode int64
	DefErr  string
	defType *Type
}

func NewKind(name string, defCode int64, defErr string) *Kind {
	k := &Kind{
		Name:    name,
		DefCode: defCode,
		DefErr:  defErr,
	}
	k.defType = k.NewType(name, defCode, defErr)
	return k
}

func (k *Kind) NewType(name string, defCode int64, defErr string) *Type {
	return &Type{
		Kind:    k,
		Name:    name,
		DefCode: defCode,
		DefErr:  defErr,
	}
}

func (k *Kind) NewError(codeAndErr ...interface{}) *Error {
	return NewError(k.defType, codeAndErr...)
}
