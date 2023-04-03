package proc

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type Worker struct {
	proc     *process
	Name     string
	StopCtx  context.Context
	WG       *sync.WaitGroup
	Fn       func(w *Worker)
	Interval time.Duration
	events   map[event]func(w *Worker, others ...interface{})
}

func NewWorker(name string, ctx context.Context, wg *sync.WaitGroup, fns ...func(w *Worker)) *Worker {
	w := &Worker{
		Name:    name,
		StopCtx: ctx,
		WG:      wg,
		events:  make(map[event]func(w *Worker, others ...interface{})),
	}
	if len(fns) > 0 {
		w.Fn = fns[0]
	}
	return w.DefaultEvent()
}

func (w *Worker) setProc(p *process) *Worker {
	w.proc = p
	return w
}

func (w *Worker) SetFunc(fn func(g *Worker)) *Worker {
	w.Fn = fn
	return w
}

func (w *Worker) checkBeforeRun() {
	if w.Fn == nil {
		panic(fmt.Sprintf("worker named [%s] func not found, can't execute", w.Name))
	}
}

func (w *Worker) Run() {
	w.checkBeforeRun()
	go w.Fn(w)
}

func (w *Worker) RunInterval() {
	w.checkBeforeRun()
	w.emit(start)

	w.WG.Add(1)
	go func(w *Worker) {
		defer w.WG.Done()
		for {
			select {
			case <-w.StopCtx.Done():
				w.emit(Exit)
				return
			case <-time.After(w.Interval):
				w.Fn(w)
			}
		}
	}(w)
}

func (w *Worker) SetInterval(interval time.Duration) {
	if interval < 0 {
		interval = 0
	}
	w.Interval = interval
}

// -----------------------------------------------------------------------------
func (w *Worker) On(e event, fn func(w *Worker, others ...interface{})) *Worker {
	if w.events == nil {
		w.events = make(map[event]func(w *Worker, others ...interface{}))
	}
	w.events[e] = fn
	return w
}

func (w *Worker) Off(e event) *Worker {
	if w.events[e] != nil {
		delete(w.events, e)
	}
	return w
}

func (w *Worker) emit(e event, others ...interface{}) {
	if fn, ok := w.events[e]; ok {
		fn(w, others...)
	}
}

func (w *Worker) DefaultEvent() *Worker {
	w.On(start, func(work *Worker, others ...interface{}) {
		log.Printf("worker [%s] loaded, start now...", work.Name)
	})
	w.On(Exit, func(work *Worker, others ...interface{}) {
		log.Printf("worker [%s] end, stopped now...", work.Name)
	})
	return w
}
