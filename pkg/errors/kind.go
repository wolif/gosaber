package errors

type kind struct {
	superKind *kind
	name      string
	code      *int64
	err       *string
}

func (k *kind) SuperKind() *kind {
	return k.superKind
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

func (k *kind) IsBaseKind() bool {
	return k.superKind == nil
}
