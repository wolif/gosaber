package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/wolif/gosaber/pkg/log"
	"time"
)

type grpHandler struct {
	Ctx     context.Context
	ErrChan chan error
	MsgChan chan string
}

func (grpHandler) Setup(sess sarama.ConsumerGroupSession) error {
	log.Infof("session {GenerationID: %d, MemberID: %s} start", sess.GenerationID(), sess.MemberID())
	return nil
}
func (grpHandler) Cleanup(sess sarama.ConsumerGroupSession) error {
	log.Infof("session {GenerationID: %d, MemberID: %s} stop", sess.GenerationID(), sess.MemberID())
	return nil
}
func (h grpHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if msg != nil {
			h.MsgChan <- string(msg.Value)
			sess.MarkMessage(msg, "")
		}
	}
}

/*
 * 使用消费组进行消费, 可以消费 多个topic 和 多个partition
 */
func (c *Client) ConsumeWithGroup(ctx context.Context) error {
	consumerGrp, err := sarama.NewConsumerGroupFromClient(conf[c.ConnName].Consumer.ConsumerGroup, c.SaramaClient)
	if err != nil {
		return err
	}
	defer func() {
		if err := consumerGrp.Close(); err != nil {
			panic(err)
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
				grpHandler{Ctx: ctx, ErrChan: c.ErrChan, MsgChan: c.MsgChan},
			)

			// 连接配置有错误
			if err != nil {
				return err
			}

			if err := ctx.Err(); err != nil {
				select {
				case c.ErrChan <- err:
				case <-time.After(1 * time.Second):
				}
			}
		}
	}
}

func (c *Client) Messages() chan string {
	return c.MsgChan
}

func (c *Client) Errors() chan error {
	return c.ErrChan
}
