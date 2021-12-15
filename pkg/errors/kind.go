package errors

import (
	"fmt"
	"strings"
)

type kind struct {
	superKind *kind
	name      string
	codeDef   *int64
	errMsgDef *string
}

func (k *kind) GetDefaultCode() int64 {
	return *(k.codeDef)
}

func (k *kind) SetDefaultCode(code int64) *kind {
	k.codeDef = &code
	return k
}

func (k *kind) GetDefaultErrMsg() string {
	return *(k.errMsgDef)
}

func (k *kind) SetDefaultErrMsg(message string) *kind {
	k.errMsgDef = &message
	return k
}

var kindsRegistered = map[string]struct{}{}

func kindNameCheck(kindName string) string {
	kindName = strings.TrimSpace(kindName)
	if kindName == "" {
		panic("kind name can't empty")
	}
	if _, ok := kindsRegistered[kindName]; ok {
		panic(fmt.Sprintf("kind name [%s] has already setted", kindName))
	}
	return kindName
}

func NewBaseKind(name string, codeDef int64, errMsgDef string) *kind {
	return &kind{
		superKind: nil,
		name:      kindNameCheck(name),
		codeDef:   &codeDef,
		errMsgDef: &errMsgDef,
	}
}

func (k *kind) Extend(name string, codeAndErrMsgDef ...interface{}) *kind {
	ret := &kind{
		superKind: k,
		name:      kindNameCheck(name),
	}
	for _, param := range codeAndErrMsgDef {
		switch p := param.(type) {
		case string:
			ret.errMsgDef = &p
		case int:
			c := int64(p)
			ret.codeDef = &c
		case int64:
			ret.codeDef = &p
		}
	}
	if ret.codeDef == nil {
		ret.codeDef = k.codeDef
	}
	if ret.errMsgDef == nil {
		ret.errMsgDef = k.errMsgDef
	}
	return ret
}
