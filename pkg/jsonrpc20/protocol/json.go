package protocol

import (
	"bytes"
	"encoding/json"
)

type Json struct{}

func (j *Json) Encode(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (j *Json) Decode(src []byte, dst interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(src))
	dec.UseNumber()
	return dec.Decode(dst)
}
