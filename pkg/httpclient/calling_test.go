package httpclient

import (
	"testing"
)

func TestClient(t *testing.T) {
	calling, err := Get("http://cn.bing.com")
	if err != nil {
		t.Error(err)
		return
	} else {
		t.Log(calling.GetRespStatusCode())
		t.Log(calling.GetRespBodyString())
	}

	// 	body := `{
	// 	"jsonrpc":"2.0",
	// 	"method": "Catalog\\Area.parent",
	// 	"params": [
	// 		246
	// 	]
	// }`
	// calling, err = Post("http://catalog-service.dev.tgs.com", body)
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// } else {
	// 	t.Log(calling.GetRespStatusCode())
	// 	t.Log(calling.GetRespBodyString())
	// }
}
