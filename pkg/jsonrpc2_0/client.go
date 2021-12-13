package jsonrpc2_0

import (
	"github.com/wolif/gosaber/pkg/jsonrpc2_0/protocol"
	"github.com/wolif/gosaber/pkg/jsonrpc2_0/transport"
	"time"
)

type Config struct {
	Addr    string
	Timeout time.Duration
}

type Client struct {
	protocol  Protocol
	transport Transport
	config    *Config
}

func (c *Client) SetProtocol(protocol Protocol) *Client {
	c.protocol = protocol
	return c
}

func (c *Client) SetTransport(transport Transport) *Client {
	c.transport = transport
	return c
}

func NewClient(config *Config) *Client {
	return &Client{
		protocol: new(protocol.Json),
		transport: &transport.HttpClient{
			Url:     config.Addr,
			Timeout: config.Timeout,
		},
		config: config,
	}
}

func (c *Client) Call(method string, params ...interface{}) *Call {
	return NewCall(c, method, params...)
}

func (c *Client) CallBatch() *CallBatch {
	return NewCallBatch(c)
}