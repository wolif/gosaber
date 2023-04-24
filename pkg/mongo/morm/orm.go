package morm

import (
	"context"
	"strings"
	"time"

	pkgmongo "github.com/wolif/gosaber/pkg/mongo"
	"github.com/wolif/gosaber/pkg/mongo/madapt"
	"github.com/wolif/gosaber/pkg/ref"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type orm struct {
	client     *mongo.Client
	db         *mongo.Database
	conn       string
	coll       string
	model      Model
	where      interface{}
	projection bson.D
	sort       bson.D
	skip       int64
	limit      int64
	ctx        context.Context
	standin    *scapegoat
}

func Get(connName ...string) (*orm, error) {
	dbConnName := "default"
	if len(connName) > 0 {
		dbConnName = connName[0]
	}
	o := new(orm)
	o.conn = dbConnName
	o.destroy()
	c, e := pkgmongo.GetClient(dbConnName)
	if e != nil {
		return nil, e
	}
	o.client = c
	o.db, e = pkgmongo.GetDB(dbConnName)
	if e != nil {
		return nil, e
	}
	o.where = bson.M{}
	o.projection = bson.D{}
	o.sort = bson.D{}
	return o, nil
}

func (o *orm) Collection(col string) *orm {
	o.coll = col
	return o
}

func (o *orm) Model(model Model) *orm {
	o.model = model
	o.coll = o.model.GetCollectionName()
	return o
}

func (o *orm) Where(filter interface{}) *orm {
	match := filter
	if ref.New(filter).IsStruct() {
		match = madapt.Match(filter)
	}
	o.where = match
	return o
}

func (o *orm) Projection(fields ...string) *orm {
	o.projection = bson.D{}
	for _, field := range fields {
		o.projection = append(o.projection, bson.E{Key: field, Value: 1})
	}
	return o
}

func (o *orm) Sort(sort ...string) *orm {
	o.sort = bson.D{}
	if len(sort) > 0 {
		for _, skv := range sort {
			if strings.HasPrefix(skv, "-") {
				o.sort = append(o.sort, bson.E{
					Key:   strings.TrimPrefix(skv, "-"),
					Value: -1,
				})
			} else {
				o.sort = append(o.sort, bson.E{
					Key:   skv,
					Value: 1,
				})
			}
		}
	}
	return o
}

func (o *orm) Paginate(pagination ...int) *orm {
	page, pageSize := resolvePagination(pagination...)
	o.skip = int64((page - 1) * pageSize)
	o.limit = int64(pageSize)
	return o
}

func (o *orm) Context(ctx context.Context, timeout ...time.Duration) *orm {
	to := 0
	if len(timeout) > 0 {
		to = int(timeout[0])
	}
	if to > 0 {
		ctx, cancel := context.WithTimeout(ctx, time.Duration(to))
		go func() {
			defer cancel()
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(to)):
				return
			}
		}()
	}
	o.ctx = ctx
	return o
}

func (o *orm) destroy() {
	o.standin = nil
	o = nil
}

func (o *orm) Transaction(fn func(mongo.SessionContext) error) error {
	return o.client.UseSession(o.ctx, func(sessCtx mongo.SessionContext) error {
		defer sessCtx.EndSession(o.ctx)
		o.Context(sessCtx)
		if err := sessCtx.StartTransaction(); err != nil {
			return err
		}
		if err := fn(sessCtx); err != nil {
			if err = sessCtx.CommitTransaction(sessCtx); err != nil {
				return err
			}
		}
		if err := sessCtx.AbortTransaction(sessCtx); err != nil {
			return err
		}
		return nil
	})
}

// ---------------------------------------------------------------------------------------------------------------------
// client     *mongo.client
// db  	      *mongo.Database
// conn       string
// coll       string
// model      Model
// where      interface{}
// projection bson.D
// sort       bson.D
// skip       int64
// limit      int64
// ctx        context.Context
func (o *orm) GetClient() *mongo.Client {
	return o.client
}

func (o *orm) GetDB() *mongo.Database {
	return o.db
}

func (o *orm) GetConnName() string {
	return o.conn
}

func (o *orm) GetCollectionName() string {
	return o.coll
}

func (o *orm) GetModel() Model {
	return o.model
}

func (o *orm) GetWhere() interface{} {
	return o.where
}

func (o *orm) GetProjection() bson.D {
	return o.projection
}

func (o *orm) GetSort() bson.D {
	return o.sort
}

func (o *orm) GetSkip() int64 {
	return o.skip
}

func (o *orm) GetLimit() int64 {
	return o.limit
}

func (o *orm) GetContext() context.Context {
	return o.ctx
}

// 替身
func (o *orm) Standin() *scapegoat {
	if o.standin == nil {
		o.standin = &scapegoat{
			orm: o,
		}
	}
	return o.standin
}
