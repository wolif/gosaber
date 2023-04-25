package abstract

import "github.com/wolif/gosaber/pkg/queue/event"

type Parser interface {
	Encoder
	Decoder
}

type Encoder interface {
	Encode(event *event.Event) ([]byte, error)
}

type Decoder interface {
	Decode(eventData []byte) (*event.Event, error)
}
