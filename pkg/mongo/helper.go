package mongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func StrID2ObjID(id string) interface{} {
	oid, _ := primitive.ObjectIDFromHex(id)
	return oid
}

func StrIDs2ObjIDs(ids []string) []interface{} {
	res := make([]interface{}, 0, len(ids))
	for _, ID := range ids {
		id, _ := primitive.ObjectIDFromHex(ID)
		res = append(res, id)
	}
	return res
}
