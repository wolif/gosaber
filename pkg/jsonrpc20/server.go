package jsonrpc20

import (
	"context"
	"reflect"
	"strings"
	"sync"

	"github.com/wolif/gosaber/pkg/jsonrpc20/protocol"
	"github.com/wolif/gosaber/pkg/ref"
)

type ModuleName = string // 模块名称
type MethodName = string // 方法名称

// jsonrpc服务
type Server struct {
	ctxType     reflect.Type                          // 方法的第一个参数, context类型
	modules     map[ModuleName]*ServerModule          // 注册的所有模块
	resolveFunc func(string) (ModuleName, MethodName) // 将request.method 解析成 module 和 method的方法
	protocol    Protocol                              // 数据编码协议
}

func NewServer(fn func(string) (ModuleName, MethodName)) *Server {
	return &Server{
		ctxType:     reflect.TypeOf((*context.Context)(nil)),
		modules:     make(map[string]*ServerModule),
		resolveFunc: fn,
		protocol:    &protocol.Json{},
	}
}

// setter ----------------------------------------------------------------------

func (s *Server) SetCtxType(t reflect.Type) *Server {
	if t.Kind() != reflect.Ptr {
		panic(ServerError("context type must be pointer"))
	}
	s.ctxType = t
	return s
}

func (s *Server) SetProtocol(p Protocol) *Server {
	s.protocol = p
	return s
}

func (s *Server) SetResolveFunc(fn func(string) (ModuleName, MethodName)) *Server {
	s.resolveFunc = fn
	return s
}

// getter ----------------------------------------------------------------------

func (s *Server) CtxType() reflect.Type {
	return s.ctxType
}

func (s *Server) Module(modName ModuleName) (*ServerModule, error) {
	if sm, ok := s.modules[modName]; ok {
		return sm, nil
	}
	return nil, ServerErrorf("module named [%s] not found", modName)
}

func (s *Server) ResolveFunc() func(string) (ModuleName, MethodName) {
	return s.resolveFunc
}

func (s *Server) Protocol() Protocol {
	return s.protocol
}

// register ----------------------------------------------------------------------

// 注册命名模块
func (s *Server) RegisterNamedModules(modules map[ModuleName]interface{}) *Server {
	for name, module := range modules {
		s.RegisterModule(module, name)
	}
	return s
}

// 注册默认名字模块
func (s *Server) RegisterModules(modules ...interface{}) *Server {
	for _, module := range modules {
		s.RegisterModule(module)
	}
	return s
}

// 注册模块, 模块必须是 结构体的指针, 否则会 中断程序
// 可以给模块取名, 否则名称为结构地的名称
func (s *Server) RegisterModule(module interface{}, name ...string) *Server {
	rm := ref.New(module)
	if !rm.IsStruct() || !rm.IsPtr() {
		panic(ServerErrorf("the module named [%s] must be a struct ptr!", name))
	}
	modName := strings.SplitN(rm.GetOriType().String(), ".", 2)[1]
	if len(name) > 0 {
		modName = name[0]
	}

	if s.modules == nil {
		s.modules = make(map[ModuleName]*ServerModule)
	}
	if _, ok := s.modules[modName]; ok {
		panic(ServerErrorf("the module named [%s] is already registered!", modName))
	}
	s.modules[modName] = &ServerModule{name: modName, server: s, entity: rm}
	s.modules[modName].init()

	return s
}

// dispatch ----------------------------------------------------------------------

// 分配并执行jsonrpc请求
// 参数 ctx 是context接口类型
// 参数 body 是请求体数据
// 返回数据
// 		当 request 只有一个请求时为 *Response, 否则为 []*Response
// 		当无法解析参数 body 时, 会返回单个包含错误的 *Response
// 附:处理多个 request 时, 会进行并发处理
func (s *Server) Dispatch(ctx interface{}, body []byte) interface{} {
	req := new(Request)
	err := s.protocol.Decode(body, req)
	if err == nil { // 单个请求
		return s.do(ctx, req)
	}
	reqs := make([]*Request, 0)
	err = s.protocol.Decode(body, &reqs)
	if err == nil { // 多个请求
		resps := make([]*Response, 0)
		wg := new(sync.WaitGroup)
		for _, req := range reqs {
			wg.Add(1)
			go func(req *Request) {
				defer wg.Done()
				resps = append(resps, s.do(ctx, req))
			}(req)
		}
		wg.Wait()
		return resps
	}

	// 无法解析请求数据
	return NewErrorResponse(E_PARSE)
}

func (s *Server) do(ctx interface{}, req *Request) *Response {
	modName, _ := s.resolveFunc(req.Method)
	module, err := s.Module(modName)
	if err != nil {
		return req.ResponseError(E_METHOD_NOT_FOUND, err.Error())
	}
	return module.Do(ctx, req)
}
