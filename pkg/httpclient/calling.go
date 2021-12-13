package httpclient

import "C"
import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
)

type Calling struct {
	ctx      context.Context
	method   string
	url      string
	header   http.Header
	body     io.Reader
	option   *Option
	client   *http.Client
	request  *http.Request
	response *http.Response
	respBody []byte
}

func (c *Calling) init() *Calling {
	c.ctx = context.TODO()
	c.method = GET
	c.url = ""
	c.header = http.Header{}
	c.body = bytes.NewReader([]byte{})
	c.option = nil
	c.client = new(http.Client)
	c.request = nil
	c.response = nil
	c.respBody = nil
	return c
}

func New(url ...string) *Calling {
	c := new(Calling).init()
	if len(url) > 0 {
		c.Url(url[0])
	}
	return c
}

func (c *Calling) Context(ctx context.Context) *Calling {
	c.ctx = ctx
	return c
}

func (c *Calling) GetContext() context.Context {
	return c.ctx
}

func (c *Calling) Method(method Method) *Calling {
	c.method = method
	return c
}

func (c *Calling) GetMethod() Method {
	return c.method
}

func (c *Calling) Url(url string) *Calling {
	c.url = url
	return c
}

func (c *Calling) GetUrl() string {
	return c.url
}

func (c *Calling) Header(headers ...interface{}) *Calling {
	if len(headers) == 2 {
		if k, ok1 := headers[0].(string); ok1 {
			if v, ok2 := headers[1].(string); ok2 {
				c.header.Set(k, v)
			}
		}
	} else if len(headers) == 1 {
		switch headers[0].(type) {
		case map[string]string:
			for k, v := range headers[0].(map[string]string) {
				c.Header(k, v)
			}
		case http.Header:
			c.header = headers[0].(http.Header)
		}
	}
	return c
}

func (c *Calling) GetHeader() http.Header {
	return c.header
}

func (c *Calling) GetHeaderByKey(key string) string {
	return c.header.Get(key)
}

func (c *Calling) Body(body interface{}) *Calling {
	switch body.(type) {
	case io.Reader:
		c.body = body.(io.Reader)
	case string:
		c.body = strings.NewReader(body.(string))
	case []byte:
		c.body = bytes.NewReader(body.([]byte))
	default:
		panic("httpbody should be []byte or string")
	}
	return c
}

func (c *Calling) GetBody() io.Reader {
	return c.body
}

func (c *Calling) Option(opt *Option) *Calling {
	c.option = opt
	return c
}

func (c *Calling) GetOption() *Option {
	return c.option
}

func (c *Calling) Client(client *http.Client) *Calling {
	c.client = client
	return c
}

func (c *Calling) GetClient() *http.Client {
	return c.client
}

func (c *Calling) Request(request *http.Request) *Calling {
	c.request = request
	return c
}

func (c *Calling) GetRequest() *http.Request {
	return c.request
}