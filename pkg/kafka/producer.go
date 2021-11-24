package kafka

import (
	"github.com/Shopify/sarama"
)

func (c *Client) SyncProduce(message *sarama.ProducerMessage) error {
	producer, err := sarama.NewSyncProducerFromClient(c.SaramaClient)
	if err != nil {
		return err
	}

	defer func() {
		if err := producer.Close(); err != nil {
			panic(err)
		}
	}()

	if message.Topic == "" {
		message.Topic = conf[c.ConnName].Producer.Topic
	}

	_, _, err = producer.SendMessage(message)
	if err != nil {
		return err
	}
	return nil
}
