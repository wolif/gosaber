package ref

import (
	"fmt"
	"testing"
)

type St struct {
	A int    `json:"a" map:"_a_"`
	B string `json:"b"`
}

func (st *St) Sum(params ...int) int {
	s := 0
	for _, p := range params {
		s += p
	}
	return s
}

func Test_New(t *testing.T) {
	st := &St{
		A: 1,
		B: "asdf",
	}

	refS := New(st)
	t.Log(refS.GetStructFields())
	t.Log(refS.CallStructMethod("Sum", 1, 2, 3))
	t.Log(refS.GetStructField("A"))
	t.Log(refS.GetStructFieldValue("A"))
	t.Log(refS.GetStructFieldValue("B"))
	t.Log(refS.GetStructFieldTag("A", "map"))
	t.Log(refS.GetStructFieldTag("B", "map"))

	st1 := &St{
		A: 2,
		B: "123",
	}

	fmt.Println()

	refS1 := New(st1)
	t.Log(refS1.GetStructFields())
	t.Log(refS1.CallStructMethod("Sum", 1, 2, 3))
	t.Log(refS1.GetStructField("A"))
	t.Log(refS1.GetStructFieldValue("A"))
	t.Log(refS1.GetStructFieldValue("B"))
	t.Log(refS1.GetStructFieldTag("A", "map"))
	t.Log(refS1.GetStructFieldTag("B", "map"))

	fmt.Println()

	m := map[string]interface{} {
		"a": 1,
		"b": "asdf",
	}

	refM := New(m)
	t.Log(refM.GetMapValue("a"))
	t.Log(refM.GetMapValue("b"))
	t.Log(refM.GetMapValue("c"))
}
