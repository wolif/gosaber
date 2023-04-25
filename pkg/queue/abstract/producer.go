package abstract

import (
	"context"

	"github.com/wolif/gosaber/pkg/queue/event"
)

type SyncProducer interface {
	SyncProduce(ctx context.Context, event *event.Event) error
}
