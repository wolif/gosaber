package proc

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type process struct {
	WG        *sync.WaitGroup
	Ctx       context.Context
	ctxCancel context.CancelFunc
	workers   map[string][]*Worker
	events    map[event]func(proc *process, work *Worker, others ...interface{})
	signalsDo map[os.Signal]func(proc *process) bool
}

func New() *process {
	ctx, cancel := context.WithCancel(context.TODO())
	return NewWithContext(new(sync.WaitGroup), ctx, cancel)
}

func NewWithContext(wg *sync.WaitGroup, ctx context.Context, ctxCancel context.CancelFunc) *process {
	return (&process{
		WG:        wg,
		Ctx:       ctx,
		ctxCancel: ctxCancel,
		workers:   make(map[string][]*Worker),
		events:    make(map[event]func(p *process, w *Worker, others ...interface{})),
		signalsDo: make(map[os.Signal]func(proc *process) bool),
	}).ListenDefault()
}

func (p *process) AddWorker(name string, fn func(w *Worker), numWorker ...int) *process {
	if _, exists := p.Worker(name); exists {
		panic(fmt.Sprintf("worker name [%s] duplicated!", name))
	}
	p.workers[name] = make([]*Worker, 0)
	num := 1
	if len(numWorker) > 0 && numWorker[0] > 0 {
		num = numWorker[0]
	}
	if num == 1 {
		p.workers[name] = append(p.workers[name], NewWorker(name, p.Ctx, p.WG).SetFunc(fn).setProc(p))
		return p
	}
	for i := 1; i <= num; i++ {
		nameTmp := fmt.Sprintf("%sNo.%d", name, i)
		p.workers[name] = append(p.workers[name], NewWorker(nameTmp, p.Ctx, p.WG).SetFunc(fn).setProc(p))
	}
	return p
}

func (p *process) Worker(name string) ([]*Worker, bool) {
	workerGroup, found := p.workers[name]
	return workerGroup, found
}

func (p *process) Run() *process {
	p.emit(start, p, nil)
	for _, workerGroup := range p.workers {
		for _, w := range workerGroup {
			w.Run()
		}
	}
	return p
}

func (p *process) RunInterval() *process {
	p.emit(start, p, nil)
	for _, workerGroup := range p.workers {
		for _, w := range workerGroup {
			w.RunInterval()
		}
	}
	return p
}

func (p *process) Listen(signals []os.Signal, fn func(proc *process) bool) *process {
	if p.signalsDo == nil {
		p.signalsDo = make(map[os.Signal]func(proc *process) bool)
	}
	for _, sign := range signals {
		p.signalsDo[sign] = fn
	}
	return p
}

func (p *process) ListenDefault() *process {
	p.Listen([]os.Signal{syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT}, func(proc *process) bool {
		proc.ctxCancel()
		proc.WG.Wait()
		proc.emit(Exit, proc, nil)
		return true
	})
	p.Listen([]os.Signal{syscall.SIGHUP}, func(proc *process) bool {
		return false
	})
	return p
}

func (p *process) Wait() {
	if p.signalsDo == nil || len(p.signalsDo) == 0 {
		p.ListenDefault()
	}
	signals := make([]os.Signal, 0)
	for sign := range p.signalsDo {
		signals = append(signals, sign)
	}
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, signals...)
	for {
		sign := <-signalChan
		p.emit(Signal, p, nil, sign)
		if p.signalsDo[sign](p) {
			return
		}
	}
}

// -----------------------------------------------------------------------------
func (p *process) On(e event, fn func(proc *process, work *Worker, others ...interface{})) *process {
	if p.events == nil {
		p.events = make(map[event]func(proc *process, work *Worker, others ...interface{}))
	}
	p.events[e] = fn
	return p
}

func (p *process) Off(e event) *process {
	if p.events != nil {
		delete(p.events, e)
	}
	return p
}

func (p *process) emit(e event, proc *process, work *Worker, others ...interface{}) {
	if fn, ok := p.events[e]; ok {
		fn(proc, work, others...)
	}
}

func (p *process) DefaultEvent() *process {
	p.On(start, func(proc *process, work *Worker, others ...interface{}) {
		log.Println("process loaded, start now...")
	})
	p.On(Signal, func(proc *process, work *Worker, others ...interface{}) {
		log.Printf("process get signal: %v ", others[0])
	})
	p.On(Exit, func(proc *process, work *Worker, others ...interface{}) {
		log.Println("process end, stopped now...")
	})
	return p
}
