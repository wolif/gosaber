package parser

import (
	"github.com/wolif/gosaber/pkg/queue/event"
	"github.com/wolif/gosaber/pkg/queue/event/parser/abstract"
	"github.com/wolif/gosaber/pkg/queue/event/parser/json"
)

var parser abstract.Parser = json.JsonParser

func SetParser(p abstract.Parser) {
	parser = p
}

func Encode(event *event.Event) ([]byte, error) {
	return parser.Encode(event)
}

func Decode(eventData []byte) (*event.Event, error) {
	return parser.Decode(eventData)
}
