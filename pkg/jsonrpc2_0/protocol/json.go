package protocol

import "encoding/json"

type Json struct{}

func (j *Json) Encode(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (j *Json) Decode(src []byte, dst interface{}) error {
	return json.Unmarshal(src, dst)
}
