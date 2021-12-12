package jsonrpc2_0

import "time"

type Config struct {
	Addr    string
	Timeout time.Duration
}

type Client struct {
	protocol  Protocol
	transprot Transport
	config    *Config
}

func NewClient(config *Config) *Client {

	return nil
}
