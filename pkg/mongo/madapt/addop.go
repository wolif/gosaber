package madapt

import (
	"log"
	"strings"
	"time"

	"github.com/wolif/gosaber/pkg/mongo"
	"github.com/wolif/gosaber/pkg/ref"
)

type AddOp = string

const (
	AddOpTimeToStr    AddOp = "TimeToStr"
	AddOpStrToTime    AddOp = "StrToTime"
	AddOpStrIDToObjID AddOp = "StrIDToObjID"
)

var addOptMap = map[AddOp]func(interface{}) interface{}{
	AddOpTimeToStr:    addOptTimeToStr,
	AddOpStrToTime:    addOptStrToTime,
	AddOpStrIDToObjID: addOptStrIDToObjID,
}

func ExtendAddOp(addOpName string, fn func(interface{}) interface{}, options ...bool) {
	if addOpName == "" {
		return
	}

	recoverFn := true
	if len(options) > 0 {
		recoverFn = options[0]
	}

	if !AddOpExist(addOpName) || recoverFn {
		addOptMap[addOpName] = fn
	}
}

func AddOpExist(opt AddOp) bool {
	_, ok := addOptMap[opt]
	return ok
}

func calAddOpValue(value interface{}, opt AddOp) interface{} {
	opt = strings.TrimSpace(opt)
	if opt == "" {
		return value
	}
	if !AddOpExist(opt) {
		log.Fatalf("addition opt [%s] madapt helper not support", opt)
	}
	return addOptMap[opt](value)
}

func addOptTimeToStr(value interface{}) interface{} {
	format := "2006-01-02 15:04:05"
	vRef := ref.New(value)
	if vRef.IsInt() {
		return time.Unix(vRef.GetValue().Int(), 0).Format(format)
	}
	return value
}

func addOptStrToTime(value interface{}) interface{} {
	format := "2006-01-02 15:04:05"
	location, _ := time.LoadLocation("Asia/Shanghai")

	vRef := ref.New(value)
	if vRef.IsString() {
		tt, err := time.ParseInLocation(format, vRef.GetValue().String(), location)
		if err == nil {
			return tt.Unix()
		}
	}

	return value
}

func addOptStrIDToObjID(value interface{}) interface{} {
	valRef := ref.New(value)

	if valRef.IsSlice() {
		ids := make([]string, 0)
		for i := 0; i < valRef.GetValue().Len(); i++ {
			id := valRef.GetValue().Index(i).String()
			ids = append(ids, id)
		}
		return mongo.StrIDs2ObjIDs(ids)
	}

	id := ""
	if valRef.IsString() {
		id = valRef.GetValue().String()
	}
	return mongo.StrID2ObjID(id)
}
