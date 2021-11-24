package utils

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/wolif/gosaber/pkg/mongo/mongoAdapt"
	"github.com/wolif/gosaber/pkg/ref"
	"reflect"
	"strings"
	"time"
)

type Model interface {
	GetCollectionName() string
}

var modelInfo = make(map[reflect.Type]map[string]string)

func fetch3TimeFields(model Model) (createTimeField string, updateTimeField string, deleteTimeField string) {
	rType := reflect.TypeOf(ref.PeelOffPtr(model))
	if rType.Kind() != reflect.Struct {
		return "", "", ""
	}

	info, found := modelInfo[rType]
	if !found {
		info = map[string]string{
			"created_at": "",
			"updated_at": "",
			"deleted_at": "",
		}

		mapBsonTagFieldName := map[string]string{}
		mapFieldNameBsonTag := map[string]string{}
		for i := 0; i < rType.NumField(); i++ {
			field := rType.Field(i)
			mapFieldNameBsonTag[field.Name] = ""
			bsonTagStr := field.Tag.Get("bson")
			bsonTagSegs := strings.SplitN(bsonTagStr, " ", 1)
			if len(bsonTagSegs) > 0 {
				bsonTag := strings.TrimSpace(bsonTagSegs[0])
				if len(bsonTag) > 0 {
					mapBsonTagFieldName[bsonTag] = field.Name
				}
				mapFieldNameBsonTag[field.Name] = bsonTag
			}
		}
		for i := 0; i < rType.NumField(); i++ {
			field := rType.Field(i)
			mongoTag := field.Tag.Get("mongo")
			tagSegsOrigin := strings.Split(mongoTag, " ")
			tagSegs := make([]string, 0)
			for _, v := range tagSegsOrigin {
				if s := strings.TrimSpace(v); s != "" {
					tagSegs = append(tagSegs, s)
				}
			}
			// 优先找 autoxxxTime tag
			if ref.IsInSlice("autoCreateTime", tagSegs) {
				if bsonTagName, found := mapFieldNameBsonTag[field.Name]; found && bsonTagName != "" {
					info["created_at"] = bsonTagName
				} else {
					info["created_at"] = field.Name
				}
			} else if ref.IsInSlice("autoUpdateTime", tagSegs) {
				if bsonTagName, found := mapFieldNameBsonTag[field.Name]; found && bsonTagName != "" {
					info["updated_at"] = bsonTagName
				} else {
					info["updated_at"] = field.Name
				}
			} else if ref.IsInSlice("autoDeleteTime", tagSegs) {
				if bsonTagName, found := mapFieldNameBsonTag[field.Name]; found && bsonTagName != "" {
					info["deleted_at"] = bsonTagName
				} else {
					info["deleted_at"] = field.Name
				}
			}
		}

		for _, c := range []string{"created_at", "updated_at", "deleted_at"} {
			// 次优先 找 bson tag
			if info[c] == "" {
				if _, found := mapBsonTagFieldName[c]; found {
					info[c] = c
					continue
				}
			}
			//// 最次找 fieldName
			//if info[c] == "" {
			//	if tag, found := mapFieldNameBsonTag[strs.CamelString(c)]; found {
			//		if len(tag) > 0 {
			//			info[c] = tag
			//		} else {
			//			info[c] = c
			//		}
			//	}
			//}
		}

		modelInfo[rType] = info
	}
	return info["created_at"], info["updated_at"], info["deleted_at"]
}

func structToBson(structData interface{}) bson.M {
	ret := bson.M{}

	refModel := ref.New(structData)
	fields, isStruct := refModel.GetStructFields()
	if !isStruct {
		return ret
	}

	for fn, _ := range fields {
		keyName := fn
		bsonTag, _ := refModel.GetStructFieldTag(fn, "bson")
		if bsonTag != "" {
			segs := strings.SplitN(bsonTag, " ", 1)
			if strings.TrimSpace(segs[0]) != "" {
				keyName = strings.TrimSpace(segs[0])
			}
		}
		ret[keyName], _ = refModel.GetStructFieldValue(fn)
	}

	return ret
}

func mapToBson(data interface{}) bson.M {
	ret := bson.M{}
	refData := ref.New(data)
	for _, mk := range refData.GetValue().MapKeys() {
		k := mk.String()
		v, _ := refData.GetMapValue(k)
		ret[k] = v
	}
	return ret
}

func dataToModelBson(data interface{}, makeNewIDIfNull ...bool) bson.M {
	ret := bson.M{}
	refData := ref.New(data)
	if refData.IsStruct() {
		ret = structToBson(data)
	} else if refData.GetType() == reflect.TypeOf(ret) {
		ret = refData.GetValue().Interface().(bson.M)
	} else if refData.IsMap() {
		ret = mapToBson(data)
	} else {
		return ret
	}

	id, found := ret["_id"]
	if !found || reflect.ValueOf(id).IsZero() {
		needMakeNewID := false
		if len(makeNewIDIfNull) > 0 {
			needMakeNewID = makeNewIDIfNull[0]
		}
		if needMakeNewID {
			ret = addMongoID(ret, "_id")
		}
	}
	return ret
}

func addMongoID(m bson.M, IDField string) bson.M {
	m[IDField] = primitive.NewObjectID()
	return m
}

func addTimeNow(m bson.M, timeFields ...string) bson.M {
	now := time.Now().Unix()
	for _, field := range timeFields {
		if field != "" {
			m[field] = now
		}
	}
	return m
}

func addTimeWhenCreate(model Model, m bson.M) bson.M {
	c, u, d := fetch3TimeFields(model)
	if d != "" {
		m["deleted_at"] = int64(0)
	}
	return addTimeNow(m, c, u)
}

func addTimeWhenUpdate(model Model, m bson.M) bson.M {
	_, u, _ := fetch3TimeFields(model)
	return addTimeNow(m, u)
}

func addTimeWhenDelete(model Model, m bson.M) bson.M {
	_, _, d := fetch3TimeFields(model)
	return addTimeNow(m, d)
}

func addTimeWhenSearch(model Model, m bson.M) bson.M {
	_, _, d := fetch3TimeFields(model)
	if d != "" {
		m[d] = int64(0)
	}
	return m
}

func HandleFilter(model Model, filter interface{}) bson.M {
	match := filter
	if ref.New(model).IsStruct() {
		match = mongoAdapt.Match(filter)
	}
	return addTimeWhenSearch(model, dataToModelBson(match))
}

func HandleCreateModel(model Model) bson.M {
	return addTimeWhenCreate(model, dataToModelBson(model, true))
}

func HandleUpdateData(model Model, data interface{}) bson.M {
	return addTimeWhenUpdate(model, dataToModelBson(data))
}

func HandleDeleteData(model Model) bson.M {
	return addTimeWhenDelete(model, bson.M{})
}

func handleSort(sort map[string]int) bson.M {
	s := bson.M{}
	if sort != nil && len(sort) > 0 {
		for field, sv := range sort {
			if sv >= 0 {
				s[field] = 1
			} else {
				s[field] = -1
			}
		}
	} else {
		s["_id"] = -1
	}
	return s
}
