package abstract

import (
	"context"

	"github.com/wolif/gosaber/pkg/queue/event"
)

type Consumer interface {
	// 获取events 和 errors 前调用
	Consume(ctx context.Context) error

	// 获取events
	Events() chan *event.Entity

	// 获取error
	Errors() chan error
}
