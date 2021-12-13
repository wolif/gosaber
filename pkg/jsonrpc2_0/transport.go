package jsonrpc2_0

import "context"

type Transport interface {
	Send(ctx context.Context, reqData []byte, header ...interface{}) (respData []byte, err error)
}
