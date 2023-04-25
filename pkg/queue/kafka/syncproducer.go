package kafka

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/wolif/gosaber/pkg/queue/event"
	"github.com/wolif/gosaber/pkg/queue/event/parser"
)

func (q *Queue) SyncProduce(_ context.Context, event *event.Entity) error {
	data, err := parser.Encode(event)
	if err != nil {
		return err
	}
	return q.Kafka.SyncProduce(&sarama.ProducerMessage{
		Topic: event.Topic,
		Key:   sarama.StringEncoder(event.Key),
		Value: sarama.StringEncoder(string(data)),
	})
}
