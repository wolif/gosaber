package errors

type Err interface {
	Kind() *kind
	IsA(k *kind, strict ...bool) bool
	Code() int64
	Error() string
}
