package kafka

import (
	"github.com/wolif/gosaber/pkg/kafka"
	"github.com/wolif/gosaber/pkg/queue/event"
)

type Queue struct {
	Client    *kafka.Client
	EventChan chan *event.Event
	ErrorChan chan error
}
