package gormhelper

import (
	"reflect"
	"testing"
)

func testFn1(value interface{}) interface{} {
	return value
}

func TestWhereOpt(t *testing.T) {
	fnSlice := []func(interface{}) interface{}{testFn1}
	t.Log(reflect.ValueOf(testFn1) == reflect.ValueOf(fnSlice[0]))
}
