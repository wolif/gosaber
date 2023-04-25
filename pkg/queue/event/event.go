package event

import (
	"github.com/wolif/gosaber/pkg/snowflake"
)

type Entity struct {
	Topic string      `json:"topic"` // @meta 事件topic 用于区分事件主题
	Type  string      `json:"type"`  // @meta 事件类型 topic的细分项
	Key   string      `json:"key"`   // @meta 对事件进行hush的依据值, 可不填, 默认为空字符串
	ID    string      `json:"id"`    // @meta 事件的唯一id, 可不填, 不填的话, 自动生成
	Data  interface{} `json:"data"`  // @data 数据
}

func New() *Entity {
	return &Entity{
		ID: snowflake.GenerateHex(),
	}
}

func (e *Entity) SetTopic(t string) *Entity {
	e.Topic = t
	return e
}

func (e *Entity) SetType(t string) *Entity {
	e.Type = t
	return e
}

func (e *Entity) SetKey(k string) *Entity {
	e.Key = k
	return e
}

func (e *Entity) SetID(id string) *Entity {
	e.ID = id
	return e
}

func (e *Entity) SetData(data interface{}) *Entity {
	e.Data = data
	return e
}
