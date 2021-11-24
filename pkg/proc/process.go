package proc

import (
	"context"
	"fmt"
	"github.com/wolif/gosaber/pkg/log"
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
}

func New() *process {
	ctx, ctxCancel := context.WithCancel(context.TODO())
	return &process{
		WG:        new(sync.WaitGroup),
		Ctx:       ctx,
		ctxCancel: ctxCancel,
		workers:   make(map[string]*Worker),
	}
}

func (p *process) NewWorker(name string, fn func(w *Worker)) *process {
	if _, isAlreadyExisted := p.GetWorkerByName(name); isAlreadyExisted {
		panic(fmt.Sprintf("worker name [%s] duplicated!", name))
	}
	p.workers[name] = NewWorker(name, p.Ctx, p.WG).SetFunc(fn)
	return p
}

func (p *process) NewWorkers(name string, fn func(w *Worker), numWorker ...int) *process {
	num := 1
	if len(numWorker) > 0 && numWorker[0] > 0{
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
	for _, g := range p.workers {
		g.Run()
	}
	return p
}

func (p *process) RunInterval() *process {
	for _, g := range p.workers {
		g.RunInterval()
	}
	return p
}

func (p *process) Wait() {
	// 等待系统信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sign := <-signalChan
		log.Infof("process get a signal %s", sign.String())
		switch sign {
		case syscall.SIGHUP:
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			p.ctxCancel()
			p.WG.Wait()
			log.Infof("workers all done, now process done")
			return
		}
	}
}
