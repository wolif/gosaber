package cooc

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type cooc struct {
	ctxArgType    reflect.Type
	methodRetType []reflect.Type
	services      map[string]*service
	emitter       *emitter
	fnAnalysisArg func(src []byte, dist interface{}) error
}

func New(ctxArgType reflect.Type, opts ...interface{}) *cooc {
	c := &cooc{
		ctxArgType:    ctxArgType,
		methodRetType: []reflect.Type{},
		fnAnalysisArg: func(src []byte, dist interface{}) error { return json.Unmarshal(src, dist) },
		services:      map[string]*service{},
		emitter:       &emitter{},
	}
	for _, opt := range opts {
		switch o := opt.(type) {
		case []reflect.Type:
			c.methodRetType = o
		case func(src []byte, dist interface{}) error:
			c.fnAnalysisArg = o
		}
	}
	return c
}

func (c *cooc) SetMethodRetType(methodRetType ...reflect.Type) *cooc {
	c.methodRetType = methodRetType
	return c
}

func (c *cooc) SetArgAnalysisFn(fn func(src []byte, dist interface{}) error) *cooc {
	c.fnAnalysisArg = fn
	return c
}

// -----------------------------------------------------------------------------

func (c *cooc) On(ev event, fn func(...interface{}), cover ...bool) {
	c.emitter.on(ev, fn, cover...)
}

func (c *cooc) ServicesOn(serviceName []string, ev event, fn func(...interface{})) {
	for _, sn := range serviceName {
		if s := c.Service(sn); s != nil {
			s.On(ev, fn)
		}
	}
}

func (c *cooc) MethodsOn(serviceMethods map[string][]string, ev event, fn func(...interface{})) {
	for sn, mns := range serviceMethods {
		for _, mn := range mns {
			if m := c.Method(sn, mn); m != nil {
				m.On(ev, fn)
			}
		}
	}
}

func (c *cooc) Off(evs ...event) {
	if c.emitter == nil {
		return
	}
	c.emitter.off(evs...)
}

func (c *cooc) ServicesOff(serviceName []string, evs ...event) {
	for _, sn := range serviceName {
		if s := c.Service(sn); s != nil {
			s.Off(evs...)
		}
	}
}

func (c *cooc) MethodsOff(serviceMethods map[string][]string, evs ...event) {
	for sn, mns := range serviceMethods {
		for _, mn := range mns {
			if m := c.Method(sn, mn); m != nil {
				m.Off(evs...)
			}
		}
	}
}

// -----------------------------------------------------------------------------

func (c *cooc) Add(name string, obj interface{}) map[string]error {
	ret := make(map[string]error)
	if len(name) == 0 {
		panic("service name can't be empty")
	}

	if c.services == nil {
		c.services = make(map[string]*service)
	}

	srv := reflect.ValueOf(obj)
	if srv.Type().Kind() != reflect.Ptr || srv.Elem().Type().Kind() != reflect.Struct {
		panic(fmt.Sprintf("service named [%s] is not a struct ptr", name))
	}

	s := &service{
		Name:         name,
		Collection:   c,
		receiver:     reflect.ValueOf(obj),
		receiverType: nil,
		methods:      map[string]*method{},
		emitter:      &emitter{},
	}

	for i := 0; i < srv.Type().NumMethod(); i++ {
		m := srv.Type().Method(i)
		if m.Type.NumIn() != 4 {
			ret[m.Name] = fmt.Errorf("number of args != 3")
			continue
		}
		ctxType := m.Type.In(1)
		if ctxType.Kind() != reflect.Ptr || ctxType.Elem() != c.ctxArgType {
			ret[m.Name] = fmt.Errorf("ctx type inconformity")
			continue
		}

		argType := m.Type.In(2)
		if argType.Kind() != reflect.Ptr || !isExportedOrBuiltin(argType) {
			ret[m.Name] = fmt.Errorf("args type inconformity")
			continue
		}

		replyType := m.Type.In(3)
		if replyType.Kind() != reflect.Ptr || !isExportedOrBuiltin(replyType) {
			ret[m.Name] = fmt.Errorf("reply type inconformity")
			continue
		}

		if m.Type.NumOut() != len(c.methodRetType) {
			ret[m.Name] = fmt.Errorf("number of output inconformity")
			continue
		}

		for i := 0; i < len(c.methodRetType); i++ {
			if m.Type.Out(i) != c.methodRetType[i] {
				ret[m.Name] = fmt.Errorf("type of output[%d] inconformity", i)
				goto ROUNDEND
			}
		}

		s.methods[m.Name] = &method{
			Name:      m.Name,
			Service:   s,
			fn:        m,
			argsType:  argType.Elem(),
			replyType: replyType.Elem(),
			emitter:   &emitter{},
		}
		ret[m.Name] = nil
	ROUNDEND:
	}

	c.services[name] = s
	return ret
}

func (c *cooc) AddMany(services map[string]interface{}) map[string]map[string]error {
	ret := make(map[string]map[string]error)
	for name, service := range services {
		ret[name] = c.Add(name, service)
	}
	return ret
}

func (c *cooc) Service(name string) *service {
	if name == "" {
		return nil
	}
	if m, ok := c.services[name]; ok {
		return m
	}
	return nil
}

func (c *cooc) Method(serviceName, methodName string) *method {
	s := c.Service(serviceName)
	if s != nil {
		return s.Method(methodName)
	}
	return nil
}
