package jsonrpc20

type Protocol interface {
	Encode(data interface{}) ([]byte, error)
	Decode(src []byte, dst interface{}) error
}
