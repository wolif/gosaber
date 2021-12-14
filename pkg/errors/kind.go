package errors

import (
	"fmt"
	"strings"
)

type kind struct {
	superKind *kind
	name      string
	codeDef   int64
	errMsgDef string
}

var kindsSetted = map[string]struct{}{}

func kindNameCheck(kindName string) string {
	kindName = strings.TrimSpace(kindName)
	if kindName == "" {
		panic("kind name can't empty")
	}
	_, ok := kindsSetted[kindName]
	if ok {
		panic(fmt.Sprintf("kind name [%s] has already setted", kindName))
	}
	return kindName
}

func NewKind(name string, superKindCodeErrMsg ...interface{}) *kind {
	name = kindNameCheck(name)
	k := &kind{name: name}
	for _, param := range superKindCodeErrMsg {
		switch p := param.(type) {
		case string:
			k.errMsgDef = p
		case int:
			k.codeDef = int64(p)
		case int64:
			k.codeDef = p
		case *kind:
			k.superKind = p
		}
	}
	if k.superKind != nil {
		if k.codeDef == 0 {
			k.codeDef = k.superKind.codeDef
		}
		if k.errMsgDef == "" {
			k.errMsgDef = k.superKind.errMsgDef
		}
	}

	if k.codeDef == 0 {
		panic("kind code default hasn't been set")
	}
	if strings.TrimSpace(k.errMsgDef) == "" {
		panic("kind err message default hasn't been set")
	}
	return k
}

func (k *kind) Extend(name string, codeAndErrMsg ...interface{}) *kind {
	return NewKind(name, append(codeAndErrMsg, k)...)
}

func (k *kind) Is(kind *kind) bool {
	if k.name != kind.name && k.superKind == nil {
		return false
	}
	if k.name == kind.name {
		return true
	}
	return k.superKind.Is(kind)
}

func (k *kind) Error(codeAndErr ...interface{}) *Error {
	return NewError(k, codeAndErr...)
}

func (k *kind) Errorf(format string, a ...interface{}) *Error {
	return k.Error(fmt.Sprintf(format, a...))
}
