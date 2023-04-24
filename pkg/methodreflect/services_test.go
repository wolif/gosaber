package methodreflect

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/wolif/gosaber/pkg/cooc/example"
)

func TestServices_Register(t *testing.T) {
	s := &Services{
		ServiceMethodDelimiter: "_",
		Method1stArgType:       reflect.TypeOf((*context.Context)(nil)).Elem(),
		MethodReturnType: []reflect.Type{
			reflect.TypeOf((*string)(nil)).Elem(),
			reflect.TypeOf((*error)(nil)).Elem(),
		},
	}

	s.Register("Serv", example.Serv)

	for serviceName, service := range s.serviceMap {
		for methodName := range service.Methods {
			fmt.Printf("serviceName: %s, methodName: %s", serviceName, methodName)
		}
	}

	t.Log(s.Get("Serv_F1"))
	t.Log(s.GetWithName("Serv", "F1"))
}
