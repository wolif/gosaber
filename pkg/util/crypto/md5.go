package crypto

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"reflect"
)

type Serializable interface {
	Serialize() []byte
}

func MD5(i interface{}) string {
	if i == nil {
		return ""
	}
	switch it := i.(type) {
	case []byte:
		return fmt.Sprintf("%x", md5.Sum(it))
	case string:
		return MD5([]byte(it))
	case int, int8, int16, int32, int64, float32, float64:
		return MD5(fmt.Sprint(it))
	case Serializable:
		return MD5(it.Serialize())
	}
	kind := reflect.TypeOf(i).Kind()
	if kind == reflect.Map || kind == reflect.Array {
		s, _ := json.Marshal(i)
		return MD5(s)
	}
	return ""
}
