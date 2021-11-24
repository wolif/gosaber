package httpclient

import (
	"net/http"
	"testing"
)

func TestSendRequest(t *testing.T) {
	resp, err := SendRequest("GET", "http://www.baidu.com", "", http.Header{}, Options{})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(string(resp))
	}
}
