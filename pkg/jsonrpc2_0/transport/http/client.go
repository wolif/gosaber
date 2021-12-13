package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	ctx          context.Context
	method       string
	url          string
	header       http.Header
	body         io.Reader
	option       *Option
	response     *http.Response
	responseBody []byte
}

func New(url ...string) *Client {
	c := &Client{
		ctx:    context.TODO(),
		method: GET,
		header: http.Header{},
		body:   strings.NewReader(""),
	}
	if len(url) > 0 {
		c.url = url[0]
	}
	return c
}

func (c *Client) SetContext(ctx context.Context) *Client {
	c.ctx = ctx
	return c
}

func (c *Client) SetMethod(method Method) *Client {
	c.method = method
	return c
}

func (c *Client) SetUrl(url string) *Client {
	c.url = url
	return c
}

func (c *Client) SetBody(httpbody interface{}) *Client {
	switch httpbody.(type) {
	case string:
		c.body = strings.NewReader(httpbody.(string))
	case []byte:
		c.body = bytes.NewReader(httpbody.([]byte))
	default:
		panic("httpbody should be []byte or string")
	}
	return c
}

func (c *Client) SetHeader(headers ...interface{}) *Client {
	if len(headers) == 2 {
		if _, ok1 := headers[0].(string); ok1 {
			if _, ok2 := headers[1].(string); ok2 {
				c.header.Set(headers[0].(string), headers[1].(string))
			}
		}
	} else if len(headers) == 1 {
		switch headers[0].(type) {
		case map[string]string:
			for k, v := range headers[0].(map[string]string) {
				c.SetHeader(k, v)
			}
		case http.Header:
			c.header = headers[0].(http.Header)
		}
	}
	return c
}

func (c *Client) SetOption(opt *Option) *Client {
	c.option = opt
	return c
}

// -----------------------------------------------------------------------------
func (c *Client) POST(url string, body ...interface{}) error {
	if len(body) == 1 {
		c.SetBody(body[0])
	}
	return c.SetUrl(url).SetMethod(POST).Do(c.ctx)
}

func (c *Client) GET(url string) error {
	return c.SetUrl(url).SetMethod(GET).Do(c.ctx)
}

// -----------------------------------------------------------------------------
func (c *Client) set(client *http.Client, request *http.Request) *Client {
	request.Header = c.header
	if c.option == nil {
		return c
	}
	if c.option.Timeout != nil {
		client.Timeout = *c.option.Timeout
	}
	return c
}

func (c *Client) Do(ctx ...context.Context) (err error) {
	if len(ctx) > 0 {
		c.SetContext(ctx[0])
	}
	if c.url == "" {
		err = fmt.Errorf("url has not set")
		return
	}
	request, err := http.NewRequestWithContext(c.ctx, c.method, c.url, c.body)
	if err != nil {
		return
	}
	httpClient := new(http.Client)
	c.set(httpClient, request)

	c.response, err = httpClient.Do(request)
	if err != nil {
		return
	}
	return
}

func (c *Client) GetResponse() *http.Response {
	return c.response
}

func (c *Client) GetResponseBody() ([]byte, error) {
	if c.responseBody == nil {
		defer c.response.Body.Close()
		rb := make([]byte, 0)
		rb, err := ioutil.ReadAll(c.response.Body)
		if err != nil {
			return nil, err
		}
		c.responseBody = rb
	}

	return c.responseBody, nil
}

func (c *Client) GetResponseStatusCode() int {
	return c.response.StatusCode
}
