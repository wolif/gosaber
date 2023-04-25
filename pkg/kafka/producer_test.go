package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/wolif/gosaber/pkg/log"
	"testing"
)

func start() {
	Init("default", &Config{
		BrokerList: []string{"10.20.1.20:9092"},
		Consumer: &Consumer{
			Topics:        []string{"ps.message_send_by_sms"},
			ConsumerGroup: "kafka-test",
			Offset:        -2,
		},
		Producer: &Producer{
			Topic: "ps.message_send_by_sms",
		},
	})

	log.Init(&log.Config{
		Level:         "debug",
		Format:        "mono",
		Output:        "stdout",
		RotationCount: 7,
		RotationTime:  "day",
		ServiceName:   "test",
	})
}

func TestSyncProduce(t *testing.T) {
	start()
	client, err := GetClient("default")
	if err != nil {
		t.Error(err)
	} else {
		x := 1
		for i := x; i <= x+5; i++ {
			err = client.SyncProduce(&sarama.ProducerMessage{
				Value: sarama.StringEncoder(fmt.Sprint(x)),
			})
			if err != nil {
				t.Error(err)
			}
		}
	}
}
