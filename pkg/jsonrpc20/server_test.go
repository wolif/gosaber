package jsonrpc20

import (
	"context"
	"encoding/json"
	"github.com/wolif/gosaber/pkg/util/strs"
	"strings"
	"testing"
)

type TestModule struct{}

type Args1 struct {
	I int    `json:"i"`
	S string `json:"s"`
}

type Reply1 struct {
	Data1 string `json:"data1"`
	Data2 string `json:"data2"`
}

func (ts *TestModule) F1(ctx *context.Context, input *Args1, output *Reply1) *ResponseError {
	output.Data1 = "hello"
	output.Data2 = "world"
	return nil
}

func TestServer(t *testing.T) {
	server := NewServer(func(s string) (ModuleName, MethodName) {
		segs := strings.SplitN(s, ".", 2)
		return segs[0], strs.UcFirst(segs[1])
	})
	module := new(TestModule)
	server.RegisterModule(module)

	reqStr := `[
		{"jsonrpc":"2.0", "method":"TestModule.F1", "params":{"i":1,"s":"1"}, "id":1},
		{"jsonrpc":"2.0", "method":"TestModule.F1", "params":{"i":1,"s":"1"}, "id":2},
		{"jsonrpc":"2.0", "method":"test.F1", "params":{"i":1,"s":"1"}, "id":3}
	]`
	ctx := context.TODO()
	resp := server.Dispatch(&ctx, []byte(reqStr))
	switch resp := resp.(type) {
	case *Response:
		if resp.IsSuccess() {
			d := new(Reply1)
			_ = json.Unmarshal(resp.Result, d)
			t.Log(resp.ID, d)
		} else {
			t.Log(resp.ID, resp.Error)
		}
	case []*Response:
		for _, r := range resp {
			if r.IsSuccess() {
				d := new(Reply1)
				_ = json.Unmarshal(r.Result, d)
				t.Log(r.ID, d)
			} else {
				t.Log(r.ID, r.Error)
			}
		}
	}
}
