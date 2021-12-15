package proc

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type Worker struct {
	Name     string
	StopCtx  context.Context
	WG       *sync.WaitGroup
	Fn       func(w *Worker)
	Interval time.Duration
}

func NewWorker(name string, ctx context.Context, wg *sync.WaitGroup, fns ...func(w *Worker)) *Worker {
	w := &Worker{
		Name:    name,
		StopCtx: ctx,
		WG:      wg,
	}
	if len(fns) > 0 {
		w.Fn = fns[0]
	}
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
	go func(w *Worker) {
		w.WG.Add(1)
		defer w.WG.Done()
		log.Printf("worker [ %s ] started", w.Name)
		for {
			select {
			case <-w.StopCtx.Done():
				log.Printf("worker [ %s ] stopped", w.Name)
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
