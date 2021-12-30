package jsonrpc20

import "time"

// 默认批量调用上线
const BatchCallLimitDef = 100

// 默认接口调用超时时间
const ClientTimeoutDef = 2 * time.Second

// jsonrpc客户端配置
type ClientConfig struct {
	Addr           string        // 服务地址
	Timeout        time.Duration // 调用超时
	BatchCallLimit int           // 批量调用时的数量上限, 默认 100
}

// 创建新的jsonrpc客户端配置(带默认选项)
func NewClientConfig(addr string) *ClientConfig {
	return &ClientConfig{
		Addr:           addr,
		Timeout:        ClientTimeoutDef,
		BatchCallLimit: BatchCallLimitDef,
	}
}

// 设置接口调用超时时间
func (c *ClientConfig) WithTimeout(t time.Duration) *ClientConfig {
	c.Timeout = t
	return c
}

// 设置服务调用地址
func (c *ClientConfig) WithAddress(addr string) *ClientConfig {
	c.Addr = addr
	return c
}

// 设置批量调用上限
func (c *ClientConfig) WithbatchCallLimit(n int) *ClientConfig {
	c.BatchCallLimit = n
	return c
}

// 获取批量调用上限
func (c *ClientConfig) FetchBatchCallLimit() int {
	if c.BatchCallLimit == 0 {
		return BatchCallLimitDef
	}
	return c.BatchCallLimit
}
