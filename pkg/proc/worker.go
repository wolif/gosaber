package proc

import (
	"context"
	"fmt"
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
	events   map[wEvent]func(w *Worker, others ...interface{})
}

func NewWorker(name string, ctx context.Context, wg *sync.WaitGroup, fns ...func(w *Worker)) *Worker {
	w := &Worker{
		Name:    name,
		StopCtx: ctx,
		WG:      wg,
		events:  make(map[wEvent]func(w *Worker, others ...interface{})),
	}
	if len(fns) > 0 {
		w.Fn = fns[0]
	}
	return w
}

func (w *Worker) emitPEv(ev pEvent) {
	if w.proc != nil {
		w.proc.emitEv(ev, w.proc, w)
	}
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
	w.emitPEv(PEvWorkerBeforeStart)
	w.emitEv(WEventBeforeStart, w)

	w.WG.Add(1)
	go func(w *Worker) {
		defer w.WG.Done()
		for {
			select {
			case <-w.StopCtx.Done():
				w.emitPEv(PEvWorkerBeforeExit)
				w.emitEv(WEventBeforeExit, w)
				return
			case <-time.After(w.Interval):
				w.Fn(w)
			}
		}
	}(w)
}

func (w *Worker) SetNextInterval(interval time.Duration) {
	if interval < 0 {
		interval = 0
	}
	w.Interval = interval
}
