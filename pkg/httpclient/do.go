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
		c.WithContext(ctx[0])
	}

	if c.Url() == nil {
		return errorf("url has not set")
	}

	if len(c.Cookies()) > 0 {
		jar, _ := cookiejar.New(nil)
		jar.SetCookies(c.Url(), c.Cookies())
	}

	if c.Request() == nil {
		request, err := http.NewRequestWithContext(c.Context(), c.Method(), c.Url().String(), c.Body())
		if err != nil {
			return errorf(err)
		}
		c.WithRequest(request)
	}

	c.fillOptions()
	var err error
	timeStart := time.Now()
	c.response, err = c.Client().Do(c.Request())
	timeExpend := time.Since(timeStart)
	c.timeExpend = &timeExpend
	if err != nil {
		return errorf(err)
	}
	return nil
}
