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
	event := event2.New("type1", 123)
	data, err := parser.Encode(event)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(data))
	}
}

func TestJsonParser_Decode(t *testing.T) {
	eventData := `{"id":"054c4efdf3800000","type":"type1","data":123}`
	event, err := parser.Decode([]byte(eventData))
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%v", event)
	}
}
