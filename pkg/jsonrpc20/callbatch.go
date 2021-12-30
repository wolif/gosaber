package jsonrpc20

import (
	"context"
)

// 客户端批量调用实体
type CallBatch struct {
	client *Client
	calls  map[interface{}]*Call
}

// constructor -----------------------------------------------------------------

// 创建批量调用实体
func NewCallBatch(client *Client) *CallBatch {
	return &CallBatch{
		client: client,
		calls:  make(map[interface{}]*Call),
	}
}

// method ----------------------------------------------------------------------

// 添加调用实体;
// 当添加量超过批量调用时的数量上限时,会报错
func (cb *CallBatch) Push(calls ...*Call) error {
	for _, c := range calls {
		if len(cb.calls) >= cb.client.config.FetchBatchCallLimit() {
			return BatchCallOverSize
		}
		cb.calls[string(c.request.ID)] = c
	}
	return nil
}

// 添加单词调用;
// 参数:
// method 是调用的方法;
// params 是参数列表;
func (cb *CallBatch) Call(method string, params ...interface{}) (*Call, error) {
	c := NewCall(cb.client, method, params...)
	if err := cb.Push(c); err != nil {
		return nil, err
	}
	return c, nil
}

// 发起调用
// 当调用过程(非结果)发生错误时, 会返回 error
// 调用的返回结果在添加的各个 *jsonrpc.Call 中
func (cb *CallBatch) Invoke(ctx context.Context) error {
	data := make([]*Request, 0)
	for _, call := range cb.calls {
		data = append(data, call.request)
	}
	reqData, err := cb.client.protocol.Encode(data)
	if err != nil {
		return ClientErrorf("encode request error:%s", err.Error())
	}
	respData, err := cb.client.transport.Send(ctx, reqData)
	if err != nil {
		return ClientErrorf("send request error: %s", err.Error())
	}

	// 可能是一些未知类型的错误, 最常见是解析不了requestbody
	resp := new(Response)
	if err = cb.client.protocol.Decode(respData, &resp); err == nil {
		if !resp.IsSuccess() {
			return resp.Error
		}
		return ClientErrorf(
			"unkonwn error type, origin data from server: %s",
			string(respData),
		)
	}

	resps := make([]*Response, 0)
	if err = cb.client.protocol.Decode(respData, &resps); err != nil {
		return ClientErrorf(
			"jcan't parse response data from server, origin data: %s",
			string(respData),
		)
	}

	for _, resp := range resps {
		if _, ok := cb.calls[string(resp.ID)]; ok {
			cb.calls[string(resp.ID)].response = resp
		}
	}

	return nil
}
