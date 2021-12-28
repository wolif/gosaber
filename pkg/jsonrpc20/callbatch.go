package jsonrpc20

import (
	"context"
	"fmt"
)

type CallBatch struct {
	client *Client
	calls  map[interface{}]*Call
}

// constructor -----------------------------------------------------------------

func NewCallBatch(client *Client) *CallBatch {
	return &CallBatch{
		client: client,
		calls:  make(map[interface{}]*Call),
	}
}

// method ----------------------------------------------------------------------

// 添加单次调用
// 当添加量超过 批量调用时的数量上限 时,会静默处理(不添加,也不报错)
func (cb *CallBatch) Push(calls ...*Call) *CallBatch {
	for _, c := range calls {
		if len(cb.calls) >= cb.client.batchCallLimit {
			break
		}
		cb.calls[string(c.request.ID)] = c
	}
	return cb
}

// 添加单词调用
// 参数 method 是调用的方法
// 参数 params 是参数列表
func (cb *CallBatch) Call(method string, params ...interface{}) *Call {
	c := NewCall(cb.client, method, params...)
	cb.Push(c)
	return c
}

// 发起调用
// 当调用过程(非结果)发生错误时, 会返回 error
func (cb *CallBatch) Invoke(ctx context.Context) error {
	data := make([]*Request, 0)
	for _, call := range cb.calls {
		data = append(data, call.request)
	}
	reqData, err := cb.client.protocol.Encode(data)
	if err != nil {
		return fmt.Errorf("jsonrpc error: encode request error:%s", err.Error())
	}
	respData, err := cb.client.transport.Send(ctx, reqData)
	if err != nil {
		return fmt.Errorf("jsonrpc error: send request error: %s", err.Error())
	}

	// 可能是一些未知类型的错误, 最常见是解析不了requestbody
	resp := new(Response)
	if err = cb.client.protocol.Decode(respData, &resp); err == nil {
		if !resp.IsSuccess() {
			return resp.Error
		}
		return fmt.Errorf(
			"jsonrpc errror: unkonwn error type, origin data from server: %s",
			string(respData),
		)
	}

	resps := make([]*Response, 0)
	if err = cb.client.protocol.Decode(respData, &resps); err != nil {
		return fmt.Errorf(
			"jsonrpc error: can't parse response data from server, origin data:%s",
			string(respData),
		)
	}

	for _, resp := range resps {
		if _, ok := cb.calls[string(resp.ID)]; ok {
			cb.calls[string(resp.ID)].response = resp
		}
	}
	for _, call := range cb.calls {
		if call.response == nil { // 有调用没有返会调用结果
			return fmt.Errorf(
				"jsonrpc error: response lacked with id %v",
				call.request.ID,
			)
		}
	}

	return nil
}
