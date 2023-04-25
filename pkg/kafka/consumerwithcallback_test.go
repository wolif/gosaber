package kafka

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"

	"github.com/Shopify/sarama"
)

func TestClient_ConsumeWithGroupAndCbFunc(t *testing.T) {
	start()
	client, err := GetClient("default")
	if err != nil {
		t.Error(err)
		return
	}

	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(context.TODO())

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := client.ConsumeWithHander(ctx, func(message *sarama.ConsumerMessage) (metadata string) {
			t.Logf("partition: %d, value: %s", message.Partition, string(message.Value))
			return
		}, func(err error) {
			t.Error(err)
		})
		if err != nil {
			t.Error(err)
		}
	}()

	// 等待系统信号
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-signals
		t.Logf("worker get a signal %s", s.String())
		switch s {
		case syscall.SIGHUP:
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			cancel()
			wg.Wait()
			t.Log("worker gorutines quited, now processor quit")
			return
		}
	}
}
