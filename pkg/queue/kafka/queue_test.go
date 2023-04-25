package kafka

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/wolif/gosaber/pkg/kafka"
	"github.com/wolif/gosaber/pkg/log"
	"github.com/wolif/gosaber/pkg/queue/event"
)

var Q = &Queue{
	Kafka:     new(kafka.Entity),
	EventChan: make(chan *event.Entity),
	ErrorChan: make(chan error),
}

func TestQ_Start(t *testing.T) {
	log.Init(&log.Config{
		Level:         "debug",
		Format:        "mono",
		Output:        "stdout",
		RotationCount: 7,
		RotationTime:  "day",
		ServiceName:   "queue-test",
	})
	kafka.Init("default", &kafka.Config{
		BrokerList: []string{"10.20.1.20:9092"},
		ConsumerConf: &kafka.ConsumerConf{
			Topics:        []string{"ps.message_send_by_sms", "ps.message_send_by_wxsub", "ps.message_send_by_wxtpl", "ps.message_send_by_wxsub_sms", "ps.message_send_by_wxtpl_sms"},
			ConsumerGroup: "kafka-test",
			Offset:        sarama.OffsetOldest,
		},
		ProducerConf: &kafka.ProducerConf{
			Topic: "ps.message_send_by_sms",
		},
	})
	var err error
	Q.Kafka, err = kafka.GetClient("default")
	if err != nil {
		t.Error(err)
	}
}

func TestQ_SyncProduce(t *testing.T) {
	TestQ_Start(t)
	e := event.New().SetTopic("ps.message_send_by_sms").SetType("type1").SetKey("123123").SetData("asdfasdfasdfasdfasdf111222333444555666777")
	err := Q.SyncProduce(context.TODO(), e)
	if err != nil {
		t.Error(err)
	}
}

func TestQ_AsyncConsume(t *testing.T) {
	TestQ_Start(t)
	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(context.TODO())

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := Q.Consume(ctx)
		if err != nil {
			t.Error(err)
			return
		}
		select {
		case <-ctx.Done():
			return
		default:
			break
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case e := <-Q.Events():
				t.Logf("%v", e)
			case err := <-Q.Errors():
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

type testHandler struct{ t *testing.T }

func (th *testHandler) Event(event *event.Entity) string {
	th.t.Log(event)
	return ""
}
func (th *testHandler) Err(err error) {
	th.t.Error(err)
}

func TestQ_ConsumeWithHandler(t *testing.T) {
	TestQ_Start(t)
	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(context.TODO())

	handler := &testHandler{t: t}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := Q.ConsumeWithHandler(ctx, handler)
		if err != nil {
			t.Error(err)
			return
		}
		select {
		case <-ctx.Done():
			return
		default:
			break
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
