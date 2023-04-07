package crypto

import (
	"testing"
)

type testStruct struct {
	StrVal string
}

func (ts testStruct) Serialize() []byte {
	return []byte(ts.StrVal)
}

func TestMD5(t *testing.T) {
	t.Log(MD5(1))
	t.Log(MD5(0.1))
	t.Log(MD5("abcd"))
	t.Log(MD5(&testStruct{"abcd"}))
}
