package morm

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/wolif/gosaber/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
)

func Init() {
	err := mongo.Init(map[string]*mongo.Config{
		"default": {
			URI:         "mongodb://10.20.1.20:27017/db_towngas_message_center&retryWrites=false",
			DBName:      "db_towngas_message_center",
			EnableTrace: true,
		},
	})
	fmt.Println(err)
}

type RecordTest struct {
	ID        string `bson:"_id"`
	Name      string `bson:"name"`
	CreatedAt int64  `bson:"created_at" morm:"create_time,time_format:20060102150405"`
	UpdatedAt int64  `bson:"updated_at" morm:"update_time,time_format:20060102150405"`
	DAt       int64  `bson:"deleted_at" morm:"delete_time,time_format:20060102150405"`
}

type RecordTestSearch struct {
	ID       string   `search:"_id,addOp:StrIDToObjID"`
	IDs      []string `search:"_id,addOp:StrIDToObjID"`
	Name     string   `search:"name"`
	NameLike string   `search:"name,op:regex"`
	Names    []string `search:"name"`
}

func (rt *RecordTest) GetCollectionName() string {
	return "t_record_test"
}

func TestStructToJson(t *testing.T) {
	r := &RecordTest{
		Name: "test1",
	}
	t.Log(Bson.FromStruct(r))
}

func TestMapToBson(t *testing.T) {
	m := map[string]interface{}{
		"abc": 1,
		"def": "123",
	}
	t.Log(Bson.FromMap(m))
}

func TestDataToModelBson(t *testing.T) {
	r := &RecordTest{
		Name: "test1",
	}
	t.Log(Bson.ToModel(r, true))
	m := map[string]interface{}{
		"abc": 1,
		"def": "123",
	}
	t.Log(Bson.ToModel(m, true))
}

func TestResolveModelOpt(t *testing.T) {
	opts := resolveModelOpts(&RecordTest{})
	t.Log(opts.createTime, opts.updateTime, opts.deleteTime)
}

// ---------------------------------------------------------------------------------------------------------------------

func TestOrm_Insert(t *testing.T) {
	Init()
	orm, err := Get()
	if err != nil {
		t.Error(err)
		return
	}
	id, err := orm.Standin().Insert(&RecordTest{
		Name: "test5",
	})

	if err != nil {
		t.Error(err)
	} else {
		t.Log(id)
	}
}

func TestOrm_Insert2(t *testing.T) {
	Init()
	orm, err := Get()
	if err != nil {
		t.Error(err)
		return
	}

	err = orm.Transaction(func(sCtx mongo2.SessionContext) error {
		_, err := orm.Standin().Insert(&RecordTest{
			ID:   "aef2ea61554e82a9582b78",
			Name: "test_duplicate",
		})
		if err != nil {
			return err
		}
		_, err = orm.Standin().Insert(&RecordTest{
			Name: "test1",
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

func TestOrm_Inserts(t *testing.T) {
	Init()
	orm, err := Get()
	if err != nil {
		t.Error(err)
		return
	}
	ids, err := orm.Standin().Inserts([]Model{
		&RecordTest{Name: "test2"},
		&RecordTest{Name: "test3"},
		&RecordTest{Name: "test4"},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(ids)
	}
}

func TestOrm_FindOne(t *testing.T) {
	Init()
	orm, err := Get()
	if err != nil {
		t.Error(err)
		return
	}
	result := new(RecordTest)
	err = orm.Model(&RecordTest{}).Where(RecordTestSearch{ID: "63f889d2b414f3526da318dd"}).Standin().FindOne(result)

	if err != nil {
		t.Error(err)
	} else {
		t.Log(result)
		s, _ := json.Marshal(result)
		t.Log(string(s))
	}
}

func TestOrm_Count(t *testing.T) {
	Init()
	orm, err := Get()
	if err != nil {
		t.Error(err)
		return
	}
	count, err := orm.Model(&RecordTest{}).Where(bson.M{"name": bson.M{"$regex": "test"}}).Standin().Count()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(count)
	}
}

func TestOrm_Search(t *testing.T) {
	Init()
	// span := lighttracer.GlobalTracer().StartSpan("mongo-test")
	orm, err := Get()
	// orm.Context(lighttracer.ContextWithSpan(context.TODO(), span))
	orm.Context(context.TODO())
	if err != nil {
		t.Error(err)
		return
	}
	result := make([]*RecordTest, 0)
	count, err := orm.Model(&RecordTest{}).
		//Where(bson.M{"name": bson.M{"$regex": "test"}}).
		//Where(RecordTestSearch{NameLike: "test"}).
		Sort("created_at", "name").
		Standin().Find(&result)

	if err != nil {
		t.Error(err)
	} else {
		t.Log(count)
		for _, r := range result {
			s, _ := json.Marshal(r)
			t.Log(string(s))
		}
	}

	// span.Finish()
}

func TestOrm_Set(t *testing.T) {
	Init()
	orm, err := Get()
	if err != nil {
		t.Error(err)
		return
	}
	cnt, err := orm.Model(&RecordTest{}).
		Where(RecordTestSearch{Name: "test2"}).
		Standin().Set(bson.M{"name": "test2-update"})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(cnt)
	}
}

func TestOrm_Sets(t *testing.T) {
	Init()
	orm, err := Get()
	if err != nil {
		t.Error(err)
		return
	}
	cnt, err := orm.Model(&RecordTest{}).
		Where(RecordTestSearch{Names: []string{"test2", "test3"}}).
		Standin().Sets(bson.M{"name": "test-update"})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(cnt)
	}
}

func TestOrm_Delete(t *testing.T) {
	Init()
	orm, err := Get()
	if err != nil {
		t.Error(err)
		return
	}
	cnt, err := orm.Model(&RecordTest{}).Where(RecordTestSearch{Name: "test1"}).Standin().Delete()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(cnt)
	}
}

func TestOrm_Aggregate(t *testing.T) {
	Init()
	orm, err := Get()
	if err != nil {
		t.Error(err)
		return
	}

	data1 := make([]*struct {
		ID   string `json:"id" bson:"_id"`
		Name string `json:"name" bson:"name"`
	}, 0)
	//data2 := make([]map[string]interface{}, 0)
	cur, err := orm.Model(&RecordTest{}).Aggregate(mongo.Pipeline.New().
		Group(bson.M{
			"_id":  "$name",
			"name": bson.M{"$first": "$name"},
		}).Extract(),
	)
	if err != nil {
		t.Error(err)
		return
	}
	err = cur.All(context.TODO(), &data1)
	if err != nil {
		t.Error(err)
		return
	}
	s1, _ := json.Marshal(data1)
	t.Log(string(s1))
}
