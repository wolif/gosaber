package cooc

import (
	"context"
	"reflect"
	"testing"

	"github.com/wolif/gosaber/pkg/cooc/example"
)

func TestCooc(t *testing.T) {
	col := New(reflect.TypeOf((*context.Context)(nil)).Elem()).
		SetMethodRetType(reflect.TypeOf((*string)(nil)).Elem(), reflect.TypeOf((*error)(nil)).Elem())
	col.Add("c", example.Serv)
	// methods := col.Add("c", example.Serv)
	// t.Log(methods)

	col.MethodsOn(map[string][]string{"c": {"F4"}}, EV_BEFORE_METHOD_INVOKING, func(i ...interface{}) {
		m := i[0].(*method)
		ctx := i[1].(*context.Context)
		args := i[2].(*string)
		t.Log(m.Name)
		t.Log(ctx)
		t.Log(*args)
		*args = "some thing other"
	})
	ctx := context.Background()
	_, _, err := col.Method("c", "F4").Invoke(&ctx, []byte(`"some stirng"`))
	if err != nil {
		t.Error(err)
	}
	_, _, err = col.Method("c", "F4").InvokeCustomized(ctx, `"asdf"`, func(fn reflect.Method, argType, replyType reflect.Type, emit func(event, ...interface{}), ctx, args interface{}) (interface{}, []interface{}, error) {
		return nil, nil, nil
	})
	if err != nil {
		t.Error(err)
	}
}
