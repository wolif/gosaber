package jsonrpc2_0

type Protocol interface {
	Encode(data interface{}) ([]byte, error)
	Decode(src []byte, dst interface{}) error
}
