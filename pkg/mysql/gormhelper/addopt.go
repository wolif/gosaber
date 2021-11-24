package gormhelper

import (
	"github.com/wolif/gosaber/pkg/log"
	"github.com/wolif/gosaber/pkg/ref"
	"strings"
	"time"
)

type addOpt = string

const (
	AddOptTimestampInt2Str   addOpt = "TimestampIntToStr"
	AddOptTimestampStr2Int64 addOpt = "TimestampStrToInt64"
)

var addOptMap = map[addOpt]func(interface{}) interface{}{
	AddOptTimestampInt2Str:   addOptTimestampInt2Str,
	AddOptTimestampStr2Int64: addOptTimestampStr2Int64,
}

func ExtendAddOpt(addOptName string, fn func(interface{}) interface{}, options ...bool) {
	if addOptName == "" {
		return
	}

	recoverFn := true
	if len(options) > 0 {
		recoverFn = options[0]
	}

	if !AddOptExist(addOptName) || recoverFn {
		addOptMap[addOptName] = fn
	}
}

func AddOptExist(opt addOpt) bool {
	_, ok := addOptMap[opt]
	return ok
}

func calAddOptValue(value interface{}, opt addOpt) interface{} {
	opt = strings.TrimSpace(opt)
	if opt == "" {
		return value
	}
	if !AddOptExist(opt) {
		log.Errorf("addition opt [%s] gorm helper not support", opt)
		return value
	}
	fn, _ := addOptMap[opt]
	return fn(value)
}

func addOptTimestampInt2Str(value interface{}) interface{} {
	format := "2006-01-02 15:04:05"
	vRef := ref.New(value)
	if vRef.IsInt() {
		return time.Unix(vRef.GetValue().Int(), 0).Format(format)
	}
	log.Errorf("convert data error, opt: [TimestampInt2Str], value: [%s], value kind: [%s]", value, vRef.GetKind())
	return value
}

func addOptTimestampStr2Int64(value interface{}) interface{} {
	format := "2006-01-02 15:04:05"
	location, _ := time.LoadLocation("Asia/Shanghai")

	vRef := ref.New(value)
	if vRef.IsString() {
		tt, err := time.ParseInLocation(format, vRef.GetValue().String(), location)
		if err == nil {
			return tt.Unix()
		}
	}
	log.Errorf("convert data error, opt: [TimestampStr2Int64], value: [%s], value kind: [%s]", value, vRef.GetKind())

	return value
}
