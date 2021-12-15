package httpclient

import (
	"time"
)

type Option struct {
	Timeout *time.Duration
}

func NewOption() *Option {
	return new(Option)
}

func (o *Option) SetTimeout(timout time.Duration) *Option {
	o.Timeout = &timout
	return o
}
