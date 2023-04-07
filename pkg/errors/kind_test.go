package errors

import (
	"testing"
)

func TestKind(t *testing.T) {
	eKind1 := NewBaseKind("kind1", 1, "error kind1 occur")
	eKind1Err1 := eKind1.New("error kind1 error1", 11)
	eKind1Err2 := eKind1.NewF("error %d occur", 1)
	t.Log(eKind1Err1.Code(), eKind1Err1.Error())
	t.Log(eKind1Err2.Code(), eKind1Err2.Error())

	eKind1_1 := eKind1.Extend("kind1_1", 2, "error kind1_1 occur")
	eKind1_1Err1 := eKind1_1.New(22)
	t.Log(eKind1_1Err1.Code(), eKind1_1Err1.Error())

	t.Log(eKind1_1Err1.IsA(eKind1_1))
	t.Log(eKind1_1Err1.IsA(eKind1))
	t.Log(eKind1_1Err1.IsA(eKind1, true))

	eKind2 := NewBaseKind("kind2", 2, "error kind2 occur")
	t.Log(eKind1Err1.IsA(eKind2))
	t.Log(eKind1_1Err1.IsA(eKind2))
}
