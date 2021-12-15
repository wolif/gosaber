package errors

import "testing"

func TestKind(t *testing.T) {
	eKind1 := NewBaseKind("kind1", 1, "error kind1 occur")
	eKind1Err1 := eKind1.Error("error kind1 error1", 11)
	eKind1Err2 := eKind1.Errorf("error %d occur", 1)
	t.Log(eKind1Err1.Code(), eKind1Err1.Error())
	t.Log(eKind1Err2.Code(), eKind1Err2.Error())

	eKind1_1 := eKind1.Extend("kind1_1", 2, "error kind1_1 occur")
	eKind1_1Err1 := eKind1_1.Error(22)
	t.Log(eKind1_1Err1.Code(), eKind1_1Err1.Error())

	t.Log(eKind1_1Err1.IsKind(eKind1_1))
	t.Log(eKind1_1Err1.IsKind(eKind1))

	eKind2 := NewBaseKind("kind2", 2, "error kind2 occur")
	t.Log(eKind1Err1.IsKind(eKind2))
	t.Log(eKind1_1Err1.IsKind(eKind2))
}
