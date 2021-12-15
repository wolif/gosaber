package httpclient

import (
	"testing"
)

func TestClient(t *testing.T) {
	calling := New("mysql://towngas_vod_rw:abcdef@10.20.1.20:3306/towngas_vod?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai#abc")
	//t.Logf("%v", calling.GetUrl())
	calling.GetUrl().Host = "abc.com:123"
	t.Log(calling.GetUrl())

	calling1, err := Get("http://cn.bing.com")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(calling1.GetRespBodyString())
	}
}
