package jsonrpc2_0

type Transport interface {
	Send(reqData []byte) (respData []byte, err error)
}
