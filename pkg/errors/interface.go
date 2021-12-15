package errors

type Err interface {
	Kind() *kind
	IsKind(*kind) bool
	Code() int64
	Error() string
}
