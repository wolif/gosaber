package errors

type Err interface {
	Kind() *Kind
	Type() *Type
	Code() int64
	Error() string
	IsKind(kind *Kind) bool
	IsType(typ *Type) bool
}
