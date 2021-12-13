package http

import (
	"testing"
)

func TestClient(t *testing.T) {
	calling := New()
	err := calling.GET("http://cn.bing.com")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(calling.GetResponseStatusCode())
		t.Log(calling.GetResponseBody())
	}
}
