package madapt

import (
	"fmt"
	"log"
	"strings"

	"github.com/wolif/gosaber/pkg/ref"
	"go.mongodb.org/mongo-driver/bson"
)

type CMPOp = string

const (
	OpEQ   CMPOp = "="
	OpNEQ  CMPOp = "neq"
	OpGT   CMPOp = "gt"
	OpGTE  CMPOp = "gte"
	OpLT   CMPOp = "lt"
	OpLTE  CMPOp = "lte"
	OpIN   CMPOp = "in"
	OpNIN  CMPOp = "nin"
	OpLIKE CMPOp = "regex"
	OpRaw  CMPOp = "raw"
)

var cmpOps = map[CMPOp]func(*bson.M, string, string, interface{}) *bson.M{
	OpEQ:   cmnCMPOp,
	OpNEQ:  cmnCMPOp,
	OpGT:   cmnCMPOp,
	OpGTE:  cmnCMPOp,
	OpLT:   cmnCMPOp,
	OpLTE:  cmnCMPOp,
	OpIN:   cmnCMPOp,
	OpNIN:  cmnCMPOp,
	OpLIKE: cmnCMPOp,
	OpRaw:  cmnCMPRaw,
}

func CmpOpExist(op CMPOp) bool {
	_, found := cmpOps[op]
	return found
}

func ExtendCmpOp(op string, fn func(match *bson.M, op, field string, value interface{}) *bson.M, opions ...bool) {
	if op == "" {
		return
	}

	recoverFn := true
	if len(opions) > 0 {
		recoverFn = opions[0]
	}

	if !CmpOpExist(op) || recoverFn {
		cmpOps[op] = fn
	}
}

func resolveOmitEmptyTag(tag string) bool {
	return ref.In(strings.ToUpper(strings.TrimSpace(tag)), "", "Y", "YES", "1", "T", "TRUE")
}

func (mh *matchResult) matchOpField(field string, op CMPOp, value interface{}, omitEmpty string, addOp AddOp) {
	field = strings.TrimSpace(field)
	if field == "" {
		return
	}
	vRef := ref.New(value)

	//确定操作符
	op = strings.TrimSpace(op)
	if op == "" {
		if vRef.IsSlice() {
			op = OpIN
		} else {
			op = OpEQ
		}
	}
	// 判断操作符是否可用
	if !CmpOpExist(op) {
		log.Fatalf("op [%s] mongo match helper not support", op)
		return
	}

	// 空值可用与否
	if resolveOmitEmptyTag(omitEmpty) && (vRef.GetValue().IsZero() || (vRef.IsSlice() && vRef.GetOriValue().Len() == 0)) {
		return
	}

	// 数据额外操作
	if addOp != "" {
		for _, one := range strings.Split(addOp, "|") {
			if o := strings.TrimSpace(one); o != "" {
				value = calAddOpValue(value, o)
			}
		}
	}

	mh = (*matchResult)(cmpOps[op]((*bson.M)(mh), op, field, value))
}

func cmnCMPOp(mh *bson.M, op, field string, value interface{}) *bson.M {
	if op == OpEQ {
		(*mh)[field] = value
	} else {
		if _, ok := (*mh)[field]; !ok {
			(*mh)[field] = bson.M{fmt.Sprintf("$%s", op): value}
		} else {
			switch v := (*mh)[field].(type) {
			case bson.M:
				v[fmt.Sprintf("$%s", op)] = value
			}
		}
	}
	return mh
}

func cmnCMPRaw(mh *bson.M, op, field string, value interface{}) *bson.M {
	(*mh)[field] = value
	return mh
}
