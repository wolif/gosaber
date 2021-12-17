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
}

func New() *process {
	ctx, ctxCancel := context.WithCancel(context.TODO())
	return &process{
		WG:        new(sync.WaitGroup),
		Ctx:       ctx,
		ctxCancel: ctxCancel,
		workers:   make(map[string]*Worker),
		events:    make(map[pEvent]func(proc *process, work *Worker, others ...interface{})),
	}
}

func NewWithContext(wg *sync.WaitGroup, ctx context.Context, ctxCancel context.CancelFunc) *process {
	return &process{
		WG:        wg,
		Ctx:       ctx,
		ctxCancel: ctxCancel,
		workers:   make(map[string]*Worker),
		events:    make(map[pEvent]func(proc *process, work *Worker, others ...interface{})),
	}
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

func (p *process) Wait() {
	p.Listen(func(sign os.Signal) bool {
		switch sign {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			p.ctxCancel()
			p.WG.Wait()
			p.emitEv(PEvProcBeforeExit, p, nil)
			return true
		}
		return false
	})
}

func (p *process) Listen(fn func(sign os.Signal) bool, signals ...os.Signal) {
	if len(signals) == 0 {
		signals = []os.Signal{syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT}
	}
	// 等待系统信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, signals...)
	for {
		sign := <-signalChan
		p.emitEv(PEvProcWhenGetSigToExit, p, nil, sign)
		if fn(sign) {
			return
		}
	}
}
