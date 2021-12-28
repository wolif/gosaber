package jsonrpc20

import (
	"github.com/wolif/gosaber/pkg/jsonrpc20/protocol"
	"github.com/wolif/gosaber/pkg/jsonrpc20/transport"
	"time"
)

const BatchCallLimitDef = 100

type Config struct {
	Addr    string        // 服务地址
	Timeout time.Duration // 调用超时
}

type Client struct {
	protocol       Protocol  // 编码格式解析
	transport      Transport // 通信协议
	config         *Config   // 配置
	batchCallLimit int       // 批量调用时的数量上限, 默认 100
}

// setter ----------------------------------------------------------------------

func (c *Client) SetProtocol(protocol Protocol) *Client {
	c.protocol = protocol
	return c
}

func (c *Client) SetTransport(transport Transport) *Client {
	c.transport = transport
	return c
}

func (c *Client) SetbatchCallLimit(n int) *Client {
	c.batchCallLimit = n
	return c
}

// constructor -----------------------------------------------------------------

func NewClient(config *Config) *Client {
	return &Client{
		protocol:       new(protocol.Json),
		transport:      &transport.Http{Url: config.Addr, Timeout: config.Timeout},
		config:         config,
		batchCallLimit: BatchCallLimitDef,
	}
}

// method ----------------------------------------------------------------------

// 调用单个接口
// 参数method 调用方法
// params 参数列表(可选)
func (c *Client) Call(method string, params ...interface{}) *Call {
	return NewCall(c, method, params...)
}

// 一次调用多个接口
func (c *Client) CallBatch() *CallBatch {
	return NewCallBatch(c)
}
