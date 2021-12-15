package errors

import "testing"

func TestKind(t *testing.T) {
	serverErr := NewKind("server error", 10000, "server error occur")
	methodNotFoundErr := serverErr.Extend("method not found error", 10001, "method [] not found")

	se1 := serverErr.Error()
	// se2 := serverErr.NewError(10010, "se2 error")

	// t.Log(se1.Code(), se1.Error())
	// t.Log(se2.Code(), se2.Error())

	mnfe := methodNotFoundErr.Error("method mnfe not found")
	// t.Log(mnfe.Code(), mnfe.Error())

	t.Log(mnfe.IsKind(serverErr))
	t.Log(mnfe.IsKind(methodNotFoundErr))
	t.Log(se1.IsKind(methodNotFoundErr))
}
