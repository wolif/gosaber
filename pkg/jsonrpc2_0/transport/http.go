package transport

import (
	"time"

	"github.com/wolif/gosaber/pkg/jsonrpc2_0/transport/http"
)

type Http struct {
	Url     string
	Timeout time.Duration
}

func (h *Http) Send(reqData []byte) (respData []byte, err error) {
	calling := http.New(h.Url).SetMethod(http.POST).SetHeader("context-type", "application/json")
	return nil, nil
}
