package errors

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
