package httpclient

import (
	"net/http"
	"time"
)

func NewOption() *Option {
	return new(Option)
}

func (c *Calling) fillOptions(client *http.Client, request *http.Request) *Calling {
	request.Header = c.header
	if c.option == nil {
		return c
	}
	if c.option.Timeout != nil {
		client.Timeout = *c.option.Timeout
	}
	return c
}

// -----------------------------------------------------------------------------

type Option struct {
	Timeout *time.Duration
}

func (o *Option) SetTimeout(timout time.Duration) *Option {
	o.Timeout = &timout
	return o
}
