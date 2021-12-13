package httpclient

import (
	"context"
	"io"
	"net/http"
)

type Method = string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
)

/**
 * otherParams type: []byte: body | string: body | io.Reader: body | context.Context: | http.Header | Option
 */
func DoMethod(method Method, url string, otherParams ...interface{}) (*Calling, error) {
	calling := New(url).Method(method)
	for _, param := range otherParams {
		switch p := param.(type) {
		case []byte, string, io.Reader:
			calling.Body(p)
		case context.Context:
			calling.Context(p)
		case http.Header:
			calling.Header(p)
		case *Option:
			calling.Option(p)
		}
	}
	err := calling.Do()
	if err != nil {
		return nil, err
	}
	return calling, nil
}

func Post(url string, otherParams ...interface{}) (*Calling, error) {
	return DoMethod(POST, url, otherParams...)
}

func Get(url string, otherParams ...interface{}) (*Calling, error) {
	return DoMethod(GET, url, otherParams...)
}

func Put(url string, otherParams ...interface{}) (*Calling, error) {
	return DoMethod(PUT, url, otherParams...)
}

func Delete(url string, otherParams ...interface{}) (*Calling, error) {
	return DoMethod(DELETE, url, otherParams...)
}
