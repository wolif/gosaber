package structtags

import "testing"

type ts struct {
	ID int64 `search:"asdf"`
	IDs []int64 `search:"d,op: in"`
}

func TestSearch(t *testing.T) {
	ts := &ts{ID: 123}
	res, err := Parse(ts, "search", map[string]string{"addOp": "some", "omitEmpty": "true"})
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(res)
}
