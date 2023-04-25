package kafka

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/wolif/gosaber/pkg/log"
)

type handler struct {
	Ctx     context.Context
	ErrChan chan error
	MsgChan chan string
}

func (handler) Setup(sess sarama.ConsumerGroupSession) error {
	log.Infof("session {GenerationID: %d, MemberID: %s} start", sess.GenerationID(), sess.MemberID())
	return nil
}
func (handler) Cleanup(sess sarama.ConsumerGroupSession) error {
	log.Infof("session {GenerationID: %d, MemberID: %s} stop", sess.GenerationID(), sess.MemberID())
	return nil
}
func (h handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if msg != nil {
			h.MsgChan <- string(msg.Value)
			session.MarkMessage(msg, "")
		}
	}
	return nil
}

/*
 * 使用消费组进行消费, 可以消费 多个topic 和 多个partition
 */
func (c *Entity) ConsumeWithGroup(ctx context.Context) error {
	consumerGrp, err := sarama.NewConsumerGroupFromClient(conf[c.ConnName].ConsumerConf.ConsumerGroup, c.SaramaClient)
	if err != nil {
		return err
	}
	defer consumerGrp.Close()

	err = consumerGrp.Consume(
		ctx,
		conf[c.ConnName].ConsumerConf.Topics,
		handler{Ctx: ctx, ErrChan: c.ErrChan, MsgChan: c.MsgChan},
	)
	// 连接配置有错误
	if err != nil {
		return err
	}

	for {
		err := consumerGrp.Consume(
			ctx,
			conf[c.ConnName].ConsumerConf.Topics,
			handler{Ctx: ctx, ErrChan: c.ErrChan, MsgChan: c.MsgChan},
		)
		// 连接配置有错误
		if err != nil {
			return err
		}
		if err := ctx.Err(); err != nil {
			return err
		}
	}
}

func (c *Entity) Messages() chan string {
	return c.MsgChan
}

func (c *Entity) Errors() chan error {
	return c.ErrChan
}
