package jsonrpc20

import (
	"context"
	"fmt"
)

type Call struct {
	client   *Client
	request  *Request
	response *Response
}

// constructor -----------------------------------------------------------------

func NewCall(client *Client, method string, params ...interface{}) *Call {
	return &Call{
		client:  client,
		request: NewRequest().setMethod(method).setParams(params...),
	}
}

// setter ----------------------------------------------------------------------

func (c *Call) SetClient(client *Client) *Call {
	c.client = client
	return c
}

// 设置调用的方法
func (c *Call) SetMethod(method string) *Call {
	c.request.setMethod(method)
	return c
}

// 设置调用参数列表
// 会被编码为 request.params = [params[0], params[1], params[2]...]
func (c *Call) SetParams(params ...interface{}) *Call {
	c.request.setParams(params)
	return c
}

// 设置调用参数
// 参数 param 本身会成为 request.params 的值(request.params = param)
func (c *Call) SetParam(param interface{}) *Call {
	c.request.setParam(param)
	return c
}

// method ----------------------------------------------------------------------

// 执行调用
// 当调用过程(非结果)出现错误时, 返回 error
func (c *Call) Invoke(ctx context.Context, params ...interface{}) error {
	if len(params) > 0 {
		c.SetParams(params...)
	}
	reqData, err := c.client.protocol.Encode(c.request)
	if err != nil {
		return fmt.Errorf("jsonrpc error: encode request error: %s", err.Error())
	}
	respData, err := c.client.transport.Send(ctx, reqData)
	if err != nil {
		return fmt.Errorf("jsonrpc error: send request error: %s", err.Error())
	}
	resp := new(Response)
	if err = c.client.protocol.Decode(respData, resp); err != nil {
		return fmt.Errorf("jsonrpc error: can't parse response data from server, origin data:%s", string(respData))
	}
	c.response = resp
	return nil
}

// 取值调用结果
// 参数 result 需要一个指针, 会被填充 response.result
// 返回值
// 当 response 中有 错误时, 返回值 errResp 会填充错误
// 其他错误会在返回值 err 中
func (c *Call) Resolve(result interface{}) (errResp *ResponseError, err error) {
	if c.response == nil {
		return nil, fmt.Errorf("jsonrpc error: the request hasn't been sent")
	}
	if c.response.Error != nil {
		return c.response.Error, nil
	}
	if c.response.Result == nil {
		return
	}
	err = c.client.protocol.Decode(c.response.Result, result)
	if err != nil {
		return nil, fmt.Errorf("jsonrpc error: decode response.result error: %s", err.Error())
	}
	return
}
