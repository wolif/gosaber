package abstract

import (
	"context"

	"github.com/wolif/gosaber/pkg/queue/event"
)

type ConsumeHandler interface {
	Event(event *event.Event) string
	Err(err error)
}

type ConsumerWithHandler interface {
	ConsumeWithHandler(ctx context.Context, handler ConsumeHandler) error
}
