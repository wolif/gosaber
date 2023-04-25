package queue

import (
	"fmt"

	kafka2 "github.com/wolif/gosaber/pkg/kafka"
	"github.com/wolif/gosaber/pkg/queue/abstract"
	"github.com/wolif/gosaber/pkg/queue/event"
	"github.com/wolif/gosaber/pkg/queue/kafka"
)

type Config struct {
	Type string
}

var QueuePool = make(map[string]abstract.QueueInterface)

func Init(name string, c *Config) error {
	switch c.Type {
	case "kafka":
		fallthrough
	default:
		client, err := kafka2.GetClient(name)
		if err != nil {
			return err
		}
		QueuePool[c.Type+"_"+name] = &kafka.Queue{
			Kafka:     client,
			EventChan: make(chan *event.Entity, 1),
			ErrorChan: make(chan error, 1),
		}
	}
	return nil
}

func GetQueue(typo, name string) (abstract.QueueInterface, error) {
	if q, found := QueuePool[typo+"_"+name]; found {
		return q, nil
	}
	return nil, fmt.Errorf("QueueInterface with type [%s] and name [%s] not found", typo, name)
}
