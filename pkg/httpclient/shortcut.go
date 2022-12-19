package httpclient

import (
	"context"
	"io"
	"net/http"
)

/**
 * otherParams type: []byte: body | string: body | io.Reader: body | context.Context: | http.Header | Option
 */
func DoMethod(method Method, url string, otherParams ...interface{}) (*Calling, error) {
	calling := New(url).WithMethod(method)
	for _, param := range otherParams {
		switch p := param.(type) {
		case []byte, string, io.Reader:
			calling.WithBody(p)
		case context.Context:
			calling.WithContext(p)
		case http.Header:
			calling.WithHeader(p)
		case *Option:
			calling.WithOption(p)
		}
	}

	if err := calling.Do(); err != nil {
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
