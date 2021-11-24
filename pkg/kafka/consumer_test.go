package kafka

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
)

func TestGetMessagesByConsumerGrp(t *testing.T) {
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
		err := client.ConsumeWithGroup(ctx)
		if err != nil {
			t.Error(err)
		}
	}()

	wg.Add(1)
	go func() {
		wg.Done()
		for {
			select {
			case msg := <-client.Messages():
				t.Log(msg)
			case err := <-client.Errors():
				t.Error(err)
			case <-ctx.Done():
				return
			}
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
