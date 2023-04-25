package json

import (
	"encoding/json"

	"github.com/wolif/gosaber/pkg/queue/event"
)

type jsonParser struct{}

var JsonParser = new(jsonParser)

func (j *jsonParser) Encode(event *event.Event) ([]byte, error) {
	data, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (j *jsonParser) Decode(eventData []byte) (*event.Event, error) {
	event := new(event.Event)
	err := json.Unmarshal(eventData, event)
	if err != nil {
		return nil, err
	}
	return event, nil
}
