package mongoAdapt

import "testing"

type cond struct {
	O int      `search:"status" opt:"gt" omitEmpty:"f"`
	A string   `search:"field"`
	B string   `search:"_id" addOpt:"MongoStrID2ObjID"`
	C []string `search:"id" addOpt:"MongoStrID2ObjID"`
	D string   `search:"time" opt:"gte" addOpt:"TimestampStrToInt64"`
	E string   `search:"title" opt:"regex"`
}

func TestMatch(t *testing.T) {
	c := &cond{
		O: 0,
		A: "abcd",
		B: "6087b4f5dd549f95b4cc8cfe",
		C: []string{"6087b4f5dd549f95b4cc8cfe"},
		D: "2021-11-02 10:10:10",
		E: "title",
	}

	res := Match(c)
	t.Log(res)
}
