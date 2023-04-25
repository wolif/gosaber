package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
)

type Config struct {
	BrokerList   []string
	ConsumerConf *ConsumerConf
	ProducerConf *ProducerConf
}

type ProducerConf struct {
	Topic     string
	Partition int32
}

type ConsumerConf struct {
	Topics        []string
	ConsumerGroup string
	Partition     int32
	Offset        int64
}

var (
	conf map[string]*Config
	Pool map[string]*Entity
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
	if conf[name].ConsumerConf.Offset != 0 {                                      // 指针
		saramaConfig.Consumer.Offsets.Initial = conf[name].ConsumerConf.Offset
	} else {
		saramaConfig.Consumer.Offsets.Initial = -1 // 指针
	}

	saramaClient, err := sarama.NewClient(conf[name].BrokerList, saramaConfig)

	if err != nil {
		return err
	}

	if Pool == nil {
		Pool = make(map[string]*Entity)
	}
	Pool[name] = &Entity{
		ConnName:     name,
		SaramaClient: saramaClient,
		SaramaConfig: saramaConfig,
		ErrChan:      make(chan error, 1),
		MsgChan:      make(chan string, 1),
	}

	return nil
}

func GetClient(name string) (*Entity, error) {
	if c, found := Pool[name]; found {
		return c, nil
	}
	return nil, fmt.Errorf("client named [%s] not found", name)
}
