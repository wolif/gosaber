package httpclient

import "fmt"

type Error string

func (e Error) Error() string {
	return "http client error: " + string(e)
}

func errorf(format interface{}, a...interface{}) Error {
	switch e := format.(type) {
	case error:
		return Error(e.Error())
	case string:
		return Error(fmt.Sprintf(e, a...))
	}
	panic("http client: code error ...")
}
