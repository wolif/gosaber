package errors

import (
	"fmt"
	"sync"
)

var kindsRegistered = map[string]struct{}{}

func kindRegister(kindName string) string {
	if kindName == "" {
		panic("kind name can't empty")
	}
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	if _, ok := kindsRegistered[kindName]; ok {
		panic(fmt.Sprintf("kind name [%s] has already setted", kindName))
	}
	kindsRegistered[kindName] = struct{}{}
	return kindName
}

func NewBaseKind(name string, code int64, err string) *kind {
	return &kind{
		superKind: nil,
		name:      kindRegister(name),
		code:      &code,
		err:       &err,
	}
}

func (k *kind) Extend(name string, codeAndErr ...interface{}) *kind {
	ret := &kind{
		superKind: k,
		name:      kindRegister(name),
	}
	for _, param := range codeAndErr {
		switch p := param.(type) {
		case string:
			ret.err = &p
		case int:
			c := int64(p)
			ret.code = &c
		case int64:
			ret.code = &p
		}
	}
	if ret.code == nil {
		ret.code = k.code
	}
	if ret.err == nil {
		ret.err = k.err
	}
	return ret
}
