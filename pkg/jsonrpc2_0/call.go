package jsonrpc2_0

import (
	"context"
	"fmt"
)

type Call struct {
	client   *Client
	request  *Request
	response *Response
}

func NewCall(client *Client, method string, params ...interface{}) *Call {
	return &Call{
		client:  client,
		request: NewRequest().SetMethod(method).SetParams(params...),
	}
}

func (c *Call) SetClient(client *Client) *Call {
	c.client = client
	return c
}

func (c *Call) SetMethod(method string) *Call {
	c.request.SetMethod(method)
	return c
}

func (c *Call) SetParams(params ...interface{}) *Call {
	c.request.SetParams(params)
	return c
}

func (c *Call) Invoke(ctx context.Context, params ...interface{}) error {
	if len(params) > 0 {
		c.SetParams(params...)
	}
	reqData, err := c.client.protocol.Encode(c.request)
	if err != nil {
		return fmt.Errorf("encode request error: %s", err.Error())
	}
	respData, err := c.client.transport.Send(ctx, reqData)
	if err != nil {
		return fmt.Errorf("send request error: %s", err.Error())
	}
	resp := new(Response)
	if err = c.client.protocol.Decode(respData, resp); err != nil {
		return fmt.Errorf("can't parse response data from server, origin data:%s", string(respData))
	}
	c.response = resp
	return nil
}

func (c *Call) Resolve(result interface{}) (errResp *ResponseError, err error) {
	if c.response == nil {
		return nil, fmt.Errorf("the request hasn't been sent")
	}
	if c.response.Error != nil {
		return c.response.Error, nil
	}
	if c.response.Result == nil {
		return
	}
	bs, err := c.client.protocol.Encode(c.response.Result)
	if err != nil {
		return nil, fmt.Errorf("encode response.result error: %s", err.Error())
	}
	err = c.client.protocol.Decode(bs, result)
	if err != nil {
		return nil, fmt.Errorf("decode response.result error: %s", err.Error())
	}
	return
}
