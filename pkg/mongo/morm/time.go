package morm

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/wolif/gosaber/pkg/ref"
	"github.com/wolif/gosaber/pkg/ref/structtags"
	"github.com/wolif/gosaber/pkg/util/strs"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	timestamp  = "timestamp"
	TimeFormat = "time_format"
	CreateTime = "create_time"
	UpdateTime = "update_time"
	DeleteTime = "delete_time"
)

func resolveTimeFields(model Model) *modelOpts {
	rm := ref.New(model)
	mt := rm.GetType()

	mormTags, err := structtags.Parse(model, "morm", map[string]string{"": "", TimeFormat: ""})
	if err != nil {
		panic(err)
	}

	for fieldInStruct, opts := range mormTags {
		for timeType, timeOpt := range map[string]*timeOpt{
			CreateTime: optCaches[mt].createTime,
			UpdateTime: optCaches[mt].updateTime,
			DeleteTime: optCaches[mt].deleteTime,
		} {
			if fieldInDB, ok := opts[timeType]; ok {
				timeOpt.fieldDB = fieldInDB
			} else if opts[""] == timeType {
				bsonTag, _ := rm.StructTagGet(fieldInStruct, "bson")
				if bsonTag = strings.TrimSpace(bsonTag); bsonTag != "" {
					if segs := strings.SplitN(bsonTag, ",", 1); len(segs) == 1 {
						if s := strings.TrimSpace(segs[0]); s != "" {
							timeOpt.fieldDB = s
							goto SSS
						}
					}
				}
				timeOpt.fieldDB = strs.SnakeCase(fieldInStruct)
			} else {
				continue
			}
		SSS:
			// 检查字段类型
			field, _ := rm.StructFieldGet(fieldInStruct)
			if field.Type.Kind() != reflect.Int64 && field.Type.Elem().Kind() != reflect.Int64 {
				panic("暂不支持 int64 和 *int64 之外的时间形式")
			}
			timeOpt.fieldStruct = fieldInStruct
			timeOpt.format = opts[TimeFormat]
		}
	}

	return optCaches[mt]
}

func FillFieldsWithTime(m bson.M, fields []string, timeFormat []string, ts ...time.Time) bson.M {
	t := time.Now()
	if len(ts) > 0 {
		t = ts[0]
	}
	for i, field := range fields {
		if field != "" {
			if timeFormat[i] == "" {
				m[field] = t.Unix()
			} else {
				timeStr := t.Format(timeFormat[i])
				var err error
				m[field], err = strconv.ParseInt(timeStr, 10, 64)
				if err != nil {
					panic(fmt.Sprintf("time string [%s] can't convert to int64", timeStr))
				}
			}
		}
	}
	return m
}

func addTimeWhenCreate(model Model, m bson.M) bson.M {
	opts := resolveModelOpts(model)
	if opts.deleteTime.fieldDB != "" {
		m[opts.deleteTime.fieldDB] = int64(0)
	}
	return FillFieldsWithTime(
		m,
		[]string{opts.createTime.fieldDB, opts.updateTime.fieldDB},
		[]string{opts.createTime.format, opts.updateTime.format},
	)
}

func addTimeWhenUpdate(model Model, m bson.M) bson.M {
	opts := resolveModelOpts(model)
	return FillFieldsWithTime(
		m,
		[]string{opts.updateTime.fieldDB},
		[]string{opts.updateTime.format},
	)
}

func addTimeWhenDelete(model Model, m bson.M) bson.M {
	opts := resolveModelOpts(model)
	return FillFieldsWithTime(
		m,
		[]string{opts.deleteTime.fieldDB},
		[]string{opts.deleteTime.format},
	)
}

func addTimeWhenWhere(model Model, m bson.M) bson.M {
	opts := resolveModelOpts(model)
	if opts.deleteTime.fieldDB != "" {
		m[opts.deleteTime.fieldDB] = int64(0)
	}
	return m
}
