package httpclient

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Calling struct {
	ctx      context.Context
	method   string
	url      *url.URL
	header   http.Header
	cookies  []*http.Cookie
	body     io.Reader
	option   *Option
	client   *http.Client
	request  *http.Request
	response *http.Response
	respBody []byte
}

func New(url ...string) *Calling {
	c := new(Calling).init()
	if len(url) > 0 {
		c.Url(url[0])
	}
	return c
}

func (c *Calling) init() *Calling {
	c.ctx = context.TODO()
	c.method = GET
	c.url = nil
	c.header = http.Header{}
	c.cookies = make([]*http.Cookie, 0)
	c.body = bytes.NewReader([]byte{})
	c.option = nil
	c.client = new(http.Client)
	c.request = nil
	c.response = nil
	c.respBody = nil
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

func (c *Calling) Url(u ...interface{}) *Calling {
	if len(u) == 0 {
		return c
	}
	switch p := u[0].(type) {
	case string:
		var err error
		c.url, err = url.Parse(p)
		if err != nil {
			panic(errorf("set url [%s] error: %s", p, err.Error()))
		}
	case *url.URL:
		c.url = p
	}
	return c
}

func (c *Calling) GetUrl() *url.URL {
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
		switch header := headers[0].(type) {
		case map[string]string:
			for k, v := range header {
				c.Header(k, v)
			}
		case http.Header:
			c.header = header
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

func (c *Calling) Cookie(cookie *http.Cookie) *Calling {
	c.cookies = append(c.cookies, cookie)
	return c
}

func (c *Calling) GetCookies() []*http.Cookie {
	return c.cookies
}

func (c *Calling) Body(body interface{}) *Calling {
	switch b := body.(type) {
	case io.Reader:
		c.body = b
	case string:
		c.body = strings.NewReader(b)
	case []byte:
		c.body = bytes.NewReader(b)
	default:
		panic(errorf("http body should be []byte | string | io.reader"))
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

func (c *Calling) fillOptions() *Calling {
	c.request.Header = c.header
	if c.option == nil {
		return c
	}
	if c.option.Timeout != nil {
		c.client.Timeout = *c.option.Timeout
	}
	return c
}
