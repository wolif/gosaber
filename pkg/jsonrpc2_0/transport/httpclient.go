package transport

import (
	"context"
	"fmt"
	"time"

	"github.com/wolif/gosaber/pkg/jsonrpc2_0/transport/http"
)

type HttpClient struct {
	Url     string
	Timeout time.Duration
}

func (h *HttpClient) Send(ctx context.Context, reqData []byte, header ...interface{}) (respData []byte, err error) {
	calling := http.New(h.Url).
		SetMethod(http.POST).
		SetHeader("context-type", "application/json").
		SetBody(reqData).
		SetContext(ctx).
		SetOption(http.NewOption().SetTimeout(h.Timeout))
	if len(header) > 0 {
		calling.SetHeader(header...)
	}
	err = calling.Do()
	if err != nil {
		return nil, fmt.Errorf("do http request error: %s", err.Error())
	}
	return calling.GetResponseBody()
}
