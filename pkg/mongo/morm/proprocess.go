package morm

import "go.mongodb.org/mongo-driver/bson"

func HandleCreateModel(model Model) bson.M {
	return addTimeWhenCreate(model, Bson.ToModel(model, true))
}

func HandleUpdateData(model Model, updateData interface{}) bson.M {
	return addTimeWhenUpdate(model, Bson.ToModel(updateData))
}

func HandleDeleteData(model Model) bson.M {
	return addTimeWhenDelete(model, bson.M{})
}

func HandleWhereData(model Model, whereData interface{}) bson.M {
	return addTimeWhenWhere(model, Bson.ToModel(whereData))
}
