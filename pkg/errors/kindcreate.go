package errors

import (
	"fmt"
	"sync"
)

var kindsRegistered = map[string]struct{}{}

func kindRegister(name string) string {
	if name == "" {
		panic("kind name can't empty")
	}
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	if _, ok := kindsRegistered[name]; ok {
		panic(fmt.Sprintf("kind name [%s] has already setted", name))
	}
	kindsRegistered[name] = struct{}{}
	return name
}

func NewBaseKind(name string, code int64, err string) *kind {
	return &kind{
		super: nil,
		name:  kindRegister(name),
		code:  &code,
		err:   &err,
	}
}

func (k *kind) Extend(name string, codeErr ...interface{}) *kind {
	ret := &kind{
		super: k,
		name:  kindRegister(name),
	}
	for _, param := range codeErr {
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
