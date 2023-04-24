package methodreflect

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

type service struct {
	Name     string
	Rcvr     reflect.Value
	Rcvrtype reflect.Type
	Methods  map[string]*serviceMethod
}

type serviceMethod struct {
	Method    reflect.Method
	ArgsType  reflect.Type
	ReplyType reflect.Type
}

type Services struct {
	ServiceMethodDelimiter string              // 服务名和方法名的分割字符串
	Method1stArgType       reflect.Type        // 服务方法第一个参数的类型, 一般为context
	MethodReturnType       []reflect.Type      // 返回值的类型(桉顺序)
	serviceMap             map[string]*service // 服务列表
}

func (s *Services) analysisMethodInput(input string) (service string, method string, err error) {
	// 将需要调用的方法分割为 serviceName 和 methodName
	names := strings.Split(input, s.ServiceMethodDelimiter)
	if len(names) != 2 {
		return "", "", fmt.Errorf("method name format error")
	}
	return names[0], names[1], nil
}

func (s *Services) RegisterMap(m map[string]interface{}) error {
	for name, obj := range m {
		if err := s.Register(name, obj); err != nil {
			return err
		}
	}
	return nil
}

func (s *Services) Register(name string, obj interface{}) error {
	// 如果serviceMap还没初始化, 就初始化
	if s.serviceMap == nil {
		s.serviceMap = make(map[string]*service, 0)
	}

	// 已经出测过相同的服务名
	if _, ok := s.serviceMap[name]; ok {
		return fmt.Errorf("server error")
	}

	// 创建名为name的服务
	ss := &service{
		Name:     name,
		Rcvr:     reflect.ValueOf(obj),
		Rcvrtype: reflect.TypeOf(obj),
		Methods:  make(map[string]*serviceMethod, 0),
	}

	// 利用反射将 结构体的方法 循环解析一遍
	for i := 0; i < ss.Rcvrtype.NumMethod(); i++ {
		method := ss.Rcvrtype.Method(i)
		mtype := method.Type

		// 如果方法的参数数量不是4个, 就不是需要注册的执行方法
		if mtype.NumIn() != 4 {
			continue
		}

		// 第二个参数 必须是指针, 并且类型为 s.Method1stArgType
		reqType := mtype.In(1)
		if reqType.Kind() != reflect.Ptr || reqType.Elem() != s.Method1stArgType {
			continue
		}

		// 第三个参数 必须是指针, 并且外部可以访问
		args := mtype.In(2)
		if args.Kind() != reflect.Ptr || !isExportedOrBuiltin(args) {
			continue
		}

		// 第四个参数 必须是指针, 并且外部可以访问
		reply := mtype.In(3)
		if reply.Kind() != reflect.Ptr || !isExportedOrBuiltin(reply) {
			continue
		}

		// 返回值个数不对, 不是执行方法
		if mtype.NumOut() != len(s.MethodReturnType) {
			continue
		}

		// 判断每个返回值的类型是否和预设的一致
		for i := 0; i < len(s.MethodReturnType); i++ {
			// 如果返回值类型不是 s.ReturnType, 就不是执行方法
			if rt := mtype.Out(i); rt != s.MethodReturnType[i] {
				goto roundEnd
			}
		}

		// 将执行方法放入map中
		ss.Methods[method.Name] = &serviceMethod{
			Method:    method,
			ArgsType:  args.Elem(),
			ReplyType: reply.Elem(),
		}

	roundEnd:
	}

	// 将服务放入map中
	s.serviceMap[name] = ss
	return nil
}

func (s *Services) GetWithName(serviceName, methodName string) (*service, *serviceMethod, error) {
	// 找不到名为 serviceName 的 service
	service, ok := s.serviceMap[serviceName]
	if !ok || service == nil {
		return nil, nil, fmt.Errorf("service named [%s] not found", serviceName)
	}

	// 找不到名为 methodName 的 method
	serviceMethod, ok := service.Methods[methodName]
	if !ok || serviceMethod == nil {
		return nil, nil, fmt.Errorf("method named [%s] not found in service [%s]", methodName, serviceName)
	}

	return service, serviceMethod, nil
}

func (s *Services) Get(method string) (*service, *serviceMethod, error) {
	// 分类 serviceName 和 methodName
	serviceName, methodName, err := s.analysisMethodInput(method)
	if err != nil {
		return nil, nil, err
	}

	return s.GetWithName(serviceName, methodName)
}

func isExported(name string) bool {
	// 是否包外部可访问(名字中首字母是否大写)
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}

func isExportedOrBuiltin(t reflect.Type) bool {
	// 循环剥离指针(解指针), 直到类型不是指针
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return isExported(t.Name()) || t.PkgPath() == "" // 包外部能访问 或者 类型所在的包名为空(golang内置类型)
}
