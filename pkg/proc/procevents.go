package proc

import (
	"fmt"
	"log"
)

type pEvent int

const (
	PEvProcBeforeStart pEvent = iota
	PEvProcWhenGetSigToExit
	PEvProcBeforeExit
	PEvWorkerBeforeStart
	PEvWorkerBeforeExit
)

func (p *process) OnEv(ev pEvent, fn func(proc *process, work *Worker, others ...interface{})) *process {
	if p.events == nil {
		p.events = make(map[pEvent]func(proc *process, work *Worker, others ...interface{}))
	}
	p.events[ev] = fn
	return p
}

func (p *process) OffEv(ev pEvent) *process {
	if p.events != nil {
		delete(p.events, ev)
	}
	return p
}

func (p *process) emitEv(ev pEvent, proc *process, work *Worker, others ...interface{}) {
	if fn, ok := p.events[ev]; ok {
		fn(proc, work, others...)
	}
}

func (p *process) SetEvDefault() *process {
	p.OnEv(PEvProcBeforeStart, func(proc *process, work *Worker, others ...interface{}) {
		log.Println("process loaded, start now...")
	})
	p.OnEv(PEvProcWhenGetSigToExit, func(proc *process, work *Worker, others ...interface{}) {
		log.Println(fmt.Sprintf("process get signal: %v ", others[0]))
	})
	p.OnEv(PEvProcBeforeExit, func(proc *process, work *Worker, others ...interface{}) {
		log.Println("process end, stopped now...")
	})
	p.OnEv(PEvWorkerBeforeStart, func(proc *process, work *Worker, others ...interface{}) {
		log.Println(fmt.Sprintf("worker [%s] loaded, start now...", work.Name))
	})
	p.OnEv(PEvWorkerBeforeExit, func(proc *process, work *Worker, others ...interface{}) {
		log.Println(fmt.Sprintf("worker [%s] end, stopped now...", work.Name))
	})
	return p
}
