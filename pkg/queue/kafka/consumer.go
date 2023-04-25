package kafka

import (
	"context"
	"sync"

	"github.com/wolif/gosaber/pkg/queue/event"
	"github.com/wolif/gosaber/pkg/queue/event/parser"
)

/*
 * 仅支持以 consumer group 的方式获取kafka的数据
 */
func (q *Queue) Consume(ctx context.Context) error {
	errChan := make(chan error, 1)
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := q.Kafka.ConsumeWithGroup(ctx)
		if err != nil {
			errChan <- err
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case message := <-q.Kafka.Messages():
				e, err := parser.Decode([]byte(message))
				if err == nil {
					q.EventChan <- e
				}
			case <-ctx.Done():
				return
			}
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

func (q *Queue) Events() chan *event.Entity {
	return q.EventChan
}

func (q *Queue) Errors() chan error {
	return q.Kafka.Errors()
}
