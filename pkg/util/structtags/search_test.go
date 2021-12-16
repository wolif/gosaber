package structtags

import "testing"

type testS struct {
	id int64 `search:"id"`
}

func TestSearch(t *testing.T) {
	sp := NewSearchParser("search", "=", "in")
	r := sp.Parse(&testS{id: 1})
	t.Log(r)
}
