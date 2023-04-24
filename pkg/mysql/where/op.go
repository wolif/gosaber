package where

import (
	"fmt"
	"log"
	"strings"

	"github.com/wolif/gosaber/pkg/ref"
)

type cmpOp = string

const (
	OpEQ        cmpOp = "="
	OpNEQ       cmpOp = "!="
	OpGT        cmpOp = ">"
	OpGTE       cmpOp = ">="
	OpLT        cmpOp = "<"
	OpLTE       cmpOp = "<="
	OpIN        cmpOp = "in"
	OpNIN       cmpOp = "not in"
	OpLIKE      cmpOp = "like"
	OpRLIKE     cmpOp = "like abc%"
	OpLLIKE     cmpOp = "like %abc"
	OpNLIKE     cmpOp = "not like"
	OpLNLIKE    cmpOp = "not like %abc"
	OpRNLIKE    cmpOp = "not like abc%"
	OpFINDINSET cmpOp = "find in set"
	OpRawSql    cmpOp = "sql"
)

var opMap = map[cmpOp]func(*WhereVal, string, string, interface{}){
	OpEQ: compare, OpNEQ: compare, OpGT: compare, OpGTE: compare, OpLT: compare, OpLTE: compare,
	OpIN: scope, OpNIN: scope,
	OpLIKE: like, OpLLIKE: like, OpRLIKE: like, OpNLIKE: like, OpLNLIKE: like, OpRNLIKE: like,
	OpFINDINSET: findInSet,
	OpRawSql:    rawSql,
}

func OpExist(op cmpOp) bool {
	_, ok := opMap[op]
	return ok
}

func ExtendOp(op string, fn func(wv *WhereVal, op cmpOp, field string, value interface{}), cover ...bool) {
	c := false
	if len(cover) > 0 {
		c = cover[0]
	}
	if !OpExist(op) || c {
		opMap[op] = fn
	}
}

func resolveOmitEmptyTag(tag string) bool {
	return ref.In(strings.ToUpper(strings.TrimSpace(tag)), "", "Y", "YES", "1")
}

func where(wv *WhereVal, field string, op string, value interface{}, needOmitEmpty string, addOp addOp) {
	field = strings.TrimSpace(field)
	if field == "" {
		return
	}

	op = strings.TrimSpace(op)
	if op == "" {
		if ref.New(value).IsSlice() {
			op = OpIN
		} else {
			op = OpEQ
		}
	}
	if !OpExist(op) {
		log.Fatalf("mysql search op [%s] not support", op)
		return
	}

	vRef := ref.New(value)
	if resolveOmitEmptyTag(needOmitEmpty) && vRef.GetValue().IsZero() {
		return
	}

	// 数据额外操作
	if addOp != "" {
		value = calAddOptValue(value, addOp)
	}

	opMap[op](wv, field, op, value)
}

func compare(wv *WhereVal, field string, op string, value interface{}) {
	rv := ref.New(value)
	if rv.IsStringOrNumber() {
		wv.Sql = fmt.Sprintf("%s AND `%s` %s ?", wv.Sql, field, op)
		wv.Value = append(wv.Value, value)
	}
}

func scope(wv *WhereVal, field string, op string, value interface{}) {
	rv := ref.New(value)
	if rv.IsSlice() && rv.GetValue().Len() > 0 {
		if rv.GetValue().Len() == 1 {
			oMap := map[string]string{OpIN: OpEQ, OpNIN: OpNEQ}
			compare(wv, field, oMap[op], rv.GetValue().Index(0).Interface())
		} else {
			wv.Sql = fmt.Sprintf("%s AND `%s` %s (?)", wv.Sql, field, strings.ToUpper(op))
			wv.Value = append(wv.Value, value)
		}
	}
}

func like(wv *WhereVal, field string, op string, value interface{}) {
	rv := ref.New(value)
	if rv.IsString() {
		vMap := map[string]string{
			OpLIKE:   fmt.Sprintf("%%%s%%", value),
			OpNLIKE:  fmt.Sprintf("%%%s%%", value),
			OpLLIKE:  fmt.Sprintf("%%%s", value),
			OpLNLIKE: fmt.Sprintf("%%%s", value),
			OpRLIKE:  fmt.Sprintf("%s%%", value),
			OpRNLIKE: fmt.Sprintf("%s%%", value),
		}
		oMap := map[string]string{
			OpLIKE:   OpLIKE,
			OpLLIKE:  OpLIKE,
			OpRLIKE:  OpLIKE,
			OpNLIKE:  OpNLIKE,
			OpLNLIKE: OpNLIKE,
			OpRNLIKE: OpNLIKE,
		}
		wv.Sql = fmt.Sprintf("%s AND `%s` %s ?", wv.Sql, field, strings.ToUpper(oMap[op]))
		wv.Value = append(wv.Value, vMap[op])
	}
}

func findInSet(wv *WhereVal, field string, op string, value interface{}) {
	rv := ref.New(value)
	if rv.IsStringOrNumber() {
		wv.Sql = fmt.Sprintf("%s AND find_in_set(?, `%s`)", wv.Sql, field)
		wv.Value = append(wv.Value, fmt.Sprint(value))
	}
}

func rawSql(wv *WhereVal, field string, op string, value interface{}) {
	wv.Sql = fmt.Sprintf("%s AND %s", wv.Sql, value)
}
