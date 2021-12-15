package errors1

import "testing"

func TestMakeError(t *testing.T) {
	k := NewKind("server error", 10000, "server error occur")
	e1 := k.NewError()
	t.Log(e1)
	t1 := k.NewType("method not found", 10001, "method [] not found")
	e2 := t1.NewError()
	t.Log(e2)
	t2 := k.NewType("invalid args", 10002, "invalid arguments")
	e3 := t2.NewError("invalid argsssssssssssssssssssssss")
	t.Log(e3)

	t.Log(e1.IsKind(k))
	t.Log(e2.IsKind(k))
	t.Log(e3.IsKind(k))

	t.Log(e2.IsType(t1))
	t.Log(e2.IsType(t2))
}
