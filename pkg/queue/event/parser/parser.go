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

func Encode(e *event.Entity) ([]byte, error) {
	return parser.Encode(e)
}

func Decode(data []byte) (*event.Entity, error) {
	return parser.Decode(data)
}
