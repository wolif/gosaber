package morm

import (
	"github.com/wolif/gosaber/pkg/mongo"
	mgo2 "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (o *orm) Aggregate(pipeline interface{}, opts ...*options.AggregateOptions) (*mgo2.Cursor, error) {
	defer o.destroy()
	return o.db.Collection(o.GetCollectionName()).Aggregate(
		o.ctx,
		pipeline,
		options.MergeAggregateOptions(append(opts, options.Aggregate().SetAllowDiskUse(true))...),
	)
}

func (o *orm) InsertOne(doc interface{}, opts ...*options.InsertOneOptions) (*mgo2.InsertOneResult, error) {
	defer o.destroy()
	return o.db.Collection(o.GetCollectionName()).InsertOne(o.ctx, doc, opts...)
}

func (o *orm) InsertMany(docs []interface{}, opts ...*options.InsertManyOptions) (*mgo2.InsertManyResult, error) {
	defer o.destroy()
	return o.db.Collection(o.GetCollectionName()).InsertMany(o.ctx, docs, opts...)
}

func (o *orm) UpdateOne(data interface{}, opts ...*options.UpdateOptions) (*mgo2.UpdateResult, error) {
	defer o.destroy()
	return o.db.Collection(o.GetCollectionName()).UpdateOne(
		o.ctx,
		HandleWhereData(o.model, o.where),
		data,
		opts...,
	)
}

func (o *orm) UpdateMany(data interface{}, opts ...*options.UpdateOptions) (*mgo2.UpdateResult, error) {
	defer o.destroy()
	return o.db.Collection(o.GetCollectionName()).UpdateMany(
		o.ctx,
		HandleWhereData(o.model, o.where),
		data,
		opts...,
	)
}

func (o *orm) UpdateByID(id interface{}, data interface{}, opts ...*options.UpdateOptions) (*mgo2.UpdateResult, error) {
	defer o.destroy()
	return o.db.Collection(o.GetCollectionName()).UpdateByID(o.ctx, id, data, opts...)
}

func (o *orm) FindOneAndUpdate(data interface{}, opts ...*options.FindOneAndUpdateOptions) *mgo2.SingleResult {
	defer o.destroy()
	return o.db.Collection(o.GetCollectionName()).FindOneAndUpdate(
		o.ctx,
		HandleWhereData(o.model, o.where),
		data,
		opts...,
	)
}

func (o *orm) FindOneAndReplace(data interface{}, opts ...*options.FindOneAndReplaceOptions) *mgo2.SingleResult {
	defer o.destroy()
	return o.db.Collection(o.GetCollectionName()).FindOneAndReplace(
		o.ctx,
		HandleWhereData(o.model, o.where),
		data,
		opts...,
	)
}

func (o *orm) FindOneAndDelete(opts ...*options.FindOneAndDeleteOptions) *mgo2.SingleResult {
	defer o.destroy()
	return o.db.Collection(o.GetCollectionName()).FindOneAndDelete(
		o.ctx,
		HandleWhereData(o.model, o.where),
		opts...,
	)
}

func (o *orm) FindOne(opts ...*options.FindOneOptions) *mgo2.SingleResult {
	defer o.destroy()
	return o.db.Collection(o.GetCollectionName()).FindOne(
		o.ctx, HandleWhereData(o.model, o.where), opts...,
	)
}

func (o *orm) Find(opts ...*options.FindOptions) (*mgo2.Cursor, error) {
	defer o.destroy()
	return o.db.Collection(o.GetCollectionName()).Find(
		o.ctx, HandleWhereData(o.model, o.where), opts...,
	)
}

func (o *orm) DeleteOne(opts ...*options.DeleteOptions) (*mgo2.DeleteResult, error) {
	defer o.destroy()
	return o.db.Collection(o.GetCollectionName()).DeleteOne(
		o.ctx, HandleWhereData(o.model, o.where), opts...,
	)
}

func (o *orm) DeleteMany(opts ...*options.DeleteOptions) (*mgo2.DeleteResult, error) {
	defer o.destroy()
	return o.db.Collection(o.GetCollectionName()).DeleteMany(
		o.ctx, HandleWhereData(o.model, o.where), opts...,
	)
}

func (o *orm) CountDocuments(opts ...*options.CountOptions) (int64, error) {
	defer o.destroy()
	return o.db.Collection(o.GetCollectionName()).CountDocuments(
		o.ctx, HandleWhereData(o.model, o.where), opts...,
	)
}

func (o *orm) EstimatedDocumentCount(opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	defer o.destroy()
	return o.db.Collection(o.GetCollectionName()).EstimatedDocumentCount(o.ctx, opts...)
}

func (o *orm) ReplaceOne(data interface{}, opts ...*options.ReplaceOptions) (*mgo2.UpdateResult, error) {
	defer o.destroy()
	return o.db.Collection(o.GetCollectionName()).ReplaceOne(
		o.ctx, HandleWhereData(o.model, o.where), data, opts...,
	)
}

func (o *orm) Distinct(fieldName string, opts ...*options.DistinctOptions) ([]interface{}, error) {
	defer o.destroy()
	return o.db.Collection(o.GetCollectionName()).Distinct(
		o.ctx, fieldName, HandleWhereData(o.model, o.where), opts...,
	)
}

func (o *orm) BulkWrite(models []mgo2.WriteModel, opts ...*options.BulkWriteOptions) (*mgo2.BulkWriteResult, error) {
	defer o.destroy()
	return o.db.Collection(o.GetCollectionName()).BulkWrite(o.ctx, models, opts...)
}

func (o *orm) Watch(pipeline *mongo.Pl, opts ...*options.ChangeStreamOptions) (*mgo2.ChangeStream, error) {
	defer o.destroy()
	return o.db.Collection(o.GetCollectionName()).Watch(o.ctx, pipeline.Extract(), opts...)
}

func (o *orm) Clone(opts ...*options.CollectionOptions) (*mgo2.Collection, error) {
	defer o.destroy()
	return o.db.Collection(o.GetCollectionName()).Clone(opts...)
}

func (o *orm) Drop() error {
	defer o.destroy()
	return o.db.Collection(o.GetCollectionName()).Drop(o.ctx)
}

func (o *orm) Indexes() mgo2.IndexView {
	defer o.destroy()
	return o.db.Collection(o.GetCollectionName()).Indexes()
}
