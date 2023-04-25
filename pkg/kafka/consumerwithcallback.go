package kafka

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/wolif/gosaber/pkg/log"
)

type callbackHandler struct {
	Ctx        context.Context
	MsgHandler func(message *sarama.ConsumerMessage) (metadata string)
}

func (callbackHandler) Setup(sess sarama.ConsumerGroupSession) error {
	log.Infof("session {GenerationID: %d, MemberID: %s} start", sess.GenerationID(), sess.MemberID())
	return nil
}
func (callbackHandler) Cleanup(sess sarama.ConsumerGroupSession) error {
	log.Infof("session {GenerationID: %d, MemberID: %s} stop", sess.GenerationID(), sess.MemberID())
	return nil
}
func (h callbackHandler) ConsumeClaim(sesssion sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if h.MsgHandler != nil && msg != nil {
			sesssion.MarkMessage(msg, h.MsgHandler(msg))
		}
	}
	return nil
}

func (c *Entity) ConsumeWithHander(ctx context.Context, msgHandler func(message *sarama.ConsumerMessage) (metadata string), errHandler func(err error)) error {
	consumerGrp, err := sarama.NewConsumerGroupFromClient(conf[c.ConnName].ConsumerConf.ConsumerGroup, c.SaramaClient)
	if err != nil {
		return err
	}
	defer consumerGrp.Close()

	for {
		err := consumerGrp.Consume(
			ctx,
			conf[c.ConnName].ConsumerConf.Topics,
			callbackHandler{Ctx: ctx, MsgHandler: msgHandler},
		)

		// 连接配置有错误
		if err != nil {
			return err
		}

		if err := ctx.Err(); err != nil {
			if errHandler != nil {
				errHandler(err)
			}
		}
	}
}
