package httpclient

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Calling) reDoInit() *Calling {
	c.response = nil
	c.respBody = nil
	return c
}

func (c *Calling) ReDo(ctx ...context.Context) error {
	return c.reDoInit().Do(ctx...)
}

func (c *Calling) Do(ctx ...context.Context) error {
	if len(ctx) > 0 {
		c.Context(ctx[0])
	}

	if c.GetUrl() == "" {
		return fmt.Errorf("url has not set")
	}

	if c.GetRequest() == nil {
		request, err := http.NewRequestWithContext(c.GetContext(), c.GetMethod(), c.GetUrl(), c.GetBody())
		if err != nil {
			return err
		}
		c.Request(request)
	}

	c.fillOptions(c.GetClient(), c.GetRequest())
	var err error
	c.response, err = c.GetClient().Do(c.GetRequest())
	if err != nil {
		return err
	}
	return nil
}
