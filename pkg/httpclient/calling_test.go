package httpclient

import (
	"testing"
)

func TestClient(t *testing.T) {
	calling := New("mysql://dbUser:dbPwd@127.0.0.1:3306/dbName?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai#abc")
	//t.Logf("%v", calling.GetUrl())
	calling.Url().Host = "abc.com:123"
	calling.Url().Path = "path/to.x"
	t.Log(calling.Url())

	// calling1, err := Get("http://cn.bing.com")
	// if err != nil {
	// 	t.Error(err)
	// } else {
	// 	t.Log(calling1.GetRespBodyString())
	// }
}
