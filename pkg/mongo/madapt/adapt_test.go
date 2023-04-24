package madapt

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

type cond struct {
	O         int      `search:"status,op:gt,omitEmpty:f"`
	A         string   `search:"field"`
	B         string   `search:"_id,addOp:StrIDToObjID"`
	C         []string `search:"id,addOp:StrIDToObjID"`
	D         string   `search:"time1,op:gte,addOp:StrToTime"`
	E         string   `search:"title,op:regex"`
	F         string   `search:"op:regex"`
	Ids       []int    `search:"IDS"`
	TimeStart string   `search:"time,op:gte,addOp:StrToTime"`
	TimeEnd   string   `search:"time,op:lte,addOp:StrToTime"`
	RawField  bson.M   `search:"one_filed,op:raw"`
}

func TestMatch(t *testing.T) {
	c := &cond{
		O:         0,
		A:         "abcd",
		B:         "6087b4f5dd549f95b4cc8cfe",
		C:         []string{"6087b4f5dd549f95b4cc8cfe"},
		D:         "2021-11-02 10:10:10",
		E:         "title",
		F:         "titleF",
		Ids:       []int{1, 2, 3},
		TimeStart: "2021-11-01 10:10:10",
		TimeEnd:   "2021-11-30 10:10:10",
		RawField:  bson.M{"$push": "time"},
	}

	res := Match(c)
	t.Log(res)
}

func TestSort(t *testing.T) {
	t.Log(Sort("-create_time", "-id", "status"))
}
