package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/wolif/gosaber/pkg/mongo"
	"testing"
	"time"
)

func Init() {
	err := mongo.Init(map[string]*mongo.Config{
		"default": {
			URI:    "mongodb://g1-mongodb-m.st1.tgs.com:27017/db_towngas_message_center",
			DBName: "db_towngas_message_center",
		},
	})
	fmt.Println(err)
}

type RecordTest struct {
	ID        string `bson:"_id"`
	Name      string `bson:"name"`
	Time      int64  `bson:"time"`
	CreatedAt int64  `bson:"created_at" mongo:"autoCreateTime"`
	UpdatedAt int64  `bson:"updated_at"`
	DAt       int64  `bson:"deleted_at" mongo:"autoDeleteTime"`
}

type RecordTestSearch struct {
	ID   string   `search:"_id" addOpt:"MongoStrID2ObjID"`
	IDs  []string `search:"_id" addOpt:"MongoStrID2ObjID"`
	Name string   `search:"name"`
}

func (rt *RecordTest) GetCollectionName() string {
	return "t_record_test"
}

func TestStructToJson(t *testing.T) {
	r := &RecordTest{
		Name: "test1",
		Time: time.Now().Unix(),
	}
	t.Log(structToBson(r))
}

func TestMapToBson(t *testing.T) {
	m := map[string]interface{}{
		"abc": 1,
		"def": "123",
	}
	t.Log(mapToBson(m))
}

func TestDataToModelBson(t *testing.T) {
	r := &RecordTest{
		Name: "test1",
		Time: time.Now().Unix(),
	}
	t.Log(dataToModelBson(r, true))
	m := map[string]interface{}{
		"abc": 1,
		"def": "123",
	}
	t.Log(dataToModelBson(m, true))
}

func TestFetch3TimeFields(t *testing.T) {
	r := &RecordTest{
		Name: "test1",
		Time: time.Now().Unix(),
	}
	t.Log(fetch3TimeFields(r))
}

func TestRecord_InsertOne(t *testing.T) {
	Init()
	id, err := Record.InsertOne(context.TODO(), &RecordTest{
		Name: "test1",
		Time: time.Now().Unix(),
	})

	if err != nil {
		t.Error(err)
	} else {
		t.Log(id)
	}
}

func TestRecord_InsertMany(t *testing.T) {
	Init()
	ids, err := Record.InsertMany(context.TODO(), &RecordTest{
		Name: "test1",
		Time: time.Now().Unix(),
	}, &RecordTest{
		Name: "test2",
		Time: time.Now().Unix(),
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(ids)
	}
}

func TestRecord_FindOne(t *testing.T) {
	Init()
	result := new(RecordTest)
	err := Record.FindOne(context.TODO(), &RecordTest{}, RecordTestSearch{
		Name: "test1",
		//ID: "6195b4e530470f13f84bee3b",
	}, result)

	if err != nil {
		t.Error(err)
	} else {
		s, _ := json.Marshal(result)
		t.Log(string(s))
	}
}

func TestRecord_SearchRecord(t *testing.T) {
	Init()
	result := make([]*RecordTest, 0)
	count, err := Record.SearchRecord(context.TODO(), &RecordTest{}, RecordTestSearch{
		IDs:  []string{"619c5280d5bb8b178ef69ca6", "619c52940de10742b2d699b6"},
	}, nil, &result)

	if err != nil {
		t.Error(err)
	} else {
		t.Log(count)
		s, _ := json.Marshal(result)
		t.Log(string(s))
	}
}

func TestRecord_UpdateOne(t *testing.T) {
	Init()
	cnt, err := Record.UpdateOne(context.TODO(), &RecordTest{}, RecordTestSearch{
		Name: "test2",
	}, bson.M{"name": "test2-update"})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(cnt)
	}
}

func TestRecord_UpdateMany(t *testing.T) {
	Init()
	cnt, err := Record.UpdateMany(context.TODO(), &RecordTest{}, RecordTestSearch{
		Name: "test2",
	}, bson.M{"name": "test2-update"})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(cnt)
	}
}

func TestRecord_Delete(t *testing.T) {
	Init()
	cnt, err := Record.Delete(context.TODO(), &RecordTest{}, RecordTestSearch{
		Name: "test1",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(cnt)
	}
}
