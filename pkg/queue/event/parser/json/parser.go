package json

import (
	"bytes"
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
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(e); err != nil {
		return nil, err
	}
	return e, nil
}
