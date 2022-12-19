package jsonrpc20

import (
	"encoding/json"
)

type ErrorCode = int

const (
	E_PARSE            ErrorCode = -32700
	E_INVALID_REQ      ErrorCode = -32600
	E_METHOD_NOT_FOUND ErrorCode = -32601
	E_INVALID_PARAMS   ErrorCode = -32602
	E_INTERNAL         ErrorCode = -32603
	E_SERVER           ErrorCode = -32000
	E_UNKNOWN          ErrorCode = -32001
)

var errMsgMap = map[ErrorCode]string{
	E_PARSE:            "syntax parse error",
	E_INVALID_REQ:      "invalid request",
	E_METHOD_NOT_FOUND: "method not found",
	E_INVALID_PARAMS:   "invalid arguments",
	E_INTERNAL:         "internal error",
	E_SERVER:           "server error",
	E_UNKNOWN:          "unknown error",
}

var DEF_ERROR_MESSAGE = E_UNKNOWN

func ErrorMessage(code ErrorCode, message ...string) string {
	if len(message) > 0 {
		return message[0]
	}
	if msg, ok := errMsgMap[code]; ok {
		return msg
	}
	return errMsgMap[DEF_ERROR_MESSAGE]
}

// ResponseError ---------------------------------------------------------------

type ResponseError struct {
	Code    ErrorCode       `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (re *ResponseError) Error() string {
	return re.Message
}

func NewResponseError(code ErrorCode, options ...interface{}) *ResponseError {
	re := &ResponseError{
		Code:    code,
		Message: ErrorMessage(code),
	}
	if len(options) >= 1 {
		switch e := options[0].(type) {
		case string:
			re.Message = e
		case error:
			re.Message = e.Error()
		}
	}
	if len(options) >= 2 {
		re.Data, _ = json.Marshal(options[1])
	}
	return re
}

func RespErr(code ErrorCode, options ...interface{}) *ResponseError {
	return NewResponseError(code, options...)
}

func RespErrParse(options ...interface{}) *ResponseError {
	return NewResponseError(E_PARSE, options...)
}

func RespErrInvalidReq(options ...interface{}) *ResponseError {
	return NewResponseError(E_INVALID_REQ, options...)
}

func RespErrMethodNotFound(options ...interface{}) *ResponseError {
	return NewResponseError(E_METHOD_NOT_FOUND, options...)
}

func RespErrInvalidParams(options ...interface{}) *ResponseError {
	return NewResponseError(E_INVALID_PARAMS, options...)
}

func RespErrInternal(options ...interface{}) *ResponseError {
	return NewResponseError(E_INTERNAL, options...)
}

func RespErrServer(options ...interface{}) *ResponseError {
	return NewResponseError(E_SERVER, options...)
}

func RespErrUnknown(options ...interface{}) *ResponseError {
	return NewResponseError(E_UNKNOWN, options...)
}

// response --------------------------------------------------------------------

type Response struct {
	JsonRpc string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *ResponseError  `json:"error,omitempty"`
	ID      json.RawMessage `json:"id"`
}

func (r *Response) setResult(result interface{}) *Response {
	r.Result, _ = json.Marshal(result)
	return r
}

func (r *Response) setError(code ErrorCode, options ...interface{}) *Response {
	r.Error = NewResponseError(code, options...)
	return r
}

func NewErrorResponse(code ErrorCode, options ...interface{}) *Response {
	return &Response{
		JsonRpc: VERSION,
		Error:   NewResponseError(code, options...),
	}
}

func (r *Response) IsSuccess() bool {
	return r.Error == nil
}

func (r *Response) IsErrorCode(code ErrorCode) bool {
	return !r.IsSuccess() && r.Error.Code == code
}
