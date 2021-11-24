package json

import (
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func ToJson(obj interface{}) (string, error) {
	bs, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func ToJsonIgnoreError(obj interface{}) string {
	jsonStr, _ := ToJson(obj)
	return jsonStr
}

func FromJson(jsonStr string, obj interface{}) error {
	return json.Unmarshal([]byte(jsonStr), obj)
}

func FromJsonBytes(bytes []byte, obj interface{}) error {
	return json.Unmarshal(bytes, obj)
}

func FromJsonIgnoreError(jsonStr string, obj interface{}) {
	_ = FromJson(jsonStr, obj)
}

func Copy(dest interface{}, src interface{}) error {
	bs, err := json.Marshal(src)
	if err != nil {
		return err
	}

	return FromJsonBytes(bs, dest)
}
