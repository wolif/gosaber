package transport

import (
	"context"
	"time"

	"github.com/wolif/gosaber/pkg/httpclient"
)

type Http struct {
	Url     string
	Timeout time.Duration
}

func (h *Http) Send(ctx context.Context, reqData []byte, header ...interface{}) (respData []byte, err error) {
	calling := httpclient.New(h.Url).
		Method(httpclient.POST).
		Header("content-type", "application/json").
		Body(reqData).
		Context(ctx).
		Option(httpclient.NewOption().SetTimeout(h.Timeout))
	if len(header) > 0 {
		calling.Header(header...)
	}
	err = calling.Do()
	if err != nil {
		return nil, err
	}
	return calling.GetRespBody()
}
