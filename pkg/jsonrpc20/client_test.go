package jsonrpc20

import (
	"context"
	"testing"
	"time"
)

type tResult struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"fullname"`
	Level    int    `json:"level"`
	ADCode   int    `json:"ad_code"`
	RealName string `json:"real_name"`
}

func TestClient_CallBatch(t *testing.T) {
	client := NewClient(&Config{
		Addr:    "http://catalog-service.dev.tgs.com",
		Timeout: time.Second * 2,
	})
	client.SetbatchCallLimit(100)

	callBatch := client.CallBatch()
	call1 := client.Call("Catalog\\Area.parent", 378)
	call2 := client.Call("Catalog\\Area.parent", 246)
	call3 := client.Call("Catalog\\Area.parent", 378)
	call4 := client.Call("Catalog\\Area.parent", 246)
	callBatch.Push(call1).Push(call2).Push(call3).Push(call4)
	call5 := callBatch.Call("Catalog\\Area.parent", 378)
	call6 := callBatch.Call("Catalog\\Area.parent", 246)
	call7 := callBatch.Call("Catalog\\Area.parent", 378)
	call8 := callBatch.Call("Catalog\\Area.parent", 246)

	err := callBatch.Invoke(context.TODO())
	if err != nil {
		t.Error(err)
		return
	}

	for _, call := range []*Call{call1, call2, call3, call4, call5, call6, call7, call8} {
		result := new(tResult)
		respErr, err := call.Resolve(result)
		if respErr != nil {
			t.Error(respErr)
		} else if err != nil {
			t.Error(err)
		} else {
			t.Log(result)
		}
	}

	call := client.Call("Catalog\\Area.parent").SetParam(246)
	err = call.Invoke(context.TODO())
	if err != nil {
		t.Error(err)
		return
	}
	result := new(tResult)
	respErr, err := call.Resolve(result)
	if respErr != nil {
		t.Error(respErr)
	} else if err != nil {
		t.Error(err)
	} else {
		t.Log(result)
	}
}
