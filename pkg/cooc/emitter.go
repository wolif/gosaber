package cooc

import (
	"fmt"
)

type event string

const (
	EV_BEFORE_METHOD_INVOKING event = "before method invoking"
	EV_AFTER_METHOD_INVOKING  event = "after method invoking"
)

type emitter struct {
	events map[event]func(...interface{})
}

func (e *emitter) on(ev event, fn func(...interface{}), cover ...bool) error {
	if e.events == nil {
		e.events = make(map[event]func(...interface{}))
	}
	if _, ok := e.events[ev]; !ok || (len(cover) > 0 && cover[0]) {
		e.events[ev] = fn
		return nil
	}
	return fmt.Errorf("public event named [%s] exists", ev)
}

func (e *emitter) off(evs ...event) {
	if e.events == nil {
		return
	}
	for _, ev := range evs {
		delete(e.events, ev)
	}
}

func (e *emitter) emit(ev event, args ...interface{}) bool {
	if e.events == nil || len(e.events) == 0 {
		return false
	}
	fn, ok := e.events[ev]
	if ok {
		fn(args...)
	}
	return true
}
