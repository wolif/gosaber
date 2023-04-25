package kafka

import (
	"github.com/wolif/gosaber/pkg/kafka"
	"github.com/wolif/gosaber/pkg/queue/event"
)

type Queue struct {
	Kafka     *kafka.Entity
	EventChan chan *event.Entity
	ErrorChan chan error
}
