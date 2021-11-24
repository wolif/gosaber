package mongoAdapt

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/wolif/gosaber/pkg/log"
	"github.com/wolif/gosaber/pkg/ref"
	"strings"
)

type CMPOpt = string

const (
	OptEQ   CMPOpt = "="
	OptNEQ  CMPOpt = "neq"
	OptGT   CMPOpt = "gt"
	OptGTE  CMPOpt = "gte"
	OptLT   CMPOpt = "lt"
	OptLTE  CMPOpt = "lte"
	OptIN   CMPOpt = "in"
	OptNIN  CMPOpt = "nin"
	OptLIKE CMPOpt = "regex"
)

type cmpOptFunc struct {
	Fn                    func(mh *bson.M, opt string, field string, value interface{}) *bson.M
	WithoutCalAddOptValue bool
}

var cmpOpts = map[CMPOpt]*cmpOptFunc{
	OptEQ:   {Fn: cmnCMPOpt},
	OptNEQ:  {Fn: cmnCMPOpt},
	OptGT:   {Fn: cmnCMPOpt},
	OptGTE:  {Fn: cmnCMPOpt},
	OptLT:   {Fn: cmnCMPOpt},
	OptLTE:  {Fn: cmnCMPOpt},
	OptIN:   {Fn: cmnCMPOpt},
	OptNIN:  {Fn: cmnCMPOpt},
	OptLIKE: {Fn: cmnCMPOpt},
}

func CmpOptExist(opt CMPOpt) bool {
	_, found := cmpOpts[opt]
	return found
}

func ExtendCmpOpt(opt string, withoutCalAddOptValue bool, fn func(mh *bson.M, opt string, field string, value interface{}) *bson.M, options ...bool) {
	if opt == "" {
		return
	}

	recoverFn := true
	if len(options) > 0 {
		recoverFn = options[0]
	}

	if !CmpOptExist(opt) || recoverFn {
		cmpOpts[opt] = &cmpOptFunc{
			Fn:                    fn,
			WithoutCalAddOptValue: withoutCalAddOptValue,
		}
	}
}

func resolveOmitEmptyTag(tag string) bool {
	return ref.IsInScope(strings.ToUpper(strings.TrimSpace(tag)), "", "Y", "YES", "1")
}

func (mh *matchResult) matchField(field string, value interface{}, needOmitEmpty string, addOpt ADDOpt) {
	if ref.New(value).IsSlice() {
		mh.matchOptField(field, OptIN, value, needOmitEmpty, addOpt)
	} else {
		mh.matchOptField(field, OptEQ, value, needOmitEmpty, addOpt)
	}
}

func (mh *matchResult) matchOptField(field string, opt CMPOpt, value interface{}, needOmitEmpty string, addOpt ADDOpt) {
	opt = strings.TrimSpace(opt)
	field = strings.TrimSpace(field)
	if opt == "" || field == "" {
		return
	}

	if !CmpOptExist(opt) {
		log.Errorf("opt [%s] mongo match helper not support", opt)
		return
	}

	optFn, _ := cmpOpts[opt]
	if optFn.WithoutCalAddOptValue {
		mh = (*matchResult)(optFn.Fn((*bson.M)(mh), opt, field, value))
		return
	}


	vRef := ref.New(value)
	if resolveOmitEmptyTag(needOmitEmpty) && vRef.IsZero() {
		return
	}

	// 数据额外操作
	if addOpt != "" {
		value = calAddOptValue(value, addOpt)
	}

	mh = (*matchResult)(optFn.Fn((*bson.M)(mh), opt, field, value))
	return
}

func cmnCMPOpt(mh *bson.M, opt string, field string, value interface{}) *bson.M {
	if opt == OptEQ {
		(*mh)[field] = value
	} else {
		(*mh)[field] = bson.M{
			fmt.Sprintf("$%s", opt): value,
		}
	}
	return mh
}
