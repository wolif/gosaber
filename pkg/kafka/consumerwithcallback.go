package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/wolif/gosaber/pkg/log"
)

type cbGrpHandler struct {
	Ctx     context.Context
	MsgFunc func(message *sarama.ConsumerMessage) (metadata string)
}

func (cbGrpHandler) Setup(sess sarama.ConsumerGroupSession) error {
	log.Infof("session {GenerationID: %d, MemberID: %s} start", sess.GenerationID(), sess.MemberID())
	return nil
}
func (cbGrpHandler) Cleanup(sess sarama.ConsumerGroupSession) error {
	log.Infof("session {GenerationID: %d, MemberID: %s} stop", sess.GenerationID(), sess.MemberID())
	return nil
}
func (h cbGrpHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if h.MsgFunc != nil && msg != nil {
			metadata := ""
			metadata = h.MsgFunc(msg)
			sess.MarkMessage(msg, metadata)
		}
	}
}

func (c *Client) ConsumeWithGroupAndCbFunc(ctx context.Context, msgFunc func(message *sarama.ConsumerMessage) (metadata string), errFunc func(err error)) error {
	consumerGrp, err := sarama.NewConsumerGroupFromClient(conf[c.ConnName].Consumer.ConsumerGroup, c.SaramaClient)
	if err != nil {
		return err
	}
	defer func() {
		if err := consumerGrp.Close(); err != nil {
			log.Fatalf("error occur when close consumer group: %v", err)
		}
	}()

	// 连接断开时尝试重新连接
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			err := consumerGrp.Consume(
				ctx,
				conf[c.ConnName].Consumer.Topics,
				cbGrpHandler{Ctx: ctx, MsgFunc:msgFunc},
			)

			// 连接配置有错误
			if err != nil {
				return err
			}

			if err := ctx.Err(); err != nil {
					if errFunc != nil {
					errFunc(err)
				}
			}
		}
	}
}
