package jsonrpc2_0

import (
	"context"
	"fmt"
)

type CallBatch struct {
	client *Client
	calls  map[int64]*Call // map[Request.ID]*Request
}

func NewCallBatch(client *Client) *CallBatch {
	return &CallBatch{
		client: client,
		calls:  make(map[int64]*Call),
	}
}

func (cb *CallBatch) Push(c *Call) *CallBatch {
	cb.calls[c.request.ID] = c
	return cb
}

func (cb *CallBatch) Call(method string, params ...interface{}) *Call {
	c := NewCall(cb.client, method, params...)
	cb.Push(c)
	return c
}

func (cb *CallBatch) Invoke(ctx context.Context) error {
	data := make([]*Request, 0)
	for _, call := range cb.calls {
		data = append(data, call.request)
	}
	reqData, err := cb.client.protocol.Encode(data)
	if err != nil {
		return fmt.Errorf("encode request error:%s", err.Error())
	}
	respData, err := cb.client.transport.Send(ctx, reqData)
	if err != nil {
		return fmt.Errorf("send request error: %s", err.Error())
	}
	resps := make([]*Response, 0)
	if err = cb.client.protocol.Decode(respData, &resps); err != nil {
		return fmt.Errorf("can't parse response data from server, origin data:%s", string(respData))
	}
	respMap := make(map[int64]*Response)
	for _, resp := range resps {
		respMap[resp.ID] = resp
	}

	for id, resp := range respMap {
		if _, ok := cb.calls[id]; ok {
			cb.calls[id].response = resp
		}
	}
	for _, call := range cb.calls {
		if call.response == nil {
			return fmt.Errorf("response lacked with id %v", call.request.ID)
		}
	}

	return nil
}