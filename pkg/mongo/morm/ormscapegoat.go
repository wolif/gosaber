package morm

import (
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 替罪羊
type scapegoat struct {
	orm *orm
}

/**
 * 插入一条数据
 */
func (si *scapegoat) Insert(model Model, options ...*options.InsertOneOptions) (string, error) {
	defer si.orm.destroy()
	si.orm.Model(model)
	res, err := si.orm.db.Collection(si.orm.coll).InsertOne(
		si.orm.ctx,
		HandleCreateModel(model),
		options...,
	)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

/**
 * 插入多条数据
 */
func (si *scapegoat) Inserts(models []Model, options ...*options.InsertManyOptions) ([]string, error) {
	defer si.orm.destroy()
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

	data := make([]interface{}, 0)
	for _, m := range models {
		data = append(data, HandleCreateModel(m))
	}
	si.orm.Model(models[0])
	res, err := si.orm.db.Collection(si.orm.coll).InsertMany(si.orm.ctx, data, options...)
	if err != nil {
		return ret, err
	}

	for _, id := range res.InsertedIDs {
		ret = append(ret, id.(primitive.ObjectID).Hex())
	}

	return ret, nil
}

/**
 * 更新一条数据
 */
func (si *scapegoat) Update(data interface{}, opts ...*options.UpdateOptions) (int64, error) {
	defer si.orm.destroy()
	updateRes, err := si.orm.db.Collection(si.orm.coll).UpdateOne(
		si.orm.ctx,
		HandleWhereData(si.orm.model, si.orm.where),
		HandleUpdateData(si.orm.model, data),
		opts...,
	)
	if err != nil {
		return 0, err
	}
	return updateRes.ModifiedCount, nil
}

/**
 * 更新多条数据
 */
func (si *scapegoat) Updates(data interface{}, opts ...*options.UpdateOptions) (int64, error) {
	defer si.orm.destroy()
	updateRes, err := si.orm.db.Collection(si.orm.coll).UpdateMany(
		si.orm.ctx,
		HandleWhereData(si.orm.model, si.orm.where),
		HandleUpdateData(si.orm.model, data),
		opts...,
	)
	if err != nil {
		return 0, err
	}
	return updateRes.ModifiedCount, nil
}

/**
 * 设置一条数据的字段
 */
func (si *scapegoat) Set(data interface{}, opts ...*options.UpdateOptions) (int64, error) {
	defer si.orm.destroy()
	updateRes, err := si.orm.db.Collection(si.orm.coll).UpdateOne(
		si.orm.ctx,
		HandleWhereData(si.orm.model, si.orm.where),
		bson.M{
			"$set": HandleUpdateData(si.orm.model, data),
		},
		opts...,
	)
	if err != nil {
		return 0, err
	}
	return updateRes.ModifiedCount, nil
}

/**
 * 设置多条数据的字段
 */
func (si *scapegoat) Sets(data interface{}, opts ...*options.UpdateOptions) (int64, error) {
	defer si.orm.destroy()
	updateRes, err := si.orm.db.Collection(si.orm.coll).UpdateMany(
		si.orm.ctx,
		HandleWhereData(si.orm.model, si.orm.where),
		bson.M{
			"$set": HandleUpdateData(si.orm.model, data),
		},
		opts...,
	)
	if err != nil {
		return 0, err
	}
	return updateRes.ModifiedCount, nil
}

/*
删除记录, 默认进行逻辑删除
*/
func (si *scapegoat) Delete(isLogicDelete ...bool) (int64, error) {
	defer si.orm.destroy()
	logicDelete := true
	if len(isLogicDelete) > 0 {
		logicDelete = isLogicDelete[0]
	}
	if logicDelete {
		deleteData := HandleDeleteData(si.orm.model)
		if len(deleteData) == 0 {
			return 0, fmt.Errorf("can't delete the orm logically: mark field not found")
		}
		updateRes, err := si.orm.db.Collection(si.orm.coll).UpdateMany(
			si.orm.ctx,
			HandleWhereData(si.orm.model, si.orm.where),
			bson.M{
				"$set": deleteData,
			},
		)
		if err != nil {
			return 0, err
		}
		return updateRes.ModifiedCount, nil
	}

	// 物理删除不考虑是否已经被逻辑删除
	deleteRes, err := si.orm.db.Collection(si.orm.coll).DeleteMany(si.orm.ctx, Bson.ToModel(si.orm.where))
	if err != nil {
		return 0, err
	}
	return deleteRes.DeletedCount, nil
}

/**
 * 数数
 */
func (si *scapegoat) Count(opts ...*options.CountOptions) (int64, error) {
	defer si.orm.destroy()
	cnt, err := si.orm.db.Collection(si.orm.coll).CountDocuments(
		si.orm.ctx,
		HandleWhereData(si.orm.model, si.orm.where),
		opts...,
	)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

/**
 * 查询一条记录
 */
func (si *scapegoat) FindOne(result interface{}, opts ...*options.FindOneOptions) error {
	defer si.orm.destroy()
	opt := options.FindOne()
	if len(si.orm.sort) > 0 {
		opt.SetSort(si.orm.sort)
	}
	if len(si.orm.projection) > 0 {
		opt.SetProjection(si.orm.projection)
	}
	res := si.orm.db.Collection(si.orm.coll).FindOne(
		si.orm.ctx,
		HandleWhereData(si.orm.model, si.orm.where),
		options.MergeFindOneOptions(append(opts, opt)...),
	)
	if res.Err() != nil {
		return res.Err()
	}
	err := res.Decode(result)
	if err != nil {
		return err
	}
	return nil
}

/**
 * 查询多条记录
 */
func (si *scapegoat) Find(result interface{}, opts ...*options.FindOptions) (int64, error) {
	defer si.orm.destroy()
	match := HandleWhereData(si.orm.model, si.orm.where)
	col := si.orm.db.Collection(si.orm.coll)
	cnt, err := col.CountDocuments(si.orm.ctx, match)
	if err != nil {
		return 0, err
	}

	defOpt := options.Find()
	if len(si.orm.sort) > 0 {
		defOpt.SetSort(si.orm.sort)
	}
	if si.orm.limit > 0 {
		defOpt.SetSkip(si.orm.skip).SetLimit(si.orm.limit)
	}
	if len(si.orm.projection) > 0 {
		defOpt.SetProjection(si.orm.projection)
	}
	cur, err := col.Find(
		si.orm.ctx,
		match,
		options.MergeFindOptions(append(opts, defOpt)...),
	)
	if err != nil {
		return 0, err
	}
	if cur.Err() != nil {
		return 0, cur.Err()
	}
	err = cur.All(si.orm.ctx, result)
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

/**
 * 查询所有记录
 */
func (si *scapegoat) FindAll(opts ...*options.FindOptions) (*mongo2.Cursor, error) {
	defer si.orm.destroy()
	defOpt := options.Find()
	if len(si.orm.sort) > 0 {
		defOpt.SetSort(si.orm.sort)
	}
	if len(si.orm.projection) > 0 {
		defOpt.SetProjection(si.orm.projection)
	}
	return si.orm.db.Collection(si.orm.coll).Find(
		si.orm.ctx,
		HandleWhereData(si.orm.model, si.orm.where),
		options.MergeFindOptions(append(opts, defOpt)...),
	)
}
