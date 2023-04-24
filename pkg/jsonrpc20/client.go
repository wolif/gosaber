package jsonrpc20

import (
	"github.com/wolif/gosaber/pkg/jsonrpc20/protocol"
	"github.com/wolif/gosaber/pkg/jsonrpc20/transport"
)

// jsonrpc客户端
type Client struct {
	config    *ClientConfig // 客户端配置
	protocol  Protocol      // 编码格式解析器
	transport Transport     // 通信协议
}

// constructor -----------------------------------------------------------------
// 创建新的jsonrpc客户端(包含默认配置)
func NewClient(config *ClientConfig) *Client {
	return &Client{
		config:    config,
		protocol:  new(protocol.Json),
		transport: &transport.Http{Url: config.Addr, Timeout: config.Timeout},
	}
}

// setter ----------------------------------------------------------------------

// 设置解析协议
func (c *Client) WithProtocol(protocol Protocol) *Client {
	c.protocol = protocol
	return c
}

// 设置通讯协议
func (c *Client) WithTransport(transport Transport) *Client {
	c.transport = transport
	return c
}

// method ----------------------------------------------------------------------

// 调用单个接口;
// 参数:
// method 调用方法;
// params 参数列表(可选);
func (c *Client) Call(method string, params ...interface{}) *Call {
	return NewCall(c, method, params...)
}

// 一次调用多个接口
func (c *Client) CallBatch() *CallBatch {
	return NewCallBatch(c)
}
