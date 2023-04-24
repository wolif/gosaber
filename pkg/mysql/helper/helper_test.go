package helper

import (
	"testing"
)

func TestMakeSortString(t *testing.T) {
	t.Log(makeSortString("id"))
	t.Log(makeSortString("-create_time"))
}