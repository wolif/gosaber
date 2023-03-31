package cooc

import "reflect"

type service struct {
	Name         string
	Collection   *cooc
	receiver     reflect.Value
	receiverType reflect.Type
	methods      map[string]*method
	emitter      *emitter
}

func (s *service) On(ev event, fn func(...interface{}), cover ...bool) {
	if s.emitter == nil {
		s.emitter = new(emitter)
	}
	s.emitter.on(ev, fn, cover...)
}

func (s *service) Off(evs ...event) {
	if s.emitter == nil {
		return
	}
	s.emitter.off(evs...)
}

func (s *service) Method(name string) *method {
	if name == "" {
		return nil
	}
	if m, ok := s.methods[name]; ok {
		return m
	}
	return nil
}
