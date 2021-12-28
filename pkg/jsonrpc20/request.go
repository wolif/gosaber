package jsonrpc20

import (
	"encoding/json"
	"github.com/wolif/gosaber/pkg/jsonrpc20/utils/idgen"
)

type Request struct {
	JsonRpc string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      json.RawMessage `json:"id"`
}

func NewRequest(id ...interface{}) *Request {
	req := &Request{
		JsonRpc: VERSION,
	}
	if len(id) > 0 {
		req.setID(id[0])
	} else {
		req.setID(idgen.Gen())
	}
	return req
}

// setter ----------------------------------------------------------------------

func (r *Request) setMethod(method string) *Request {
	r.Method = method
	return r
}

func (r *Request) setParams(params ...interface{}) *Request {
	if len(params) == 0 {
		return r
	}
	r.Params, _ = json.Marshal(params)
	return r
}

func (r *Request) setParam(param interface{}) *Request {
	r.Params, _ = json.Marshal(param)
	return r
}

func (r *Request) setID(id interface{}) *Request {
	switch id := id.(type) {
	case int, int8, int16, int32, int64, string:
		r.ID, _ = json.Marshal(id)
	default:
		panic("request id should be type integer or string")
	}
	return r
}

// create Response -------------------------------------------------------------

func (r *Request) ResponseResult(result interface{}) *Response {
	return (&Response{JsonRpc: r.JsonRpc, ID: r.ID}).setResult(result)
}

func (r *Request) ResponseError(code ErrorCode, options ...interface{}) *Response {
	return (&Response{JsonRpc: r.JsonRpc, ID: r.ID}).setError(code, options...)
}
