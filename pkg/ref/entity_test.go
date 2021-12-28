package ref

import (
	"fmt"
	"testing"
)

type ts struct {
	ID int64 `search:"id"`
}

func (ts *ts) Show() {
	fmt.Printf("id = %d", ts.ID)
}

func TestEntity(t *testing.T) {
	ts := &ts{ID: 123}
	refTs := New(ts)
	t.Log(refTs.IsString())
	t.Log(refTs.IsNumber())
	t.Log(refTs.IsStruct())

	i := 123
	refI := New(i)
	t.Log(!refI.IsStringOrNumber() && !refI.IsSlice())

	refTs.StructMethodCall("Show")
}
