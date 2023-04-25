package kafka

import "github.com/Shopify/sarama"

type Entity struct {
	ConnName     string
	SaramaClient sarama.Client
	SaramaConfig *sarama.Config
	ErrChan      chan error
	MsgChan      chan string
}
