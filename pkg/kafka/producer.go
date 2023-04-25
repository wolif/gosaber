package kafka

import (
	"github.com/Shopify/sarama"
)

func (c *Entity) SyncProduce(message *sarama.ProducerMessage) error {
	producer, err := sarama.NewSyncProducerFromClient(c.SaramaClient)
	if err != nil {
		return err
	}

	defer producer.Close()

	if message.Topic == "" {
		message.Topic = conf[c.ConnName].ProducerConf.Topic
	}

	_, _, err = producer.SendMessage(message)
	if err != nil {
		return err
	}
	return nil
}
