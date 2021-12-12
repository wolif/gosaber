package jsonrpc2_0

import "fmt"

type call struct {
	client   *Client
	request  *Request
	response *Response
}

func NewCall(client *Client, method string, params ...interface{}) *call {
	return &call{
		client:  client,
		request: NewRequest().SetMethod(method).SetParams(params...),
	}
}

func (c *call) SetClient(client *Client) *call {
	c.client = client
	return c
}

func (c *call) SetMethod(method string) *call {
	c.request.SetMethod(method)
	return c
}

func (c *call) Invoke(params ...interface{}) error {
	c.request.SetParams(params...)
	reqData, err := c.client.protocol.Encode(c.request)
	if err != nil {
		return fmt.Errorf("encode request error: %s", err.Error())
	}
	respData, err := c.client.transprot.Send(reqData)
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

func (c *call) Resolve(result interface{}) (errResp *ResponseError, err error) {
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
