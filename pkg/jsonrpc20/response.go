package jsonrpc20

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
	E_PARSE:            "语法解析错误",
	E_INVALID_REQ:      "无效请求",
	E_METHOD_NOT_FOUND: "找不到方法",
	E_INVALID_PARAMS:   "无效的参数",
	E_INTERNAL:         "内部错误",
	E_SERVER:           "服务端错误",
	E_UNKNOWN:          "错误类型未知",
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

type ResponseError struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
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
		if msg, ok := options[0].(string); ok {
			re.Message = msg
		}
	}
	if len(options) >= 2 {
		re.Data = options[1]
	}
	return re
}

type Response struct {
	JsonRpc string         `json:"jsonrpc"`
	Result  interface{}    `json:"result,omitempty"`
	Error   *ResponseError `json:"error,omitempty"`
	ID      int64          `json:"id"`
}
