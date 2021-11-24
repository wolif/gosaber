package strs

import "testing"

func TestUcFirst(t *testing.T) {
	t.Log(UcFirst("123232abc"))
	t.Log(UcFirst("example"))
	t.Log(UcFirst("Example"))
}

func TestLcFirst(t *testing.T) {
	t.Log(LcFirst("123232abc"))
	t.Log(LcFirst("example"))
	t.Log(LcFirst("Example"))
}

func TestCamelString(t *testing.T) {
	t.Log(CamelString("created_at"))
}
