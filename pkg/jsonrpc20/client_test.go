package jsonrpc20

import (
	"context"
	"testing"
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
	client := NewClient(NewClientConfig("http://catalog-service.dev.tgs.com"))

	callBatch := client.CallBatch()
	call1 := client.Call("Catalog\\Area.parent", 378)
	call2 := client.Call("Catalog\\Area.parent", 246)
	call3 := client.Call("Catalog\\Area.parent", 378)
	call4 := client.Call("Catalog\\Area.parent", 246)
	callBatch.Push(call1, call2, call3, call4)
	call5, _ := callBatch.Call("Catalog\\Area.parent", 378)
	call6, _ := callBatch.Call("Catalog\\Area.parent", 246)
	call7, _ := callBatch.Call("Catalog\\Area.parent", 378)
	call8, _ := callBatch.Call("Catalog\\Area.parent", 246)

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
}

func TestClient_Call(t *testing.T) {
	client := NewClient(NewClientConfig("http://catalog-service.dev.tgs.com"))

	call := client.Call("Catalog\\Area.parent").SetParam([]int{246})
	err := call.Invoke(context.TODO())
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
