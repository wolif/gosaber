package utils

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/wolif/gosaber/pkg/mongo"
	"github.com/wolif/gosaber/pkg/mongo/mongoAdapt"
	"reflect"
)

type record struct {
	dbConnName string
}

var Record = &record{dbConnName: "default"}

func (r *record) SetDBConnName(name string) {
	r.dbConnName = name
}

func (r *record) InsertMany(c context.Context, models ...Model) ([]string, error) {
	ret := make([]string, 0)
	if len(models) <= 0 {
		return ret, nil
	}

	t := reflect.TypeOf(models[0])
	for i := 0; i < len(models); i++ {
		if reflect.TypeOf(models[i]) != t {
			return ret, fmt.Errorf("different types in models")
		}
	}

	db, err := mongo.GetDB(r.dbConnName)
	if err != nil {
		return ret, err
	}

	data := make([]interface{}, 0)
	for _, m := range models {
		data = append(data, HandleCreateModel(m))
	}
	res, err := db.Collection(models[0].GetCollectionName()).InsertMany(c, data)
	if err != nil {
		return ret, err
	}

	for _, id := range res.InsertedIDs {
		ret = append(ret, id.(primitive.ObjectID).Hex())
	}

	return ret, nil
}

func (r *record) InsertOne(c context.Context, model Model) (string, error) {
	db, err := mongo.GetDB(r.dbConnName)
	if err != nil {
		return "", err
	}

	res, err := db.Collection(model.GetCollectionName()).InsertOne(
		c,
		HandleCreateModel(model),
	)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *record) Save(c context.Context, model Model, filter interface{}, data interface{}) (int64, error) {
	db, err := mongo.GetDB(r.dbConnName)
	if err != nil {
		return 0, err
	}

	updateRes, err := db.Collection(model.GetCollectionName()).UpdateOne(
		c,
		HandleFilter(model, filter),
		HandleUpdateData(model, data),
	)
	if err != nil {
		return 0, err
	}
	return updateRes.ModifiedCount, nil
}

func (r *record) UpdateOne(c context.Context, model Model, filter interface{}, data interface{}) (int64, error) {
	db, err := mongo.GetDB(r.dbConnName)
	if err != nil {
		return 0, err
	}

	updateRes, err := db.Collection(model.GetCollectionName()).UpdateOne(
		c,
		HandleFilter(model, filter),
		bson.M{
			"$set": HandleUpdateData(model, data),
		},
	)
	if err != nil {
		return 0, err
	}
	return updateRes.ModifiedCount, nil
}

func (r *record) UpdateMany(c context.Context, model Model, filter interface{}, data interface{}) (int64, error) {
	db, err := mongo.GetDB(r.dbConnName)
	if err != nil {
		return 0, err
	}

	updateRes, err := db.Collection(model.GetCollectionName()).UpdateMany(
		c,
		HandleFilter(model, filter),
		bson.M{
			"$set": HandleUpdateData(model, data),
		},
	)
	if err != nil {
		return 0, err
	}
	return updateRes.ModifiedCount, nil
}

func (r *record) Delete(c context.Context, model Model, filter interface{}, isLogicDelete ...bool) (int64, error) {
	db, err := mongo.GetDB(r.dbConnName)
	if err != nil {
		return 0, err
	}

	logicDelete := true
	if len(isLogicDelete) > 0 {
		logicDelete = isLogicDelete[0]
	}
	if logicDelete {
		updateRes, err := db.Collection(model.GetCollectionName()).UpdateMany(
			c,
			HandleFilter(model, filter),
			bson.M{
				"$set": HandleDeleteData(model),
			},
		)
		if err != nil {
			return 0, err
		}
		return updateRes.ModifiedCount, nil
	}

	match := filter
	if reflect.TypeOf(filter).Kind() == reflect.Struct {
		match = mongoAdapt.Match(filter)
	}
	// 物理删除不考虑是否已经被逻辑删除
	deleteRes, err := db.Collection(model.GetCollectionName()).DeleteMany(c, match)
	if err != nil {
		return 0, err
	}
	return deleteRes.DeletedCount, nil
}

func (r *record) FindOne(c context.Context, model Model, filter interface{}, result interface{}) error {
	db, err := mongo.GetDB(r.dbConnName)
	if err != nil {
		return err
	}

	res := db.Collection(model.GetCollectionName()).FindOne(
		c,
		HandleFilter(model, filter),
	)
	if res.Err() != nil {
		return res.Err()
	}
	err = res.Decode(result)
	if err != nil {
		return err
	}
	return nil
}

func (r *record) SearchRecord(c context.Context, model Model, filter interface{}, sort map[string]int, result interface{}, pagination ...int) (int64, error) {
	db, err := mongo.GetDB(r.dbConnName)
	if err != nil {
		return 0, err
	}

	col := db.Collection(model.GetCollectionName())
	match := HandleFilter(model, filter)
	cnt, err := col.CountDocuments(c, match)
	if err != nil {
		return 0, err
	}

	page, pageSize := resolvePagination(pagination...)
	fo := options.Find().
		SetSort(handleSort(sort)).
		SetSkip(int64((page - 1) * pageSize)).
		SetLimit(int64(pageSize))
	cur, err := col.Find(c, match, fo)
	if err != nil {
		return 0, err
	}
	if cur.Err() != nil {
		return 0, cur.Err()
	}
	err = cur.All(c, result)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

func (r *record) Search(c context.Context, model Model, filter interface{}, sort map[string]int) (*mongo2.Cursor, error) {
	db, err := mongo.GetDB(r.dbConnName)
	if err != nil {
		return nil, err
	}

	return db.Collection(model.GetCollectionName()).Find(
		c,
		HandleFilter(model, filter),
		options.Find().SetSort(handleSort(sort)),
	)
}
