package where

import "testing"

type searchCond struct {
	ID        int64   `search:"id"`
	IDs       []int64 `search:"id"`
	IDsNotIn  []int64 `search:"idNIN,op:not in"`
	Name      string  `search:"name,op:like"`
	Title     string  `search:"title,op:like abc%"`
	Age       int     `search:"age"`
	ArrayElem int     `search:"set,op:find in set"`
	Time      int64   `search:"time_at,op:>=,addOp:TimestampIntToStr"`
	Mark      string  `search:"mark,op:sql"`
}

func TestMake(t *testing.T) {
	sc := &searchCond{
		IDs:       []int64{2},
		IDsNotIn:  []int64{3},
		Name:      "abcdef",
		Title:     "title",
		Age:       10,
		ArrayElem: 9,
		Time:      1000000000,
		Mark:      "mark <> 1",
	}

	t.Log(Make(sc))
}
