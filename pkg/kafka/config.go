package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
)

type Config struct {
	BrokerList []string
	Consumer   *Consumer
	Producer   *Producer
}

type Producer struct {
	Topic     string
	Partition int32
}

type Consumer struct {
	Topics        []string
	ConsumerGroup string
	Partition     int32
	Offset        int64
}

var (
	conf       map[string]*Config
	ClientPool map[string]*Client
)

func Init(name string, c *Config) error {
	if conf == nil {
		conf = make(map[string]*Config)
	}
	conf[name] = c

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true                                 // 生产成功返回
	saramaConfig.Producer.Return.Errors = true                                    // 生产时有错误报出
	saramaConfig.Consumer.Return.Errors = true                                    // 消费时有错误报出
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky // 重新分配策略
	if conf[name].Consumer.Offset != 0 {                                          // 指针
		saramaConfig.Consumer.Offsets.Initial = conf[name].Consumer.Offset
	} else {
		saramaConfig.Consumer.Offsets.Initial = -1 // 指针
	}

	saramaClient, err := sarama.NewClient(conf[name].BrokerList, saramaConfig)

	if err != nil {
		return err
	}

	if ClientPool == nil {
		ClientPool = make(map[string]*Client)
	}
	ClientPool[name] = &Client{
		ConnName:     name,
		SaramaClient: saramaClient,
		SaramaConfig: saramaConfig,
		ErrChan:      make(chan error, 1),
		MsgChan:      make(chan string, 1),
	}

	return nil
}

func GetClient(name string) (*Client, error) {
	if c, found := ClientPool[name]; found {
		return c, nil
	}
	return nil, fmt.Errorf("client named [%s] not found", name)
}
