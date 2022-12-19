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
		WithMethod(httpclient.POST).
		WithHeader("content-type", "application/json").
		WithBody(reqData).
		WithContext(ctx).
		WithOption(httpclient.NewOption().SetTimeout(h.Timeout))
	if len(header) > 0 {
		calling.WithHeader(header...)
	}
	err = calling.Do()
	if err != nil {
		return nil, err
	}
	return calling.GetRespBody()
}
