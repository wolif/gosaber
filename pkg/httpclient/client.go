package httpclient

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type HttpClient interface {
	SendRequest(ctx context.Context, method string, url string, body io.Reader, headers map[string]string, timeout time.Duration) ([]byte, error)
	Post(ctx context.Context, url string, body io.Reader, headers map[string]string, timeout time.Duration) ([]byte, error)
	Get(ctx context.Context, url string, headers map[string]string, timeout time.Duration) ([]byte, error)
}

type Options struct {
	Timeout time.Duration
}

func SendRequest(method string, url string, body string, header http.Header, options Options) ([]byte, error) {
	request, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header = header

	client := http.Client{}

	if options.Timeout != 0 {
		client.Timeout = options.Timeout
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return responseBody, fmt.Errorf("%s %s return %d", method, url, response.StatusCode)
	}

	return responseBody, nil
}

func Post(url string, body string, header http.Header, options Options) ([]byte, error) {
	return SendRequest("POST", url, body, header, options)
}

func Get(url string, body string, header http.Header, options Options) ([]byte, error) {
	return SendRequest("GET", url, body, header, options)
}
