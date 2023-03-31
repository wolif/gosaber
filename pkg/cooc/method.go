package cooc

import (
	"fmt"
	"reflect"
)

type method struct {
	Name      string
	Service   *service
	fn        reflect.Method
	argsType  reflect.Type
	replyType reflect.Type
	emitter   *emitter
}

func (m *method) On(ev event, fn func(...interface{}), cover ...bool) {
	if m.emitter == nil {
		m.emitter = new(emitter)
	}
	m.emitter.on(ev, fn, cover...)
}

func (m *method) Off(evs ...event) {
	if m.emitter == nil {
		return
	}
	m.emitter.off(evs...)
}

func (m *method) emit(ev event, args ...interface{}) {
	if m.emitter.emit(ev, args...) {
		return
	}
	if m.Service.emitter.emit(ev, args...) {
		return
	}
	m.Service.Collection.emitter.emit(ev, args...)
}

func (m *method) Invoke(ctx interface{}, args []byte) (reply interface{}, methodReturn []interface{}, invokeErr error) {
	argsTmp := reflect.New(m.argsType)
	replyTmp := reflect.New(m.replyType)
	methodReturn = []interface{}{}
	invokeErr = m.Service.Collection.fnAnalysisArg(args, argsTmp.Interface())
	if invokeErr != nil {
		invokeErr = fmt.Errorf("can't analysis the args [ %s ], in service [%s] method [%s]: %v", args, m.Service.Name, m.Name, invokeErr)
		return
	}

	m.emit(EV_BEFORE_METHOD_INVOKING, m, ctx, argsTmp.Interface())
	vals := m.fn.Func.Call([]reflect.Value{
		m.Service.receiver,
		reflect.ValueOf(ctx),
		argsTmp,
		replyTmp,
	})
	for _, v := range vals {
		methodReturn = append(methodReturn, v.Interface())
	}
	reply = replyTmp.Interface()
	m.emit(EV_AFTER_METHOD_INVOKING, m, ctx, reply, &methodReturn)
	return
}

func (m *method) InvokeCustomized(
	ctx interface{},
	args interface{},
	fn func(fn reflect.Method, argType, replyType reflect.Type, emit func(event, ...interface{}), ctx, args interface{}) (interface{}, []interface{}, error)) (reply interface{}, methodReturn []interface{}, invokeErr error,
) {
	return fn(m.fn, m.argsType, m.replyType, m.emit, ctx, args)
}
