package jsonrpc20

import "context"

// 通信协议
type Transport interface {
	// 发送数据并返回结果
	Send(ctx context.Context, reqData []byte, header ...interface{}) (respData []byte, err error)
}
