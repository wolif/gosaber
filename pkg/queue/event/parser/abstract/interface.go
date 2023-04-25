package abstract

import "github.com/wolif/gosaber/pkg/queue/event"

type Parser interface {
	Encoder
	Decoder
}

type Encoder interface {
	Encode(e *event.Entity) ([]byte, error)
}

type Decoder interface {
	Decode(data []byte) (*event.Entity, error)
}
