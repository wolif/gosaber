package morm

import (
	"reflect"
	"strings"

	"github.com/wolif/gosaber/pkg/ref"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type bsonConv struct{}

var Bson = new(bsonConv)

func (bc *bsonConv) FromStruct(structData interface{}) bson.M {
	ret := bson.M{}

	refModel := ref.New(structData)
	fields, isStruct := refModel.StructFields()
	if !isStruct {
		return ret
	}

	for fn, _ := range fields {
		keyName := fn
		bsonTag, _ := refModel.StructTagGet(fn, "bson")
		if bsonTag != "" {
			if s := strings.TrimSpace(strings.SplitN(bsonTag, " ", 1)[0]); s != "" {
				keyName = s
			}
		}
		ret[keyName], _ = refModel.StructValueGet(fn)
	}

	return ret
}

func (bc *bsonConv) FromMap(data interface{}) bson.M {
	ret := bson.M{}
	refData := ref.New(data)
	for _, mk := range refData.GetValue().MapKeys() {
		v, _ := refData.MapGet(mk.String())
		ret[mk.String()] = v
	}
	return ret
}

func (bc *bsonConv) ToModel(data interface{}, newIDIfNull ...bool) bson.M {
	ret := bson.M{}
	refData := ref.New(data)
	if refData.IsStruct() {
		ret = bc.FromStruct(data)
	} else if refData.GetType() == reflect.TypeOf(ret) {
		ret = refData.GetValue().Interface().(bson.M)
	} else if refData.IsMap() {
		ret = bc.FromMap(data)
	} else {
		return ret
	}

	if id, found := ret["_id"]; !found || reflect.ValueOf(id).IsZero() {
		if len(newIDIfNull) > 0 && newIDIfNull[0] {
			ret["_id"] = primitive.NewObjectID()
		}
	}
	return ret
}
