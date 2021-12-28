package jsonrpc20

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/wolif/gosaber/pkg/jsonrpc20/utils"
	"github.com/wolif/gosaber/pkg/ref"
)

type ServerModule struct {
	name        ModuleName
	server      *Server
	entity      *ref.Entity
	methodsName map[MethodName]struct{}
}

// setter ----------------------------------------------------------------------

// 初始化模块,并注册模块方法
// 如果方法不符合注册条件会静默处理(不注册也不报错)
func (sm *ServerModule) init() {
	sm.methodsName = make(map[string]struct{})
	methods, _ := sm.entity.StructMethods()
	for _, method := range methods {
		if sm.registerMethod(method) == nil {
			sm.methodsName[method.Name] = struct{}{}
		}
	}
}

// 注册模块方法
// 如果方法不符合注册条件会返回错误
func (sm *ServerModule) registerMethod(method *reflect.Method) error {
	if method.Type.NumIn() != 4 {
		return fmt.Errorf(
			"jsonrpc error: the method [%s] need 3 args",
			method.Name,
		)
	}
	if method.Type.In(1) != sm.server.ctxType {
		return fmt.Errorf(
			"jsonrpc error: the method [%s]args 0 should be type [%s]",
			method.Name,
			method.Type.In(1).String(),
		)
	}
	if method.Type.In(2).Kind() != reflect.Ptr || method.Type.In(3).Kind() != reflect.Ptr {
		return fmt.Errorf(
			"jsonrpc error: the method [%s] args 1 & 2 should be a ptr",
			method.Name,
		)
	}
	if method.Type.NumOut() != 1 {
		return fmt.Errorf(
			"jsonrpc error: the method [%s] need 1 return values",
			method.Name,
		)
	}
	if method.Type.Out(0) != reflect.TypeOf((*ResponseError)(nil)) {
		return fmt.Errorf(
			"jsonrpc error: the method [%s] return value 0 should be type [*jsonrpc20.Response]",
			method.Name,
		)
	}

	for i := 1; i <= 3; i++ {
		if !utils.IsSymbolExportedOrBuiltin(method.Type.In(i)) {
			return fmt.Errorf(
				"jsonrpc error: the method [%s] args 1 & 2 should be avild",
				method.Name,
			)
		}
	}

	return nil
}

// getter ----------------------------------------------------------------------

func (sm *ServerModule) Name() ModuleName {
	return sm.name
}

func (sm *ServerModule) Server() *Server {
	return sm.server
}

func (sm *ServerModule) Method(metName MethodName) (*reflect.Method, error) {
	_, ok := sm.methodsName[metName]
	if ok {
		method, _ := sm.entity.StructMethodGet(metName)
		return method, nil
	}
	return nil, fmt.Errorf("josnrpc error: method named [%s] in module [%s]", metName, sm.name)
}

// method ----------------------------------------------------------------------

// 执行方法调用
func (sm *ServerModule) Do(c interface{}, req *Request, metName MethodName) *Response {
	method, err := sm.Method(metName)
	if err != nil {
		return req.ResponseError(E_METHOD_NOT_FOUND, err.Error())
	}

	if reflect.TypeOf(c) != sm.server.ctxType {
		return req.ResponseError(
			E_INTERNAL,
			fmt.Sprintf("jsonrpc error: context type is not [%s]", sm.server.ctxType.String()),
		)
	}

	ctx := reflect.NewAt(sm.server.ctxType.Elem(), unsafe.Pointer(&c))
	input := reflect.New(method.Type.In(2).Elem())
	output := reflect.New(method.Type.In(3).Elem())

	inputData, _ := sm.server.protocol.Encode(req.Params)
	err = sm.server.protocol.Decode(inputData, input.Interface())
	if err != nil {
		return req.ResponseError(
			E_INTERNAL,
			fmt.Sprintf("jsonrpc error: [%s] parse params error: %s", req.Method, err.Error()),
		)
	}

	retVals := method.Func.Call([]reflect.Value{
		sm.entity.GetOriValue(),
		ctx,
		input,
		output,
	})

	if retVals[0].IsNil() {
		return req.ResponseResult(output.Interface())
	} else {
		r := retVals[0].Interface().(*ResponseError)
		return req.ResponseError(r.Code, r.Message)
	}
}
