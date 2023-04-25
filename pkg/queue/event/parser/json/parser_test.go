package json

import (
	"testing"

	event2 "github.com/wolif/gosaber/pkg/queue/event"
	"github.com/wolif/gosaber/pkg/snowflake"
)

var parser = new(jsonParser)

func start() {
	snowflake.Init(&snowflake.Config{})
}

func TestJsonParser_Encode(t *testing.T) {
	start()
	data, err := parser.Encode(event2.New().SetType("type1").SetData(123))
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(data))
	}
}

func TestJsonParser_Decode(t *testing.T) {
	data := `{"id":"054c4efdf3800000","type":"type1","data":123}`
	e, err := parser.Decode([]byte(data))
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%v", e)
	}
}
