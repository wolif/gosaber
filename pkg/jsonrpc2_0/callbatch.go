package jsonrpc2_0

import (
	"fmt"
)

type callBatch struct {
	client *Client
	calls  map[interface{}]*call // map[Request.ID]*Request
}

func NewCallBatch(client *Client) *callBatch {
	return &callBatch{
		client: client,
		calls:  make(map[interface{}]*call),
	}
}

func (cb *callBatch) Push(c *call) *callBatch {
	cb.calls[c.request.ID] = c
	return cb
}

func (cb *callBatch) Call(method string, params ...interface{}) *call {
	c := NewCall(cb.client, method, params...)
	cb.Push(c)
	return c
}

func (cb *callBatch) Invoke() error {
	data := make([]*Request, 0)
	for _, call := range cb.calls {
		data = append(data, call.request)
	}
	reqData, err := cb.client.protocol.Encode(data)
	if err != nil {
		return fmt.Errorf("encode request error:%s", err.Error())
	}
	respData, err := cb.client.transprot.Send(reqData)
	if err != nil {
		return fmt.Errorf("send request error: %s", err.Error())
	}
	resps := make([]*Response, 0)
	if err = cb.client.protocol.Decode(respData, resps); err != nil {
		return fmt.Errorf("can't parse response data from server, origin data:%s", string(respData))
	}
	respMap := make(map[interface{}]*Response)
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
