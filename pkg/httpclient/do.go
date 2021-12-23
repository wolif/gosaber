package httpclient

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"time"
)

func (c *Calling) reDoInit() *Calling {
	c.response = nil
	c.respBody = nil
	c.timeExpend = nil
	return c
}

func (c *Calling) ReDo(ctx ...context.Context) error {
	return c.reDoInit().Do(ctx...)
}

func (c *Calling) Do(ctx ...context.Context) error {
	if len(ctx) > 0 {
		c.Context(ctx[0])
	}

	if c.GetUrl() == nil {
		return errorf("url has not set")
	}

	if len(c.GetCookies()) > 0 {
		jar, _ := cookiejar.New(nil)
		jar.SetCookies(c.GetUrl(), c.GetCookies())
	}

	if c.GetRequest() == nil {
		request, err := http.NewRequestWithContext(c.GetContext(), c.GetMethod(), c.GetUrl().String(), c.GetBody())
		if err != nil {
			return errorf(err)
		}
		c.Request(request)
	}

	c.fillOptions()
	var err error
	timeStart := time.Now()
	c.response, err = c.GetClient().Do(c.GetRequest())
	timeExpend := time.Since(timeStart)
	c.timeExpend = &timeExpend
	if err != nil {
		return errorf(err)
	}
	return nil
}
