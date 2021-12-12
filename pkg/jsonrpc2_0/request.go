package jsonrpc2_0

import "github.com/wolif/gosaber/pkg/jsonrpc2_0/utils/idgen"

type Request struct {
	JsonRpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	ID      interface{} `json:"id"`
}

func NewRequest(id ...interface{}) *Request {
	req := &Request{
		JsonRpc: VERSION,
	}
	if len(id) > 0 {
		req.ID = id[0]
	} else {
		req.ID = idgen.GenID()
	}
	return req
}

func (r *Request) SetMethod(method string) *Request {
	r.Method = method
	return r
}

func (r *Request) SetParams(params ...interface{}) *Request {
	l := len(params)
	if l == 0 {
		return r
	}
	if l == 1 {
		r.Params = params[0]
	} else {
		r.Params = params
	}

	return r
}

func (r *Request) ResponseResult(result interface{}) *Response {
	return &Response{
		JsonRpc: r.JsonRpc,
		Result:  result,
		ID:      r.ID,
	}
}

func (r *Request) ResponseError(code ErrorCode, options ...interface{}) *Response {
	return &Response{
		JsonRpc: r.JsonRpc,
		Error:   NewResponseError(code, options...),
		ID:      r.ID,
	}
}
