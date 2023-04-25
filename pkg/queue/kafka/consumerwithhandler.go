package kafka

import (
	"context"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/wolif/gosaber/pkg/queue/abstract"
	"github.com/wolif/gosaber/pkg/queue/event/parser"
)

func (q *Queue) ConsumeWithHandler(ctx context.Context, handler abstract.ConsumeHandler) error {
	errChan := make(chan error, 1)
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := q.Kafka.ConsumeWithHander(ctx, func(message *sarama.ConsumerMessage) (metadata string) {
			ev, err := parser.Decode(message.Value)
			if err != nil {
				handler.Err(err)
				return "event decode error"
			}
			return handler.Event(ev)
		}, func(err error) {
			handler.Err(err)
		})
		if err != nil {
			errChan <- err
			return
		}
	}()

	select {
	case <-ctx.Done():
		wg.Wait()
		return nil
	case err := <-errChan:
		return err
	}
}
