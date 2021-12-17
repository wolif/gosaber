package proc

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type process struct {
	WG        *sync.WaitGroup
	Ctx       context.Context
	ctxCancel context.CancelFunc
	workers   map[string]*Worker
	events    map[pEvent]func(proc *process, work *Worker, others ...interface{})
	signalsDo map[os.Signal]func(proc *process) bool
}

func New() *process {
	ctx, ctxCancel := context.WithCancel(context.TODO())
	return NewWithContext(new(sync.WaitGroup), ctx, ctxCancel)
}

func NewWithContext(wg *sync.WaitGroup, ctx context.Context, ctxCancel context.CancelFunc) *process {
	return (&process{
		WG:        wg,
		Ctx:       ctx,
		ctxCancel: ctxCancel,
		workers:   make(map[string]*Worker),
		events:    make(map[pEvent]func(proc *process, work *Worker, others ...interface{})),
		signalsDo: make(map[os.Signal]func(proc *process) bool),
	}).ListenDefault()
}

func (p *process) NewWorker(name string, fn func(w *Worker)) *process {
	if _, isAlreadyExisted := p.GetWorkerByName(name); isAlreadyExisted {
		panic(fmt.Sprintf("worker name [%s] duplicated!", name))
	}
	p.workers[name] = NewWorker(name, p.Ctx, p.WG).SetFunc(fn).setProc(p)
	return p
}

func (p *process) NewWorkers(name string, fn func(w *Worker), numWorker ...int) *process {
	num := 1
	if len(numWorker) > 0 && numWorker[0] > 0 {
		num = numWorker[0]
	}
	if num == 1 {
		return p.NewWorker(name, fn)
	}
	for i := 1; i <= num; i++ {
		p.NewWorker(fmt.Sprintf("%s-%d", name, i), fn)
	}
	return p
}

func (p *process) GetWorkerByName(name string) (*Worker, bool) {
	worker, found := p.workers[name]
	return worker, found
}

func (p *process) Run() *process {
	p.emitEv(PEvProcBeforeStart, p, nil)
	for _, g := range p.workers {
		g.Run()
	}
	return p
}

func (p *process) RunInterval() *process {
	p.emitEv(PEvProcBeforeStart, p, nil)
	for _, g := range p.workers {
		g.RunInterval()
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
		proc.emitEv(PEvProcBeforeExit, proc, nil)
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
		p.emitEv(PEvProcWhenGetSigToExit, p, nil, sign)
		if p.signalsDo[sign](p) {
			return
		}
	}
}
