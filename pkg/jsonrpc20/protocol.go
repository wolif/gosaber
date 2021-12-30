package jsonrpc20

// 数据解析工具
type Protocol interface {
	// 编码
	Encode(data interface{}) ([]byte, error)
	// 解码
	Decode(src []byte, dst interface{}) error
}
