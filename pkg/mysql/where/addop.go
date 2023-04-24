package where

import (
	"log"
	"strings"
	"time"

	"github.com/wolif/gosaber/pkg/ref"
)

type addOp = string

const (
	AddOptTimestampInt2Str   addOp = "TimestampIntToStr"
	AddOptTimestampStr2Int64 addOp = "TimestampStrToInt64"
)

var addOpMap = map[addOp]func(interface{}) interface{}{
	AddOptTimestampInt2Str:   addOpTimestampInt2Str,
	AddOptTimestampStr2Int64: addOpTimestampStr2Int64,
}

func ExtendAddOp(addOpName string, fn func(interface{}) interface{}, options ...bool) {
	if addOpName == "" {
		return
	}

	recoverFn := true
	if len(options) > 0 {
		recoverFn = options[0]
	}

	if !AddOptExist(addOpName) || recoverFn {
		addOpMap[addOpName] = fn
	}
}

func AddOptExist(op addOp) bool {
	_, ok := addOpMap[op]
	return ok
}

func calAddOptValue(value interface{}, op addOp) interface{} {
	op = strings.TrimSpace(op)
	if op == "" {
		return value
	}
	if !AddOptExist(op) {
		log.Fatalf("mysql addOp [%s] not support", op)
		return value
	}
	return addOpMap[op](value)
}

func addOpTimestampInt2Str(value interface{}) interface{} {
	format := "2006-01-02 15:04:05"
	vRef := ref.New(value)
	if vRef.IsInt() {
		return time.Unix(vRef.GetValue().Int(), 0).Format(format)
	} else {
		log.Fatalf("mysql addOp [TimestampIntToStr] not support other type param except int")
	}
	return value
}

func addOpTimestampStr2Int64(value interface{}) interface{} {
	format := "2006-01-02 15:04:05"
	location, _ := time.LoadLocation("Asia/Shanghai")

	vRef := ref.New(value)
	if vRef.IsString() {
		tt, err := time.ParseInLocation(format, vRef.GetValue().String(), location)
		if err == nil {
			return tt.Unix()
		}
		log.Fatalf("mysql addOp [OptTimestampStr2Int64] not support param with other time format except 2006-01-02 15:04:05")
	}
	return value
}
