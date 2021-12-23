package dotenv

import (
	"encoding/json"
	"testing"
)

type TestConf struct {
	Test *Test
}

type Test struct {
	Attr *Attr
}

type Attr struct {
	I int64
	S string
	S2 string
	B bool
	A []int64
	N json.RawMessage
}

func TestLoad(t *testing.T) {
	tc := new(TestConf)
	err := Load(".env", tc)
	if err != nil {
		t.Error(err)
	}
	jsonData, err := json.Marshal(tc)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%s", jsonData)
}
