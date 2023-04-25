package json

import (
	"encoding/json"

	"github.com/wolif/gosaber/pkg/queue/event"
)

type jsonParser struct{}

var JsonParser = new(jsonParser)

func (j *jsonParser) Encode(e *event.Entity) ([]byte, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (j *jsonParser) Decode(data []byte) (*event.Entity, error) {
	e := new(event.Entity)
	err := json.Unmarshal(data, e)
	if err != nil {
		return nil, err
	}
	return e, nil
}
