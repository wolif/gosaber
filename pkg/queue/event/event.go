package event

import (
	"github.com/wolif/gosaber/pkg/snowflake"
)

type Event struct {
	Topic string      `json:"topic"` // @meta 事件topic 用于区分事件主题
	Type  string      `json:"type"`  // @meta 事件类型 topic的细分项
	Key   string      `json:"key"`   // @meta 事件hush的依据, 可不填, 默认为空字符串
	ID    string      `json:"id"`    // @meta 唯一id, 可不填, 不填的话, 自动生成
	Value interface{} `json:"data"`  // @data 数据
}

func New(eventType string, data interface{}) *Event {
	return NewWithID(snowflake.GenerateHex(), eventType, "", data)
}

func NewWithKey(eventType string, key string, data interface{}) *Event {
	return NewWithID(snowflake.GenerateHex(), eventType, key, data)
}

func NewWithID(id, eventType, key string, data interface{}) *Event {
	return &Event{
		Key:   key,
		ID:    id,
		Type:  eventType,
		Value: data,
	}
}
