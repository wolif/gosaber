package morm

import (
	"reflect"

	"github.com/wolif/gosaber/pkg/ref"
)

type Model interface {
	GetCollectionName() string
}

type timeOpt struct {
	fieldStruct string
	fieldDB     string
	format      string
}

type modelOpts struct {
	createTime *timeOpt
	updateTime *timeOpt
	deleteTime *timeOpt
}

var optCaches = make(map[reflect.Type]*modelOpts)

func resolveModelOpts(model Model) *modelOpts {
	t := ref.New(model).GetType()
	if opts, ok := optCaches[t]; ok {
		return opts
	}
	optCaches[t] = &modelOpts{
		createTime: &timeOpt{},
		updateTime: &timeOpt{},
		deleteTime: &timeOpt{},
	}
	return resolveTimeFields(model)
}
